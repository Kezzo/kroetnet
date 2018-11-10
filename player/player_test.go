package player

import (
	"testing"
	"time"
)

func TestPlayerMove(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8, time.Now()}

	x, y := GetPosition(p.X, p.Y, 255, 127)
	if x != 400 {
		t.Errorf("Player Move returned a wrong x result")
	}

	x, y = GetPosition(p.X, p.Y, 127, 255)
	if y != 400 {
		t.Errorf("Player Move returned a wrong y result")
	}
}
func TestPlayerMoveNegative(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8, time.Now()}

	x, y := GetPosition(p.X, p.Y, 0, 127)
	if x != -400 {
		t.Errorf("Player Move returned a wrong x result")
	}

	x, y = GetPosition(p.X, p.Y, 127, 0)
	if y != -400 {
		t.Errorf("Player Move returned a wrong y result")
	}
}
