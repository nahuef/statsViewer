package main

import "strings"

// SimplifyDate ...
func SimplifyDate(d string) string {
	d = strings.ReplaceAll(d, "2018", "18")
	d = strings.ReplaceAll(d, "2019", "19")
	d = strings.ReplaceAll(d, "2020", "20")
	d = strings.ReplaceAll(d, "2021.", "")
	d = strings.ReplaceAll(d, ".", "/")

	return d
}
