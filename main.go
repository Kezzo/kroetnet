package main

import (
	"os"
)

func main() {
	playerCount := 2
	port := ":2448"
	if os.Getenv("GO_ENV") == "DEV" {
		port = ":0"
	}
	match := newMatch(playerCount, 15, port)
	match.startServer()
}
