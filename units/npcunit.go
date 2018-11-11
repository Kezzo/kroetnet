package units

import "kroetnet/collision"

// NPCUnit ...
type NPCUnit struct {
	ID       byte
	Team     byte
	X        int32
	Y        int32
	Rotation byte
	Collider collision.Collider
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
