package main

import (
	"strconv"
	"strings"
)

// Extract ...
type Extract struct {
	line      string
	fileName  string
	challenge *Challenge
}

func (e *Extract) extractData() {
	e.name()
	e.score()
	e.sensScale()
	e.hsens()
	e.vsens()
	e.fov()
	e.DateAndTime()
}

func (e *Extract) name() {
	if strings.Contains(e.line, "Scenario:,") {
		e.challenge.Name = strings.Split(e.line, ",")[1]
	}
}

func (e *Extract) sensScale() {
	if strings.Contains(e.line, "Sens Scale:,") {
		e.challenge.SensScale = strings.Split(e.line, ",")[1]
	}
}

func (e *Extract) score() {
	if strings.Contains(e.line, "Score:,") {
		scoreStr := strings.Split(e.line, ",")[1]
		scoreFloat, _ := strconv.ParseFloat(scoreStr, 1)
		e.challenge.Score = float64(int(scoreFloat*10)) / 10
	}
}

func (e *Extract) hsens() {
	if strings.Contains(e.line, "Horiz Sens:,") {
		hsensStr := strings.Split(e.line, ",")[1]
		hsensFloat, _ := strconv.ParseFloat(hsensStr, 1)
		e.challenge.HSens = float64(int(hsensFloat*10)) / 10
	}
}

func (e *Extract) vsens() {
	if strings.Contains(e.line, "Vert Sens:,") {
		vsensStr := strings.Split(e.line, ",")[1]
		vsensFloat, _ := strconv.ParseFloat(vsensStr, 1)
		e.challenge.VSens = float64(int(vsensFloat*10)) / 10
	}
}

func (e *Extract) fov() {
	if strings.Contains(e.line, "FOV:,") {
		fovStr := strings.Split(e.line, ",")[1]
		fovFloat, _ := strconv.ParseFloat(fovStr, 1)
		e.challenge.FOV = float64(int(fovFloat*10)) / 10
	}
}

// DateAndTime ...
func (e *Extract) DateAndTime() {
	datetimeAndExtension := (strings.Split(e.fileName, " - Challenge - "))[1]
	datetime := strings.Split(datetimeAndExtension, " ")[0]
	dateAndTime := strings.Split(datetime, "-")

	e.challenge.Datetime = datetime
	e.challenge.Date = dateAndTime[0]
	e.challenge.Time = dateAndTime[1]
}
