package main

import "log"

// Check ...
func Check(e error) {
	if e != nil {
		log.Println(e)
		EnterToExit()
	}
}
