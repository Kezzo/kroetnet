package main

// Game holds current game properties
type Game struct {
	players              []Player
	State                int
	StateChangeTimestamp int64
	start                int
	end                  int64
}
