package main

import (
	"bufio"
	"html/template"
	"os"
	"sort"
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

// Scenario ...
type Scenario struct {
	fileName    string
	Name        string
	TimesPlayed int
	Challenges  []Challenge
	Chart       template.HTML
	ByDateMax   map[string]float64
	ByDateAvg   map[string]float64
}

// Stats ...
type Stats struct {
	Scenarios   map[string]*Scenario
	Sorted      []*Scenario
	TotalScens  int
	TotalPlayed int
}

// TODO:
// var dateDesc = "dateDesc"
// var dateAsc = "dateAsc"

var timesPlayed = "timesPlayed"

// for kDate, vScore := range scen.ByDateMax {
// 	date = append(date, kDate)
// 	scores = append(scores, opts.LineData{Value: vScore})
// }

func (s *Stats) forEach() {
	for _, scen := range s.Scenarios {
		ByDate := make(map[string][]float64, 0)
		for _, chall := range scen.Challenges {
			ByDate[chall.Date] = append(ByDate[chall.Date], chall.Score)
		}

		max, avg := group(ByDate)
		scen.ByDateMax = max
		scen.ByDateAvg = avg
		// fmt.Printf("scen %+v:", scen)
	}
}

func group(ByDate map[string][]float64) (map[string]float64, map[string]float64) {
	ByDateMax := make(map[string]float64, 0)
	ByDateAvg := make(map[string]float64, 0)

	for k, v := range ByDate {
		var max float64
		var avg float64
		var sum float64
		for i, e := range v {
			if i == 0 || e > max {
				max = e
			}
			sum += e
		}

		avg = sum / float64(len(v))

		ByDateMax[k] = float64(int(max*10)) / 10
		ByDateAvg[k] = float64(int(avg*10)) / 10
	}

	return ByDateMax, ByDateAvg
}

func (s *Stats) sortBy(condition string) {
	var sorted []*Scenario
	for _, scen := range s.Scenarios {
		sorted = append(sorted, scen)
	}

	switch condition {
	case timesPlayed:
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].TimesPlayed > sorted[j].TimesPlayed
		})
	}

	s.Sorted = sorted
}

// ParseStats ...
func ParseStats(files []os.FileInfo) Stats {
	stats := Stats{
		Scenarios: make(map[string]*Scenario),
	}

	for _, file := range files {
		if file.IsDir() == true {
			continue
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
	}

	// Sort & group
	stats.forEach()
	AddCharts(&stats)
	stats.sortBy(timesPlayed)
	return stats
}
