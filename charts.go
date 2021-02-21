package main

import (
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// ScenarioLineChart ...
func ScenarioLineChart(scen *Scenario) template.HTML {
	line := charts.NewLine()
	line.Renderer = newSnippetRenderer(line, line.Validate)
	line.SetGlobalOptions(chartGlobalOpts(max, scen.Name, len(scen.ByDateMax), scen.Highscore, scen.LowestAvgScore))

	maxDates := []string{}
	maxScores := []opts.LineData{}
	for _, dateScore := range scen.ByDateMax {
		for date, chall := range dateScore {
			maxDates = append(maxDates, SimplifyDate(date))
			maxScores = append(maxScores, opts.LineData{
				Name:  fmt.Sprintf("%v: %v. FOV: %v. %v", SimplifyDate(date), chall.Score, chall.FOV, chall.SensStr()),
				Value: chall.Score,
			})
		}
	}

	avgDates := []string{}
	avgScores := []opts.LineData{}
	for _, dateScore := range scen.ByDateAvg {
		for date, data := range dateScore {
			avgDates = append(avgDates, SimplifyDate(date))
			avgScores = append(avgScores, opts.LineData{
				Name:  fmt.Sprintf("%v: %v. Grouped: %v", SimplifyDate(date), data.Score, data.Grouped),
				Value: data.Score,
			})
		}
	}

	line.SetXAxis(maxDates).
		AddSeries("Max scores", maxScores).
		AddSeries("Average scores", avgScores).
		SetSeriesOptions(seriesOpts...)

	return renderToHTML(line)
}

// PerformanceChart ...
func PerformanceChart(uniqueDays *map[string]int) template.HTML {
	progress := charts.NewLine()
	progress.Renderer = newSnippetRenderer(progress, progress.Validate)

	// Order map by date
	orderedDates := []map[string]int{}
	for k, v := range *uniqueDays {
		orderedDates = append(orderedDates, map[string]int{k: v})
	}
	sort.SliceStable(orderedDates, func(i, j int) bool {
		var iDate int
		for k := range orderedDates[i] {
			iDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		var jDate int
		for k := range orderedDates[j] {
			jDate, _ = strconv.Atoi(strings.ReplaceAll(k, ".", ""))
		}
		return iDate < jDate
	})

	lowestAvgPb := 99
	dates := []string{}
	avgPercentagePBs := []opts.LineData{}
	for _, dateAndAvgPercentagePB := range orderedDates {
		for date, avgPercentagePB := range dateAndAvgPercentagePB {
			if avgPercentagePB <= 0 {
				continue
			}
			if lowestAvgPb > avgPercentagePB {
				lowestAvgPb = avgPercentagePB
			}
			dates = append(dates, SimplifyDate(date))
			avgPercentagePBs = append(avgPercentagePBs, opts.LineData{
				Name:  SimplifyDate(date) + " " + strconv.Itoa(avgPercentagePB) + "%",
				Value: avgPercentagePB,
			})
		}
	}

	progress.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Experimental performance tracker",
			Subtitle: "Data points are average scores for every scenario played that day, converted into a percentage of your current highscore.",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "value",
			Max:       100,
			Min:       10 * (lowestAvgPb / 10),
			AxisLabel: &yAxisLabelFormatter,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
			Show:      true,
			Formatter: "{b}",
		}),
		ToolBoxOpts("performance"),
		xAxisOpts,
		initOpts,
	)

	progress.SetXAxis(dates).
		AddSeries("PB %", avgPercentagePBs).
		SetSeriesOptions(seriesOpts...)

	return renderToHTML(progress)
}

var max = "max"
var avg = "avg"

var seriesOpts = []charts.SeriesOpts{
	charts.WithLabelOpts(opts.Label{Show: true, Color: "black"}),
	charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
}

func titleOpts(scenName string, length int) charts.GlobalOpts {
	return charts.WithTitleOpts(opts.Title{
		Title:    scenName,
		Subtitle: fmt.Sprintf("Grouped by day, %v datapoints.", length),
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

func toolBoxFeatures(fileName string) *opts.ToolBoxFeature {
	return &opts.ToolBoxFeature{
		SaveAsImage: saveAsImage(fileName),
		DataZoom:    &dataZoom,
	}
}

// ToolBoxOpts ...
func ToolBoxOpts(fileName string) charts.GlobalOpts {
	return charts.WithToolboxOpts(opts.Toolbox{
		Show:    true,
		Feature: toolBoxFeatures(fileName),
	})
}

var tooltipOpts = charts.WithTooltipOpts(opts.Tooltip{
	Trigger:   "item",
	TriggerOn: "mousemove|click",
	Show:      true,
	Formatter: "{b}",
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
		Type:      "value",
		Max:       10 * ((int(highscore*1.05) + 9) / 10),
		Min:       10 * (int(lowestAvg*0.95) / 10),
		AxisLabel: &yAxisLabelFormatter,
	})
}

var legendOpts = charts.WithLegendOpts(opts.Legend{
	Show: true,
})

var initOpts = charts.WithInitializationOpts(opts.Initialization{
	AssetsHost: "static/",
})

func chartGlobalOpts(groupedBy string, scenName string, length int, hs float64, ls float64) (charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts, charts.GlobalOpts) {
	return titleOpts(scenName, length), ToolBoxOpts(scenName), tooltipOpts, xAxisOpts, yAxisOpts(hs, ls), legendOpts, initOpts
}
