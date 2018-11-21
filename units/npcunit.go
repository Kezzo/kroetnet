package units

import (
	"kroetnet/collision"
	"math"
)

// NPCUnit ...
type NPCUnit struct {
	ID            byte
	Team          byte
	UnitType      byte
	X             int32
	Y             int32
	Rotation      byte
	HealthPercent int32
	Collider      collision.Collider
}

// GetID ...
func (npc *NPCUnit) GetID() byte {
	return npc.ID
}

// GetTeam ...
func (npc *NPCUnit) GetTeam() byte {
	return npc.Team
}

// GetUnitType ...
func (npc *NPCUnit) GetUnitType() byte {
	return npc.UnitType
}

// GetCollider ...
func (npc *NPCUnit) GetCollider() collision.Collider {
	return npc.Collider
}

// GetPosition ...
func (npc *NPCUnit) GetPosition() (int32, int32) {
	return npc.X, npc.Y
}

// GetRotation ...
func (npc *NPCUnit) GetRotation() byte {
	return npc.Rotation
}

// AddDamage ...
func (npc *NPCUnit) AddDamage(damageToAdd int32) {
	newHealth := npc.HealthPercent - damageToAdd

	if newHealth <= 0 {
		newHealth = 0
		// dead
	}

	npc.HealthPercent = newHealth
}

// GetHealthPercent ...
func (npc *NPCUnit) GetHealthPercent() int32 {
	return npc.HealthPercent
}

// IsPlayer ...
func (npc *NPCUnit) IsPlayer() bool {
	return false
}

// SetPosition ...
func (npc *NPCUnit) SetPosition(xPos int32, yPos int32, xTranslation byte, yTranslation byte) {
	resX, resY := xPos, yPos

	movX, movY := getTranslation(xTranslation, yTranslation)

	npc.X = int32(math.Min(xmax, math.Max(-xmax, float64(resX+movX))))
	npc.Y = int32(math.Min(ymax, math.Max(-ymax, float64(resY+movY))))
}
