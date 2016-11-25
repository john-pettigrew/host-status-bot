package main

import "log"

func main() {
	err := startLoop()
	if err != nil {
		log.Fatal(err)
	}
}
