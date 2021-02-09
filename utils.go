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
