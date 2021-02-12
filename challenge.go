package main

import (
	"fmt"
)

// Challenge contains data from a single CSV file.
type Challenge struct {
	Name string
	// Datetime  string
	Date string
	// Time      string
	Score     float64
	SensScale string
	HSens     float64
	VSens     float64
	FOV       float64
}

// SensStr formats the sensibility output checking if vertical and horizontal sesibilities are different or not.
func (c *Challenge) SensStr() string {
	if c.HSens == c.VSens {
		return fmt.Sprintf("Sens: %v %v", c.HSens, c.SensScale)
	}

	return fmt.Sprintf("Vsens: %v, Hsens: %v %v", c.VSens, c.HSens, c.SensScale)
}

// IsValid checks if challenge fields are valid.
func (c *Challenge) IsValid() bool {
	if c.Name == "" || c.Date == "" || c.SensScale == "" {
		return false
	}
	if c.Score == 0 || c.HSens <= 0 || c.VSens <= 0 || c.FOV <= 0 {
		return false
	}

	return true
}
