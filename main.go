package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

func main() {
	//go:embed templates/*
	var templates embed.FS

	files, err := ioutil.ReadDir("./stats")
	if err != nil {
		log.Fatal(err)
	}

	var statsCount = make(map[string]int)

	type statCount struct {
		Name  string
		Count int
	}

	for _, file := range files {
		if file.IsDir() == true {
			continue
		}

		scenName := (strings.Split(file.Name(), " - Challenge - "))[0]

		if _, ok := statsCount[scenName]; ok {
			statsCount[scenName]++
		} else {
			statsCount[scenName] = 1
		}
	}

	var sorted []statCount
	for name, count := range statsCount {
		sorted = append(sorted, statCount{name, count})
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	t, err := template.ParseFS(templates, "templates/statsViewer.html")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("statsViewer.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(f, sorted)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	fileName := "statsViewer.txt"
	err = ioutil.WriteFile(fileName, []byte(""), 0777)
	if err != nil {
		log.Fatal(err)
	}

	file, error := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0777)
	if error != nil {
		log.Fatal(error)
	}
	defer file.Close()

	for _, value := range sorted {
		line := fmt.Sprintln(value.Name+":", value.Count)
		fmt.Print(line)

		if _, err := file.WriteString(line); err != nil {
			log.Fatal(err)
		}
	}
}
