package main

import (
	"kroetnet/msg"
	"math"
	"net"
)

var xmax int32 = 24000
var ymax int32 = 24000
var unitSpeed float64 = 400

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
	if Xtrans < 127 {
		movement := math.Max(stepSize*float64(Xtrans), 1.)
		summand := int32(-1 * unitSpeed * movement)
		resX += summand
	} else if Xtrans > 127 {
		movement := math.Round(stepSize * float64(Xtrans/2.))
		summand := int32(unitSpeed * movement)
		resX += summand
	}

	if resX > xmax {
		resX = xmax
	} else if resX < -xmax {
		resX = -xmax
	}

	if Ytrans < 127 {
		movement := math.Max(stepSize*float64(Ytrans), 1.)
		summand := int32(-1 * unitSpeed * movement)
		resY += summand
	} else if Ytrans > 127 {
		movement := math.Round(stepSize * float64(Ytrans/2.))
		summand := int32(unitSpeed * movement)
		resY += summand
	}

	if resY > ymax {
		resY = ymax
	} else if resY < -ymax {
		resY = -ymax
	}

	// p.X, p.Y = resX, resY

	return resX, resY
}

func (p Player) validateMoves(buffer []msg.InputMsg) (int32, int32) {
	for _, v := range buffer {
		p.X, p.Y = p.move(v)
	}
	return p.X, p.Y
}
