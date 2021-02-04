package main

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// AddCharts ...
func AddCharts(stats *Stats) {
	for _, scen := range stats.Scenarios {
		LineChart(scen)
	}
}

// LineChart ...
func LineChart(scen *Scenario) {
	lineMax := charts.NewLine()
	lineMax.Renderer = newSnippetRenderer(lineMax, lineMax.Validate)
	lineMax.SetGlobalOptions(chartGlobalOpts(max, scen.Name, len(scen.ByDateMax)))

	lineAvg := charts.NewLine()
	lineAvg.Renderer = newSnippetRenderer(lineAvg, lineAvg.Validate)
	lineAvg.SetGlobalOptions(chartGlobalOpts(avg, scen.Name, len(scen.ByDateAvg)))

	maxDates := []string{}
	maxScores := []opts.LineData{}
	for kDate, vScore := range scen.ByDateMax {
		maxDates = append(maxDates, kDate)
		maxScores = append(maxScores, opts.LineData{Value: vScore})
	}

	avgDates := []string{}
	avgScores := []opts.LineData{}
	for kDate, vScore := range scen.ByDateAvg {
		avgDates = append(avgDates, kDate)
		avgScores = append(avgScores, opts.LineData{Value: vScore})
	}

	lineMax.SetXAxis(maxDates).
		AddSeries("Score", maxScores).
		SetSeriesOptions(seriesOpts)
	var htmlSnippet = renderToHTML(lineMax)
	scen.ChartByDateMax = htmlSnippet

	lineAvg.SetXAxis(avgDates).
		AddSeries("Score", avgScores).
		SetSeriesOptions(seriesOpts)
	var htmlSnippetAvg = renderToHTML(lineAvg)
	scen.ChartByDateAvg = htmlSnippetAvg
}

var max = "max"
var avg = "avg"

var seriesOpts = charts.WithLabelOpts(opts.Label{Show: true, Color: "black"})

func titleOpts(group string, scenName string, length int) charts.GlobalOpts {
	var groupedBy string
	switch group {
	case "max":
		groupedBy = "highest"
	case "avg":
		groupedBy = "average"
	}

	return charts.WithTitleOpts(opts.Title{
		Title:    scenName,
		Subtitle: fmt.Sprintf("Grouped by %v score. Datapoints (days): %v", groupedBy, length),
	})
}

var dataZoom = opts.ToolBoxFeatureDataZoom{
	Show:  true,
	Title: map[string]string{"zoom": "Zoom", "back": "Restore"},
}

func saveAsImage(fileName string) *opts.ToolBoxFeatureSaveAsImage {
	return &opts.ToolBoxFeatureSaveAsImage{
		Show:  true,
		Name:  fileName,
		Title: "Download",
	}
}

func toolBoxFeatures(scenName string) *opts.ToolBoxFeature {
	return &opts.ToolBoxFeature{
		SaveAsImage: saveAsImage(scenName),
		DataZoom:    &dataZoom,
	}
}

func toolBoxOpts(scenName string) charts.GlobalOpts {
	return charts.WithToolboxOpts(opts.Toolbox{
		Show:    true,
		Feature: toolBoxFeatures(scenName),
	})
}

var tooltipOpts = charts.WithTooltipOpts(opts.Tooltip{
	Show: true,
})

var xAxisLabelFormatter = opts.AxisLabel{
	Rotate: 45,
}

var xAxisOpts = charts.WithXAxisOpts(opts.XAxis{
	AxisLabel: &xAxisLabelFormatter,
})

var yAxisLabelFormatter = opts.AxisLabel{}

var yAxisOpts = charts.WithYAxisOpts(opts.YAxis{
	Type:      "value",
	Scale:     true,
	AxisLabel: &yAxisLabelFormatter,
})

func chartGlobalOpts(groupedBy string, scenName string, length int) (charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts) {
	return titleOpts(groupedBy, scenName, length), toolBoxOpts(scenName), tooltipOpts, xAxisOpts, yAxisOpts
}
