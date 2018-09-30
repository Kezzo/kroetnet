package main

import (
	"kroetnet/msg"
	"net"
)

var xmax int32 = 2400
var ymax int32 = 2400
var unitSpeed float64 = 100

// Player details
type Player struct {
	ipAddr   net.Addr
	id       int
	X        int32
	Y        int32
	rotation byte
}

func (p Player) move(input msg.InputMsg) (int32, int32) {
	Xtrans, Ytrans := input.XTranslation, input.YTranslation
	resX, resY := p.X, p.Y
	stepSize := 1. / 128.
	// TODO check boundaries
	if Xtrans < 128 {
		resX += int32(-1 * unitSpeed * (stepSize * float64(Xtrans)))
	} else if Xtrans > 128 {
		resX += int32(unitSpeed * (stepSize * float64(Xtrans)))
	}

	if resX > xmax {
		resX = xmax
	} else if resX < -xmax {
		resX = -xmax
	}

	if Ytrans < 128 {
		resY += int32(-1 * unitSpeed * (stepSize * float64(Ytrans)))
	} else if Ytrans > 128 {
		resY += int32(unitSpeed * (stepSize * float64(Ytrans)))
	}

	if resY > ymax {
		resY = ymax
	} else if resY < -ymax {
		resY = -ymax
	}

	return resX, resY
}

func (p Player) validateMoves(buffer []msg.InputMsg) (int32, int32) {
	for _, v := range buffer {
		p.X, p.Y = p.move(v)
	}
	return p.X, p.Y
}
