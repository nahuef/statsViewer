package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gosuri/uiprogress"
)

var extractor = Extract{}

// Challenge ...
type Challenge struct {
	Name      string
	Datetime  string
	Date      string
	Time      string
	Score     float64
	SensScale string
	HSens     float64
	VSens     float64
	FOV       float64
}

func (c *Challenge) sensStr() string {
	if c.HSens == c.VSens {
		return fmt.Sprintf("Sens: %v %v", c.HSens, c.SensScale)
	}

	return fmt.Sprintf("Vsens: %v, Hsens: %v %v", c.VSens, c.HSens, c.SensScale)
}

// Scenario ...
type Scenario struct {
	fileName    string
	Name        string
	TimesPlayed int
	Challenges  []Challenge
	Highscore   float64
	Lowscore    float64
	LowestAvg   float64
	ByDateMax   []map[string]Challenge
	// ByDateAvg []interface{}: [0]float64 score, [1]int # of grouped challenges
	ByDateAvg   []map[string][]interface{}
	ChartByDate template.HTML
}

// Stats ...
type Stats struct {
	Scenarios         map[string]*Scenario
	SortedTimesPlayed []*Scenario
	UniqueDays        []string
	DaysPlayed        int
	TotalScens        int
	TotalPlayed       int
}

func scenarioWorker(scen *Scenario, sortedTimesPlayed *[]*Scenario, uniqueDays *[]string, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	scen.Lowscore = scen.Challenges[0].Score
	ByDate := map[string][]Challenge{}

	for _, chall := range scen.Challenges {
		if chall.Score > scen.Highscore {
			scen.Highscore = chall.Score
		}
		if chall.Score < scen.Lowscore {
			scen.Lowscore = chall.Score
		}

		ByDate[chall.Date] = append(ByDate[chall.Date], chall)
	}

	// Group challenges per date
	// max: a key per date containing one challenge
	// avg: a key per date containing a slice with average score and number of grouped challenges
	max, avg := Group(ByDate)

	// maps into a slice so we can sort them by date
	for k, v := range max {
		scen.ByDateMax = append(scen.ByDateMax, map[string]Challenge{k: v})
	}
	scen.LowestAvg = scen.Highscore
	for k, v := range avg {
		scen.ByDateAvg = append(scen.ByDateAvg, map[string][]interface{}{k: v})
		score := v[0].(float64)
		if score < scen.LowestAvg {
			scen.LowestAvg = score
		}
	}

	// Actually sort them by date (descending)
	sort.SliceStable(scen.ByDateMax, func(i, j int) bool {
		var iDate int
		for k := range scen.ByDateMax[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range scen.ByDateMax[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})
	sort.SliceStable(scen.ByDateAvg, func(i, j int) bool {
		var iDate int
		for k := range scen.ByDateAvg[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range scen.ByDateAvg[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})

	mux.Lock()
	defer mux.Unlock()
	for k := range ByDate {
		if !ContainsString(*uniqueDays, k) {
			*uniqueDays = append(*uniqueDays, k)
		}
	}
	*sortedTimesPlayed = append(*sortedTimesPlayed, scen)

	// Less than 2 datapoints or 3 challenges => skip chart
	if scen.TimesPlayed <= 2 || len(ByDate) <= 1 {
		return
	}
	AddLineChart(scen)
}

func (s *Stats) forEachScenario() {
	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	for _, scen := range s.Scenarios {
		wg.Add(1)
		go scenarioWorker(scen, &s.SortedTimesPlayed, &s.UniqueDays, &wg, mux)
	}
	wg.Wait()

}

func fileWorker(stats *Stats, file os.FileInfo, wg *sync.WaitGroup, mux *sync.Mutex, bar *uiprogress.Bar) {
	defer wg.Done()
	if file.IsDir() == true {
		return
	}

	// Open file
	f, err := os.Open(StatsPath + file.Name())
	Check(err)
	defer f.Close()

	// New challenge
	challenge := Challenge{}

	// Read line by line
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		err = s.Err()
		Check(err)

		extractor := Extract{line: line, fileName: file.Name(), challenge: &challenge}
		extractor.extractData()
	}

	mux.Lock()
	stats.TotalPlayed++
	if _, ok := stats.Scenarios[challenge.Name]; ok {
		stats.Scenarios[challenge.Name].TimesPlayed++
		stats.Scenarios[challenge.Name].Challenges = append(stats.Scenarios[challenge.Name].Challenges, challenge)
	} else {
		stats.TotalScens++
		stats.Scenarios[challenge.Name] = &Scenario{
			fileName:    file.Name(),
			Name:        challenge.Name,
			TimesPlayed: 1,
			Challenges:  []Challenge{challenge},
		}
	}
	mux.Unlock()
	bar.Incr()
}

// ParseStats ...
func ParseStats(files []os.FileInfo) Stats {
	stats := Stats{
		Scenarios: map[string]*Scenario{},
	}

	bar := uiprogress.AddBar(len(files)).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Files (%d/%d)", b.Current(), len(files))
	})
	uiprogress.Start()

	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go fileWorker(&stats, file, &wg, mux, bar)
	}
	wg.Wait()
	uiprogress.Stop()

	stats.forEachScenario()
	stats.DaysPlayed = len(stats.UniqueDays)
	sort.SliceStable(stats.SortedTimesPlayed, func(i, j int) bool {
		return stats.SortedTimesPlayed[i].TimesPlayed > stats.SortedTimesPlayed[j].TimesPlayed
	})

	return stats
}
