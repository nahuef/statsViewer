package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testScenName = "TEST Thin Gauntlet"
var date = "2020.11.13"
var timeScen = "22.08.18"
var datetime = fmt.Sprintf("%s-%s", date, timeScen)
var testFileName = fmt.Sprintf("%s - Challenge - %s Stats.csv", testScenName, datetime)
var testFilePath = StatsPath + testFileName

func TestStatsParse(t *testing.T) {
	// Create test file
	ioutil.WriteFile(testFilePath, []byte{}, 0777)
	defer os.Remove(testFilePath)
	// Add data
	file, error := os.OpenFile(testFilePath, os.O_WRONLY|os.O_APPEND, 0777)
	Check(error)
	defer file.Close()
	_, err := file.WriteString(`
	Kills:,99
	Deaths:,0
	Fight Time:,3.109
	Avg TTK:,0.606
	Damage Done:,2475.0
	Damage Taken:,0.0
	Midairs:,0
	Midaired:,0
	Directs:,0
	Directed:,0
	Distance Traveled:,0.0
	Score:,915.981323
	Scenario:,TEST Thin Gauntlet
	Hash:,a5a9fea2fb99346f14d742b8d19777b3
	Game Version:,2.0.2.0

	Input Lag:,0
	Max FPS (config):,144.0
	Sens Scale:,Valorant
	Horiz Sens:,0.31
	Vert Sens:,0.31
	FOV:,103.0
	Hide Gun:,false
	Crosshair:,plus.png
	Crosshair Scale:,1.0
	Crosshair Color:,FFFF00
	`)
	Check(err)

	fi, err := os.Stat(testFilePath)
	Check(err)

	files := []os.FileInfo{fi}
	stats := ParseStats(files)

	// fmt.Printf("%+v", stats.Scenarios[testScenName])
	expected := Scenario{
		fileName:    testFileName,
		Name:        testScenName,
		TimesPlayed: 1,
		Challenges: []Challenge{
			{
				Name: testScenName,
				// Datetime:  "2020.11.13-22.08.18",
				Date: "2020.11.13",
				// Time:      "22.08.18",
				Score:     915.9,
				SensScale: "Valorant",
				HSens:     0.3,
				VSens:     0.3,
				FOV:       103,
			},
		},
		Highscore:      915.9,
		Lowscore:       915.9,
		LowestAvgScore: 915.9,
		ByDateMax: []map[string]Challenge{
			{
				"2020.11.13": Challenge{
					Name:      "TEST Thin Gauntlet",
					Date:      "2020.11.13",
					Score:     915.9,
					SensScale: "Valorant",
					HSens:     0.3,
					VSens:     0.3,
					FOV:       103,
				},
			},
		},
		ByDateAvg: []map[string]DateAvg{
			{
				"2020.11.13": DateAvg{
					Score:        915.9,
					Grouped:      1,
					PercentagePB: 100,
				},
			},
		},
		ByDateScores: []map[string][]float64{
			{
				"2020.11.13": []float64{
					915.9,
				},
			},
		},
		ByDateWMA: []map[string]DateWMA{
			{
				"2020.11.13": DateWMA{
					Avg:     915.9,
					Grouped: 1,
				},
			},
		},
	}

	reflect.DeepEqual(&expected, stats.Scenarios[testScenName])
	fmt.Println(reflect.DeepEqual(&expected, stats.Scenarios[testScenName]))
	assert.Equal(t, &expected, stats.Scenarios[testScenName], "Output Scenario should equal to expected")
}

func TestExtractDate(t *testing.T) {
	challenge := Challenge{}
	extractor := Extract{
		fileName:  testFileName,
		challenge: &challenge,
	}
	extractor.DateAndTime()

	// assert.Equal(t, challenge.Datetime, "2020.11.13-22.08.18", "Datetime should be equal")
	assert.Equal(t, challenge.Date, "2020.11.13", "Date should be equal")
	// assert.Equal(t, challenge.Time, "22.08.18", "Time should be equal")
}
