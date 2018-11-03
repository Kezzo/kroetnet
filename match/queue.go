package match

import (
	"sync"
)

// PastState holds past state information
type PastState struct {
	Frame byte
	Xpos,
	Ypos int32
	Xtrans,
	Ytrans byte
}

// Queue is a FIFO structure
type Queue struct {
	nodes []*PastState
	size  int
	mutex *sync.Mutex
}

// Push adds a past state to the queue.
func (q *Queue) Push(n *PastState) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.nodes) < q.size {
		q.nodes = append(q.nodes, n)
	} else {
		q.nodes = q.nodes[1:]
		q.nodes = append(q.nodes, n)
	}
}

// Pop removes and returns a node from the queue
func (q *Queue) Pop() *PastState {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.nodes) == 0 {
		return nil
	}
	node := q.nodes[0]
	q.nodes = q.nodes[1:]
	return node
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*PastState, size),
		size:  size,
		mutex: &sync.Mutex{},
	}
}
