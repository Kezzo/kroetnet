package collision

// Collider ...
type Collider interface {
	Update(xPos int32, yPos int32, rotation byte)
	IsColliding(collider Collider) bool
}
