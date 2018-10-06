package main

func main() {
	playerCount := 1
	game := newGame(playerCount, 5, ":2448")
	game.startServer()
}
