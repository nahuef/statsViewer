package main

import (
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func getWords(scens *map[string]*Scenario) []string {
	var b strings.Builder
	for name := range *scens {
		fmt.Fprintf(&b, "%v ", name)
	}

	return strings.Split(b.String(), " ")
}

func filterWords(words []string) []string {
	remove := []string{"10", "but", "No", "no", "NO", "has", "fs", "and", "AND", "gp9", "v1", "V1", "v2", "V2", "v3",
		"V3", "v4", "V4", "v5", "V5", "II", "ii", "60s", "+2", "hs", "HS", "66", "I", "i", "OUT", "GP9", "_", "l", "L",
		"Your", "your", "-", "|", "||", "|||", "III", "||||", "the", "THE", "The", "off", "bots", "Bots", "by", "49",
		"LR", "Wrist)", "55", "49", "15", "65", "37", "28", "(Smooth", "12", "57", "is", "ar", "es", "cata", "very",
		"acc", "beta", "01", "13", "62", "5s", "VL", "30", "59", "1T", "50", "SF"}

	result := []string{}
	for i := len(words) - 1; i >= 0; i-- {
		trimmedW := strings.Trim(words[i], "\t \n")

		for _, rem := range remove {
			if len(trimmedW) == 1 || trimmedW == rem {
				result = append(words[:i], words[i+1:]...)
				break
			}
		}
	}

	return result
}

func wordOcurrences(words []string) map[string]int {
	m := map[string]int{}
	for _, word := range words {
		m[word] = m[word] + 1
	}

	return m
}

func sortAndRemoveWords(weightedWords map[string]int) []map[string]int {
	sorted := make([]map[string]int, len(weightedWords))

	for word, count := range weightedWords {
		sorted = append(sorted, map[string]int{word: count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		var iV, jV int
		for _, v := range sorted[i] {
			iV = v
		}
		for _, v := range sorted[j] {
			jV = v
		}

		return iV > jV
	})
	toTrim := float64(len(sorted)) * float64(0.5)
	trimmed := sorted[:int(toTrim)]

	return trimmed
}

func generateWCData(data []map[string]int) []opts.WordCloudData {
	items := make([]opts.WordCloudData, 0)
	for _, kv := range data {
		for k, v := range kv {
			items = append(items, opts.WordCloudData{Name: k, Value: v})
		}
	}

	return items
}

// WordCloud ...
func WordCloud(scens *map[string]*Scenario) template.HTML {
	wc := charts.NewWordCloud()
	wc.Renderer = newSnippetRenderer(wc, wc.Validate)
	wc.SetGlobalOptions(ToolBoxOpts("wordcloud"))

	words := getWords(scens)
	filteredWords := filterWords(words)
	weightedWords := wordOcurrences(filteredWords)
	trimmed := sortAndRemoveWords(weightedWords)

	wc.AddSeries("Word Cloud", generateWCData(trimmed)).
		SetSeriesOptions(charts.WithWorldCloudChartOpts(
			opts.WordCloudChart{
				SizeRange: []float32{12, 80},
				Shape:     "circle",
			}))

	wordcloud := renderToHTML(wc)
	return wordcloud
}
