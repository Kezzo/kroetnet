package units

import (
	"testing"
	"time"
)

func TestPlayerMove(t *testing.T) {
	p := Player{0, 0, 0, 0, 8, 100, 0, nil, nil, time.Now()}

	x, y := GetPlayerPosition(p.X, p.Y, 255, 127)
	if x != 250 {
		t.Errorf("Player Move returned a wrong x result")
	}

	x, y = GetPlayerPosition(p.X, p.Y, 127, 255)
	if y != 250 {
		t.Errorf("Player Move returned a wrong y result")
	}
}
func TestPlayerMoveNegative(t *testing.T) {
	p := Player{0, 0, 0, 0, 8, 100, 0, nil, nil, time.Now()}

	x, y := GetPlayerPosition(p.X, p.Y, 0, 127)
	if x != -250 {
		t.Errorf("Player Move returned a wrong x result")
	}

	x, y = GetPlayerPosition(p.X, p.Y, 127, 0)
	if y != -250 {
		t.Errorf("Player Move returned a wrong y result")
	}
}
