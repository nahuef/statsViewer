package main

import (
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

// StatsPath ...
var StatsPath, statsNotFound = GetStatsPath()

func main() {
	defer func() {
		if error := recover(); error != nil {
			fmt.Println("Error:", error)
			EnterToExit()
		}
	}()

	start := time.Now()

	files, err := ioutil.ReadDir(StatsPath)
	if err != nil {
		if StatsPath == DefaultPath {
			cwd, err := os.Getwd()
			Check(err)
			StatsPath = "current working directory " + cwd
		}
		log.Printf("Error: %v\"stats\" folder not found, make sure the path is right %v", statsNotFound, StatsPath)
		EnterToExit()
	}
	fmt.Println("\"stats\" folder found!\nParsing files... This may take a few minutes!")

	stats := ParseStats(files)
	fmt.Println("Files parsed. Creating HTML file...")

	// Output HTML
	t, err := template.ParseFiles("static/statsViewerTpl.html")
	Check(err)
	f, err := os.Create(statsViewerHTML)
	Check(err)
	err = t.Execute(f, stats)
	Check(err)
	f.Close()

	fmt.Println("Success!")
	fmt.Println(time.Now().Sub(start))
	exec.Command("cmd", "/C", "start", statsViewerHTML).Run()
}
