package main

// Game holds current game properties
type Game struct {
	players              []Player
	State                int
	Frame                byte
	StateChangeTimestamp int64
	start                int
	end                  int64
	statesMap            map[int]Queue
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
