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
	game := newGame(playerCount, 15, port)
	game.startServer()
}
