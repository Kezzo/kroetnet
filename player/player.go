package player

import (
	"kroetnet/collision"
	"kroetnet/utility"
	"math"
	"net"
	"time"
)

var xmax float64 = 24000
var ymax float64 = 24000
var unitSpeed float64 = 400
var colliderSize int32 = 2000

// Player details
type Player struct {
	IPAddr   net.Addr
	ID       byte
	X        int32
	Y        int32
	Rotation byte
	LastMsg  time.Time
	Collider collision.Collider
}

// EmptyPlayer ...
var EmptyPlayer = Player{}

// NewPlayer ...
func NewPlayer(ID byte, xPos int32, yPos int32, ipAddr net.Addr) *Player {
	return &Player{
		IPAddr:   ipAddr,
		ID:       ID,
		X:        xPos,
		Y:        yPos,
		Collider: &collision.CircleCollider{Xpos: xPos, Ypos: yPos, Radius: colliderSize}}
}

// GetPosition ...
func GetPosition(xPos int32, yPos int32, xTranslation byte, yTranslation byte) (int32, int32) {
	resX, resY := xPos, yPos

	movX, movY := GetTranslation(xTranslation, yTranslation)

	resX = int32(math.Min(xmax, math.Max(-xmax, float64(resX+movX))))
	resY = int32(math.Min(ymax, math.Max(-ymax, float64(resY+movY))))

	return resX, resY
}

// GetTranslation ...
func GetTranslation(xTranslation byte, yTranslation byte) (int32, int32) {
	movX, movY := 0., 0.

	if xTranslation != 127 {
		movX = utility.Lerp(-1, 1, utility.InverseLerp(0, 255, float64(xTranslation)))
	}

	if yTranslation != 127 {
		movY = utility.Lerp(-1, 1, utility.InverseLerp(0, 255, float64(yTranslation)))
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

	return int32(math.Round(unitSpeed * movX)), int32(math.Round(unitSpeed * movY))
}
