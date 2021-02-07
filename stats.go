package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
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
	ByDateAvg   []map[string]float64
	ChartByDate template.HTML
}

// Stats ...
type Stats struct {
	Scenarios         map[string]*Scenario
	SortedTimesPlayed []*Scenario
	TotalScens        int
	TotalPlayed       int
}

func (s *Stats) forEachScenario() {
	var sortedTimesPlayed []*Scenario

	for _, scen := range s.Scenarios {
		sortedTimesPlayed = append(sortedTimesPlayed, scen)
		sort.SliceStable(sortedTimesPlayed, func(i, j int) bool {
			return sortedTimesPlayed[i].TimesPlayed > sortedTimesPlayed[j].TimesPlayed
		})

		scen.Lowscore = scen.Challenges[0].Score

		ByDate := map[string][]Challenge{}
		for _, chall := range scen.Challenges {
			ByDate[chall.Date] = append(ByDate[chall.Date], chall)
			if chall.Score > scen.Highscore {
				scen.Highscore = chall.Score
			}
			if chall.Score < scen.Lowscore {
				scen.Lowscore = chall.Score
			}
		}

		max, avg := Group(ByDate)

		for k, v := range max {
			scen.ByDateMax = append(scen.ByDateMax, map[string]Challenge{k: v})
		}

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

		scen.LowestAvg = scen.Highscore
		for k, v := range avg {
			scen.ByDateAvg = append(scen.ByDateAvg, map[string]float64{k: v})
			if v < scen.LowestAvg {
				scen.LowestAvg = v
			}
		}

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

		if scen.TimesPlayed <= 2 || len(ByDate) <= 1 {
			continue
		}

		AddLineChart(scen)
	}

	s.SortedTimesPlayed = sortedTimesPlayed
}

// ParseStats ...
func ParseStats(files []os.FileInfo) Stats {
	stats := Stats{
		Scenarios: map[string]*Scenario{},
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

	stats.forEachScenario()
	return stats
}
