package main

func main() {
	playerCount := 2
	port := ":2448"
	game := newGame(playerCount, 15, port)
	game.startServer()
}
