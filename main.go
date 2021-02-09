package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"
)

var statsViewer = "StatsViewer"
var statsViewerHTML = statsViewer + ".html"
var defaultPath = "./stats"

// StatsPath ...
var StatsPath, statsNotFound = GetStatsPath()

func main() {
	start := time.Now()

	files, err := ioutil.ReadDir(StatsPath)
	if err != nil {
		log.Println(statsNotFound + " \"stats\" folder not found.\nPress \"enter\" key to exit.")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}
	fmt.Println("\"stats\" folder found!\nParsing files... \nThis may take a few minutes!")

	stats := ParseStats(files)
	fmt.Println("Files parsed. \nCreating HTML file...")

	// Output HTML
	t, err := template.ParseFiles("templates/" + statsViewerHTML)
	Check(err)
	f, err := os.Create(statsViewerHTML)
	Check(err)
	err = t.Execute(f, stats)
	Check(err)
	f.Close()

	// Output TXT
	var fileName = statsViewer + ".txt"
	err = ioutil.WriteFile(fileName, []byte(""), 0777)
	Check(err)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0777)
	Check(err)
	defer file.Close()
	for _, scenario := range stats.SortedTimesPlayed {
		line := fmt.Sprintln(scenario.Name+":", scenario.TimesPlayed)

		_, err := file.WriteString(line)
		Check(err)
	}

	fmt.Println("Success!")
	exec.Command("cmd", "/C", "start", statsViewerHTML).Run()

	time2 := time.Now()
	fmt.Println(time2.Sub(start))
}
