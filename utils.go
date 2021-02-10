package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DefaultPath is a "stats" folder at binary's directory.
var DefaultPath = "./stats/"

// GetStatsPath ...
func GetStatsPath() (statsPath string, errStr string) {
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		errStr += "No config.json file found."
	} else {
		var parsedConfig map[string]interface{}
		err = json.Unmarshal(config, &parsedConfig)
		if err != nil {
			errStr = "Error reading config file. Make sure it is a valid JSON."
		} else {
			statsPath = filepath.Clean(parsedConfig["stats_path"].(string)) + "\\"
		}
	}

	if statsPath == "" {
		statsPath = DefaultPath
	}

	return statsPath, errStr
}

// SimplifyDate ...
func SimplifyDate(d string) string {
	sep := "/"
	d = strings.ReplaceAll(d, ".", sep)
	d = reorderDate(d, sep)
	d = strings.ReplaceAll(d, "2018", "18")
	d = strings.ReplaceAll(d, "2019", "19")
	d = strings.ReplaceAll(d, "2020", "20")
	d = strings.ReplaceAll(d, "/2021", "")

	return d
}

func reorderDate(d, sep string) string {
	dateUnits := strings.Split(d, sep)
	dateUnits[0], dateUnits[1], dateUnits[2] = dateUnits[1], dateUnits[2], dateUnits[0]

	return strings.Join(dateUnits, sep)
}

// ContainsString checks if the slice has the contains value in it.
// Source: https://github.com/stretchr/stew/blob/master/slice/contains.go
func ContainsString(slice []string, contains string) bool {
	for _, value := range slice {
		if value == contains {
			return true
		}
	}
	return false
}

// EnterToExit ...
func EnterToExit() {
	log.Println("Press \"enter\" key to exit.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(1)
}
