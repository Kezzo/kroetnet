package units

import "kroetnet/collision"

// Unit ...
type Unit interface {
	GetID() byte
	GetTeam() byte
	GetUnitType() byte
	GetCollider() collision.Collider
	GetPosition() (int32, int32)
	GetRotation() byte
	AddDamage(damageToAdd int32)
	AddHeal(healToAdd int32)
	GetHealthPercent() int32
	IsPlayer() bool
	SetPosition(xPos int32, yPos int32, xTranslation byte, yTranslation byte)
}
