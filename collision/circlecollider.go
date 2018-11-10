package collision

import "kroetnet/utility"

// CircleCollider ...
type CircleCollider struct {
	Xpos   int32
	Ypos   int32
	Radius int32
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
