package main

func main() {
	playerCount := 1
	port := ":2448"
	game := newGame(playerCount, 5, port)
	game.startServer()
}
