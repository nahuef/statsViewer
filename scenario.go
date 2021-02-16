package main

import "html/template"

// DateAvg has processed values for a specific date.
type DateAvg struct {
	Score        float64
	Grouped      int
	PercentagePB int
}

// Scenario contains all the scenario data including challenges.
type Scenario struct {
	fileName       string
	Name           string
	TimesPlayed    int
	Challenges     []Challenge
	Highscore      float64
	Lowscore       float64
	LowestAvgScore float64
	ByDateMax      []map[string]Challenge
	ByDateAvg      []map[string]DateAvg
	ChartByDate    template.HTML
}
