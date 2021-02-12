package main

import "html/template"

// Scenario contains all the scenario data including challenges.
type Scenario struct {
	fileName    string
	Name        string
	TimesPlayed int
	Challenges  []Challenge
	Highscore   float64
	Lowscore    float64
	LowestAvg   float64
	ByDateMax   []map[string]Challenge
	ByDateAvg   []map[string][]interface{} // []interface{}: [0]float64 score, [1]int # of grouped challenges
	ChartByDate template.HTML
}
