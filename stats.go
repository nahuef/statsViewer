package main

import (
	"bufio"
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
}

// Stats ...
type Stats struct {
	Scenarios   map[string]*Scenario
	Sorted      []*Scenario
	UniqueScens int
	TotalPlayed int
}

// TODO:
// var dateDesc = "dateDesc"
// var dateAsc = "dateAsc"

var timesPlayed = "timesPlayed"

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
			stats.UniqueScens++
			stats.Scenarios[challenge.Name] = &Scenario{
				fileName:    file.Name(),
				Name:        challenge.Name,
				TimesPlayed: 1,
				Challenges:  []Challenge{challenge},
			}
		}
	}

	stats.sortBy(timesPlayed)
	return stats
}
