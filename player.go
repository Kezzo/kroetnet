package main

import (
	"kroetnet/msg"
	"math"
	"net"
)

var xmax float64 = 24000
var ymax float64 = 24000
var unitSpeed float64 = 400

// Player details
type Player struct {
	ipAddr   net.Addr
	id       int
	X        int32
	Y        int32
	rotation byte
}

func (p *Player) move(input msg.InputMsg) (int32, int32) {
	Xtrans, Ytrans := input.XTranslation, input.YTranslation
	resX, resY := p.X, p.Y

	movX, movY := 0., 0.

	if Xtrans != 127 {
		movX = lerp(-1, 1, inverseLerp(0, 255, float64(Xtrans)))
	}

	if Ytrans != 127 {
		movY = lerp(-1, 1, inverseLerp(0, 255, float64(Ytrans)))
	}

	combinedMov := math.Abs(movX) + math.Abs(movY)

	// restrict diagonnal movement so it doesn't feels much faster than straight movement
	// not restricting to combined max translation of 1, because then diagonal movement feels too slow
	if combinedMov > 1.5 {
		overhead := math.Abs(1.5-combinedMov) / 2

		if movX > 0 {
			movX = movX - overhead
		} else {
			movX = movX + overhead
		}

		if movY > 0 {
			movY = movY - overhead
		} else {
			movY = movY + overhead
		}
	}

	resX += int32(math.Round(unitSpeed * movX))
	resX = int32(math.Min(xmax, math.Max(-xmax, float64(resX))))

	resY += int32(math.Round(unitSpeed * movY))
	resY = int32(math.Min(ymax, math.Max(-ymax, float64(resY))))

	return resX, resY
}

func (p *Player) validateMoves(buffer []msg.InputMsg) (int32, int32) {
	for _, v := range buffer {
		p.X, p.Y = p.move(v)
	}
	return p.X, p.Y
}

// accepts a min value and max value and a value between 0, 1. Will return the linearly interpolated value between the two first given values.
// i.e. (0, 255, 0.5) = 127
func lerp(from float64, to float64, value float64) float64 {
	return from*(1-value) + to*value
}

// similiar to lerp, but accepts a value that lies between the first two values and returns a value between 0 and 1 depending on where the values lies
// i.e. (0, 255, 127) = 0.5
func inverseLerp(from float64, to float64, value float64) float64 {
	return (value - from) / (to - from)
}
