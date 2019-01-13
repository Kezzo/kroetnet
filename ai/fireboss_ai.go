package ai

import (
	"kroetnet/units"
	"kroetnet/utility"
	"math"
)

// FireBoss ...
type FireBoss struct {
	unitData     units.Unit
	playerTarget *Target
}

// Target ...
type Target struct {
	unit     units.Unit
	distance float64
}

// Init ...
func (ai *FireBoss) Init(unit units.Unit) {
	ai.unitData = unit
}

// Tick ...
func (ai *FireBoss) Tick(unitsArr []units.Unit, updatedUnitIDs map[byte]bool) map[byte]bool {

	// unit is dead
	if ai.unitData.GetHealthPercent() == 0 {
		return updatedUnitIDs
	}

	possiblePlayerTargets := []Target{}

	for _, unitData := range unitsArr {
		if unitData.IsPlayer() == true {
			aiX, aiY := ai.unitData.GetPosition()
			uX, uY := unitData.GetPosition()
			dist := math.Sqrt(math.Pow(float64(aiX-uX), 2) + math.Pow(float64(aiY-uY), 2))
			possiblePlayerTargets = append(possiblePlayerTargets, Target{unitData, dist})
		}
	}

	min := possiblePlayerTargets[0].distance
	ai.playerTarget = &possiblePlayerTargets[0]
	for _, target := range possiblePlayerTargets {
		if target.distance < min {
			ai.playerTarget = &Target{target.unit, target.distance}
		}
	}
	x, y := ai.unitData.GetPosition()
	targetX, targetY := ai.playerTarget.unit.GetPosition()
	xDiff, yDiff := float64(targetX-x), float64(targetY-y)
	xTrans, yTrans := 127, 127
	// rotation := ai.playerTarget.unit.GetRotation()
	newRotation := utility.RadianToDegress(math.Atan2(yDiff, xDiff))
	// log.Println("Rotation", newRotation)
	maxSpeed := units.UnitSpeed * 2
	if xDiff > maxSpeed || xDiff < -maxSpeed && yDiff > maxSpeed || yDiff < -maxSpeed {
		if newRotation <= -90 && newRotation >= -180 { // top right
			trans := math.Abs(newRotation) - 90
			transXRatio := trans / 90.0
			transYRation := 1 - trans/90.0
			xTrans, yTrans = int(transXRatio*40), int(transYRation*40)
		} else if newRotation <= 0 && newRotation >= -90 { // top left
			trans := math.Abs(newRotation) - 90
			transXRatio := trans / 90.0
			transYRation := 1 - trans/90.0
			xTrans, yTrans = int(transXRatio*-220), int(transYRation*40)
		} else if newRotation <= 90 && newRotation >= 0 { // bottom left
			trans := newRotation
			transXRatio := trans / 90.0
			transYRation := 1 - trans/90.0
			xTrans, yTrans = int(transXRatio*-40), int(transYRation*-220)
		} else if newRotation <= 180 && newRotation >= 90 { // bottom right
			trans := math.Abs(newRotation)
			transXRatio := trans / 90.0
			transYRation := 1 - trans/90.0
			xTrans, yTrans = int(transXRatio*40), int(transYRation*220)
		}
	}

	byteRotation := byte(0)
	if newRotation < 0 {
		byteRotation = byte((newRotation * -2) / 360 * 255)
	} else {
		byteRotation = byte(newRotation / 360 * 255)
	}

	// log.Println("ByteRotation", byteRotation)
	ai.unitData.SetPosition(x, y, byte(xTrans), byte(yTrans))
	ai.unitData.GetCollider().Update(x, y, byteRotation) // maybe need offset
	updatedUnitIDs[ai.unitData.GetID()] = true

	return updatedUnitIDs
}
