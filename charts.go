package main

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// AddCharts ...
func AddCharts(stats *Stats) {
	for _, scen := range stats.Scenarios {
		LineShowLabel(scen)
	}
}

// LineShowLabel ...
func LineShowLabel(scen *Scenario) {
	line := charts.NewLine()
	line.Renderer = newSnippetRenderer(line, line.Validate)

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			// Title:    scen.Name,
			Subtitle: fmt.Sprintf("Total days: %v", len(scen.ByDateMax)),
		}),
	)

	dates := []string{}
	maxScores := make([]opts.LineData, 0)
	for kDate, vScore := range scen.ByDateMax {
		dates = append(dates, kDate)
		maxScores = append(maxScores, opts.LineData{Value: vScore})
	}

	line.SetXAxis(dates).
		AddSeries("Score", maxScores).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}),
		)

	var htmlSnippet = renderToHTML(line)
	scen.Chart = htmlSnippet
}
