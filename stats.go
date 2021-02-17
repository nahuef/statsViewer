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

// Stats ...
type Stats struct {
	Scenarios         map[string]*Scenario
	SortedTimesPlayed []*Scenario
	UniqueDays        map[string]int
	DaysPlayed        int
	TotalScens        int
	TotalPlayed       int
	PerformanceChart  template.HTML
	WordCloud         template.HTML
}

func newStats() *Stats {
	return &Stats{
		Scenarios:  make(map[string]*Scenario),
		UniqueDays: make(map[string]int),
	}
}

func scenarioWorker(scen *Scenario, sortedTimesPlayed *[]*Scenario, uniqueDays *map[string]int, wg *sync.WaitGroup, mux *sync.Mutex) {
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
	groupedMax, groupedAvg := Group(ByDate, scen.Highscore, scen.Name)

	// Maps into a slice so we can sort them
	for date, challenge := range groupedMax {
		scen.ByDateMax = append(scen.ByDateMax, map[string]Challenge{date: challenge})
	}
	scen.LowestAvgScore = scen.Highscore
	for date, dateAvg := range groupedAvg {
		scen.ByDateAvg = append(scen.ByDateAvg, map[string]DateAvg{date: dateAvg})
		if dateAvg.Score < scen.LowestAvgScore {
			scen.LowestAvgScore = dateAvg.Score
		}
	}

	// Actually sort by date (descending)
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
	*sortedTimesPlayed = append(*sortedTimesPlayed, scen)
	// Less than 2 datapoints or 3 challenges => skip chart
	if scen.TimesPlayed <= 2 || len(ByDate) <= 1 {
		return
	}
	scen.ChartByDate = ScenarioLineChart(scen)
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

	s := bufio.NewScanner(f)
	// For each line
	for s.Scan() {
		line := s.Text()
		Check(s.Err())

		extractor := Extract{line: line, fileName: file.Name(), challenge: &challenge}
		extractor.extractData()
	}

	if valid := challenge.IsValid(); valid == false {
		return
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
func ParseStats(files []os.FileInfo) *Stats {
	stats := newStats()

	bar := uiprogress.AddBar(len(files)).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Files (%d/%d)", b.Current(), len(files))
	})
	uiprogress.Start()

	mux := &sync.Mutex{}
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go fileWorker(stats, file, &wg, mux, bar)
	}
	wg.Wait()
	uiprogress.Stop()

	stats.forEachScenario()
	percentagesToPBByDate(stats)
	stats.WordCloud = WordCloud(&stats.Scenarios)
	stats.PerformanceChart = PerformanceChart(&stats.UniqueDays)
	stats.DaysPlayed = len(stats.UniqueDays)
	sort.SliceStable(stats.SortedTimesPlayed, func(i, j int) bool {
		return stats.SortedTimesPlayed[i].TimesPlayed > stats.SortedTimesPlayed[j].TimesPlayed
	})

	return stats
}

func percentagesToPBByDate(s *Stats) {
	// Each date will have a list of scenario averages
	dates := make(map[string][]int)

	for _, scen := range s.Scenarios {
		for _, dateAndDateAvg := range scen.ByDateAvg {
			for date, DateAvg := range dateAndDateAvg {
				dates[date] = append(dates[date], DateAvg.PercentagePB)
			}
		}
	}

	for date, percentages := range dates {
		amount := len(percentages)
		if amount <= 5 {
			continue
		}

		var sum int
		for _, v := range percentages {
			sum += v
		}
		avgPercentage := sum / amount

		s.UniqueDays[date] = avgPercentage
	}
}
