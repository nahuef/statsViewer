package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var statsViewer = "statsViewer"
var statsViewerHTML = statsViewer + ".html"
var defaultPath = "./stats"

// StatsPath ...
var StatsPath, statsNotFound = getStatsPath()

func main() {
	files, err := ioutil.ReadDir(StatsPath)
	if err != nil {
		log.Println(statsNotFound + " \"stats\" folder not found. \n Press \"enter\" key to exit.")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}

	stats := ParseStats(files)

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
	for _, scenario := range stats.Sorted {
		line := fmt.Sprintln(scenario.Name+":", scenario.TimesPlayed)

		_, err := file.WriteString(line)
		Check(err)
	}

	fmt.Println("Success!")
	exec.Command("cmd", "/C", "start", statsViewerHTML).Run()
}

func getStatsPath() (path string, errStr string) {
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		errStr += "No config.json file found."
	} else {
		var parsedConfig map[string]interface{}
		err = json.Unmarshal(config, &parsedConfig)
		if err != nil {
			errStr = "Error reading config.json."
		} else {
			path = parsedConfig["stats_path"].(string) + "\\"
		}
	}

	if path == "" {
		path = defaultPath
	}

	return path, errStr
}
