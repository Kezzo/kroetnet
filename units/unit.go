package units

import "kroetnet/collision"

// Unit ...
type Unit interface {
	GetTeam() byte
	GetCollider() collision.Collider
	GetPosition() (int32, int32)
	GetRotation() byte
}
