package main

import (
	"fmt"
	"kroetnet/msg"
	"testing"
)

func TestPlayerMove(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8}
	im := msg.InputMsg{MessageID: 0,
		PlayerID: 16, XTranslation: 200,
		YTranslation: 100, Rotation: 0, Frame: 10}
	x, y := p.move(im)
	if x != 156 {
		t.Errorf("Player Move returned a wrong result")
	}
	if y != -78 {
		t.Errorf("Player Move returned a wrong result")
	}
}

func TestPlayerValidateStates(t *testing.T) {
	p := Player{nil, 16, 0, 0, 8}
	imArr := []msg.InputMsg{
		{MessageID: 0, PlayerID: 16, XTranslation: 200, YTranslation: 100, Rotation: 0, Frame: 10},
		{MessageID: 0, PlayerID: 16, XTranslation: 200, YTranslation: 160, Rotation: 0, Frame: 11}}
	x, y := p.validateMoves(imArr)
	fmt.Println("Move result: ", x, " ", y)
	if x != 312 {
		t.Errorf("Player Move returned a wrong result")
	}
	if y != 47 {
		t.Errorf("Player Move returned a wrong result")
	}
}
