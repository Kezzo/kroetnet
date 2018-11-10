package collision

import (
	"kroetnet/utility"
	"math"
)

// ConeCollider ...
type ConeCollider struct {
	Xpos            int32
	Ypos            int32
	length          int32
	rotationDegress float64
	widthDegrees    float64
}

// Update ...
func (cc *ConeCollider) Update(xPos int32, yPos int32, rotation byte) {
	cc.Xpos = xPos
	cc.Xpos = yPos
	cc.rotationDegress = utility.Lerp(0, 360, utility.InverseLerp(0, 255, float64(rotation)))
}

// IsColliding ...
func (cc *ConeCollider) IsColliding(collider Collider) bool {

	otherCC, ok := collider.(*CircleCollider)

	if ok {
		distanceToOther := utility.GetDistance(cc.Xpos, cc.Ypos, otherCC.Xpos, otherCC.Ypos) - otherCC.Radius

		// too far away
		if distanceToOther > cc.length {
			return false
		}

		// check angle
		rightAnglePositionX, rightAnglePositionY := float64(cc.Xpos), float64(cc.Ypos+distanceToOther)
		rotationRadian := utility.DegreesToRadian(cc.rotationDegress)

		rightAnglePositionXWithRot := int32(rightAnglePositionX*math.Cos(rotationRadian) + rightAnglePositionY*math.Sin(rotationRadian))
		rightAnglePositionYWithRot := int32(rightAnglePositionY*math.Cos(rotationRadian) - rightAnglePositionX*math.Sin(rotationRadian))
		distanceFromConeCenterToOther := utility.GetDistanceInFloat(rightAnglePositionXWithRot, rightAnglePositionYWithRot, otherCC.Xpos, otherCC.Ypos) - float64(otherCC.Radius)

		distanceToOtherFloat := float64(distanceToOther)
		// Law of Cosines
		angleToOther := math.Acos(
			((distanceToOtherFloat * distanceToOtherFloat) +
				(distanceToOtherFloat * distanceToOtherFloat) -
				(distanceFromConeCenterToOther * distanceFromConeCenterToOther)) / (2 * distanceToOtherFloat * distanceToOtherFloat))

		if math.IsNaN(angleToOther) {
			return false
		}

		angleToOther = utility.RadianToDegress(angleToOther)

		// a 90 degress code has a width degress of 45, because we check the angle differences to the center of the cone
		return angleToOther <= cc.widthDegrees
	}

	return false
}
