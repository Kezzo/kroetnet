package main

func main() {
	playerCount := 1
	port := ":2448"
	game := newGame(playerCount, 15, port)
	game.startServer()
}
