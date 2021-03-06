package match

import (
	"testing"
)

func TestQInit(t *testing.T) {
	q := NewQueue(1)
	if len(q.nodes) != 1 {
		t.Errorf("Failed to create Queue with given size")
	}
}

func TestQPushPop(t *testing.T) {
	q := NewQueue(1)
	ps := &PastState{1, 120, 245, 100, 100}
	q.Push(ps)
	if q.Pop() != ps {
		t.Errorf("Failed to Pop pushed Element in Q")
	}
}

func TestQPushPeek(t *testing.T) {
	q := NewQueue(1)
	ps := &PastState{1, 120, 245, 100, 100}
	q.Push(ps)
	if q.nodes[0] != ps {
		t.Errorf("Failed to peek pushed Element in Q")
	}
}

func TestQPushPushPop(t *testing.T) {
	q := NewQueue(2)
	ps := &PastState{1, 120, 245, 123, 124}
	q.Push(ps)
	ps1 := &PastState{1, 20, 24, 213, 121}
	q.Push(ps1)
	if q.Pop() != ps {
		t.Errorf("Failed to Pop pushed Element in Q")
	}
	if q.Pop() != ps1 {
		t.Errorf("Failed to Pop pushed Element in Q")
	}
}
