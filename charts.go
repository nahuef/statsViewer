package main

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// LineChart ...
func LineChart(scen *Scenario) {
	lineMax := charts.NewLine()
	lineMax.Renderer = newSnippetRenderer(lineMax, lineMax.Validate)
	lineMax.SetGlobalOptions(chartGlobalOpts(max, scen.Name, len(scen.ByDateMax), scen.Highscore, scen.LowestAvg))

	lineAvg := charts.NewLine()
	lineAvg.Renderer = newSnippetRenderer(lineAvg, lineAvg.Validate)
	lineAvg.SetGlobalOptions(chartGlobalOpts(avg, scen.Name, len(scen.ByDateAvg), scen.Highscore, scen.LowestAvg))

	maxDates := []string{}
	maxScores := []opts.LineData{}
	for _, dateScore := range scen.ByDateMax {
		for k, v := range dateScore {
			maxDates = append(maxDates, SimplifyDate(k))
			maxScores = append(maxScores, opts.LineData{Value: v})
		}
	}

	avgDates := []string{}
	avgScores := []opts.LineData{}
	for _, dateScore := range scen.ByDateAvg {
		for k, v := range dateScore {
			avgDates = append(avgDates, SimplifyDate(k))
			avgScores = append(avgScores, opts.LineData{Value: v})
		}
	}

	lineMax.SetXAxis(maxDates).
		AddSeries("Max score", maxScores).
		AddSeries("Avg score", avgScores).
		SetSeriesOptions(seriesOpts)
	var htmlSnippet = renderToHTML(lineMax)
	scen.ChartByDateMax = htmlSnippet
}

var max = "max"
var avg = "avg"

var seriesOpts = charts.WithLabelOpts(opts.Label{Show: true, Color: "black"})

func titleOpts(scenName string, length int) charts.GlobalOpts {
	return charts.WithTitleOpts(opts.Title{
		Title:    scenName,
		Subtitle: fmt.Sprintf("Grouped by days, %v datapoints.", length),
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

func yAxisOpts(highscore, lowestAvg float64) charts.GlobalOpts {
	return charts.WithYAxisOpts(opts.YAxis{
		Type: "value",
		// Scale:     true,
		Max:       10 * ((int(highscore*1.05) + 9) / 10),
		Min:       10 * (int(lowestAvg*0.95) / 10),
		AxisLabel: &yAxisLabelFormatter,
	})
}

func chartGlobalOpts(groupedBy string, scenName string, length int, hs float64, ls float64) (charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts) {
	return titleOpts(scenName, length), toolBoxOpts(scenName), tooltipOpts, xAxisOpts, yAxisOpts(hs, ls)
}
