package units

import "kroetnet/collision"

// NPCUnit ...
type NPCUnit struct {
	ID            byte
	Team          byte
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
