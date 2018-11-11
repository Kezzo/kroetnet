package collision

import "kroetnet/utility"

// CircleCollider ...
type CircleCollider struct {
	Xpos   int32
	Ypos   int32
	Radius int32
}

// NewCircleCollider ...
func NewCircleCollider(xPos int32, yPos int32, radius int32) *CircleCollider {
	return &CircleCollider{
		Xpos:   xPos,
		Ypos:   yPos,
		Radius: radius}
}

// Update ...
func (cc *CircleCollider) Update(xPos int32, yPos int32, rotation byte) {
	cc.Xpos = xPos
	cc.Ypos = yPos
}

// IsColliding ...
func (cc *CircleCollider) IsColliding(collider Collider) bool {
	otherCC, ok := collider.(*CircleCollider)

	if ok {
		return utility.GetDistance(cc.Xpos, cc.Ypos, otherCC.Xpos, otherCC.Ypos) <= otherCC.Radius+cc.Radius
	}

	return false
}
