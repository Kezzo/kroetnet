package main

import (
	"log"
	"net"
	"time"
)

// Game holds current game properties
type Game struct {
	players              []Player
	State                int
	Frame                byte
	StateChangeTimestamp int64
	start                time.Time
	end                  time.Time
	statesMap            []Queue
}

// PastState holds past state information
type PastState struct {
	Frame byte
	Xpos,
	Ypos int32
	Xtans,
	Ytrans byte
}

// Queue is a FIFO structure
type Queue struct {
	nodes []*PastState
	size  int
	head  int
	tail  int
	count int
}

// Push adds a past state to the queue.
func (q *Queue) Push(n *PastState) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*PastState, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	// log.Println("HEAD ", q.head, " Tail ", q.tail, "LEN ", len(q.nodes))
	q.count++
}

// Pop removes and returns a node from the queue
func (q *Queue) Pop() *PastState {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*PastState, size),
		size:  size,
	}
}

var emptyPlayer = Player{}

// AddPlayer adds servers to the game
func AddPlayer(addr net.Addr) int {
	playerID := -1
	nextPlayerID := 0
	for i := 0; i < len(game.players); i++ {
		if game.players[i] != emptyPlayer {
			nextPlayerID = i + 1
			// find player with same addr from udp packet
			if game.players[i].ipAddr == addr {
				playerID = game.players[i].id
				break
			}
		}
	}
	// player no in match yet & match not full
	if playerID == -1 && game.players[len(game.players)-1] == emptyPlayer {
		game.players[nextPlayerID] = Player{id: nextPlayerID, ipAddr: addr}
		playerID = nextPlayerID
		game.statesMap[playerID] = *NewQueue(15)
		return playerID
	}
	return -1
}

// CheckGameFull changes the gamestate when all players joined
func CheckGameFull(pc net.PacketConn, addr net.Addr) {
	// last player joined and match is full
	if game.players[len(game.players)-1] != emptyPlayer {
		// wait for all players
		for _, v := range game.players {
			sendGameStart(pc, v.ipAddr)
		}
		// skip state 1 for now
		time.Sleep(time.Second)
		game.State = 2
		// frame tick every 33 ms
		go doEvery(33*time.Millisecond, incFrame)
		log.Println("GAME STARTED")

		game.start = time.Now()
		// game ends after 2 minutes
		game.end = time.Now().Add(time.Minute * 1)
		// game.State = 1
		// game.StateChangeTimestamp = time.Now().Unix()
	}

}
