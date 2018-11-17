package units

import "kroetnet/collision"

// Unit ...
type Unit interface {
	GetTeam() byte
	GetID() byte
	GetCollider() collision.Collider
	GetPosition() (int32, int32)
	GetRotation() byte
	AddDamage(damageToAdd int32)
	GetHealthPercent() int32
	IsPlayer() bool
}
