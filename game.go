package main

import "kroetnet/msg"

// Game holds current game properties
type Game struct {
	players              []Player
	State                int
	StateChangeTimestamp int64
	start                int
	end                  int64
	playerBuffers        map[int][]msg.InputMsg
}
