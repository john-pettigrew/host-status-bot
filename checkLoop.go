package main

import (
	"fmt"
	"time"
)

func startLoop() {
	for {
		go loop()
		time.Sleep(time.Second)
	}
}

func loop() {
	fmt.Println("loop")
	// Get sites from DB

	// Check each site

}
