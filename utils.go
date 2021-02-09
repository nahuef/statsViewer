package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// GetStatsPath ...
func GetStatsPath() (path string, errStr string) {
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

// SimplifyDate ...
func SimplifyDate(d string) string {
	d = strings.ReplaceAll(d, "2018", "18")
	d = strings.ReplaceAll(d, "2019", "19")
	d = strings.ReplaceAll(d, "2020", "20")
	d = strings.ReplaceAll(d, "2021.", "")
	d = strings.ReplaceAll(d, ".", "/")

	return d
}
