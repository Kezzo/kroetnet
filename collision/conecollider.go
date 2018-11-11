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

// NewConeCollider ...
func NewConeCollider(xPos int32, yPos int32, rotation byte, length int32, widthDegrees float64) *ConeCollider {
	return &ConeCollider{
		Xpos:            xPos,
		Ypos:            yPos,
		length:          length,
		rotationDegress: utility.Lerp(0, 360, utility.InverseLerp(0, 255, float64(rotation))),
		widthDegrees:    widthDegrees}
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
		distanceToOther := utility.GetDistance(cc.Xpos, cc.Ypos, otherCC.Xpos, otherCC.Ypos)

		// too far away
		if (distanceToOther - otherCC.Radius) > cc.length {
			return false
		}

		// check angle
		rotationRadian := utility.DegreesToRadian(cc.rotationDegress)

		// first find rotate end point of collider based on distance to target at '0,0'
		colliderDirectionX, colliderDirectionY := float64(0), float64(distanceToOther)
		colliderDirectionXWithRot := int32(colliderDirectionX*math.Cos(rotationRadian) + colliderDirectionY*math.Sin(rotationRadian))
		colliderDirectionYWithRot := int32(colliderDirectionY*math.Cos(rotationRadian) - colliderDirectionX*math.Sin(rotationRadian))

		// then shift that to actual collider position
		colliderPositionDirectionXWithRot := cc.Xpos + colliderDirectionXWithRot
		colliderPositionDirectionYWithRot := cc.Ypos + colliderDirectionYWithRot

		// get distance to target to get length of third trianle leg the collider, that hit range end point and the target create
		distanceFromConeCenterToOther := math.Max(0, utility.GetDistanceInFloat(colliderPositionDirectionXWithRot, colliderPositionDirectionYWithRot, otherCC.Xpos, otherCC.Ypos)-float64(otherCC.Radius))

		distanceToOtherFloat := float64(distanceToOther)

		// Law of Cosines
		// find angle of the leg that goes to the target and the one that goes in the direction of the collider
		angleToOther := math.Acos(
			((distanceToOtherFloat * distanceToOtherFloat) +
				(distanceToOtherFloat * distanceToOtherFloat) -
				(distanceFromConeCenterToOther * distanceFromConeCenterToOther)) / (2 * distanceToOtherFloat * distanceToOtherFloat))

		if math.IsNaN(angleToOther) {
			return false
		}

		angleToOther = utility.RadianToDegress(angleToOther)

		//log.Println("Conecollider is angle to target with degress of: ", angleToOther)

		// a 90 degress code has a width degress of 45, because we check the angle differences to the center of the cone
		return angleToOther <= cc.widthDegrees
	}

	return false
}
