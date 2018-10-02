package main

import (
	"fmt"
	"kroetnet/msg"
	"testing"
)

func TestPlayerMoveX(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8}
	im := msg.InputMsg{MessageID: 0,
		PlayerID: 16, XTranslation: 255,
		YTranslation: 0, Frame: 10}
	x, y := p.move(im)
	if x != 100 {
		t.Errorf("Player Move returned a wrong x result")
	}
	if y != 0 {
		t.Errorf("Player Move returned a wrong y result")
	}
}

func TestPlayerMoveY(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8}
	im := msg.InputMsg{MessageID: 0,
		PlayerID: 16, XTranslation: 0,
		YTranslation: 255, Frame: 10}
	x, y := p.move(im)
	if x != 0 {
		t.Errorf("Player Move returned a wrong x result")
	}
	if y != 100 {
		t.Errorf("Player Move returned a wrong y result")
	}
}
func TestPlayerMoveNegativeY(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8}
	im := msg.InputMsg{MessageID: 0,
		PlayerID: 16, XTranslation: 0,
		YTranslation: 0, Frame: 10}
	x, y := p.move(im)
	if x != 0 {
		t.Errorf("Player Move returned a wrong x result")
	}
	if y != -100 {
		t.Errorf("Player Move returned a wrong y result")
	}
}

func TestPlayerValidateStates(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8}
	imArr := []msg.InputMsg{
		{MessageID: 0, PlayerID: 16, XTranslation: 255, YTranslation: 0,
			Rotation: 0, Frame: 10},
		{MessageID: 0, PlayerID: 16, XTranslation: 0, YTranslation: 255,
			Rotation: 0, Frame: 11}}
	x, y := p.validateMoves(imArr)
	fmt.Println("MOVES ", x, y)
	if x != 100 {
		t.Errorf("Player Move returned a wrong result")
	}
	if y != 100 {
		t.Errorf("Player Move returned a wrong result")
	}
}
