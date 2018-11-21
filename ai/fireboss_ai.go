package ai

import (
	"kroetnet/units"
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
	xDiff, yDiff := float64(x-targetX), float64(y-targetY)
	xTrans, yTrans := 127, 127
	rotation := ai.playerTarget.unit.GetRotation()
	maxSpeed := units.UnitSpeed * 2
	if xDiff > maxSpeed || xDiff < -maxSpeed && yDiff > maxSpeed || yDiff < -maxSpeed {
		if xDiff > 0 && yDiff > 0 {
			xTrans, yTrans = 80, 80
			rotation = 60
		} else if xDiff < 0 && yDiff < 0 {
			xTrans, yTrans = 205, 205
			rotation = 255
		} else if xDiff < 0 && yDiff > 0 {
			xTrans, yTrans = 205, 80
			rotation = 0
		} else if xDiff > 0 && yDiff < 0 {
			xTrans, yTrans = 80, 205
			rotation = 120
		}
		ai.unitData.SetPosition(x, y, byte(xTrans), byte(yTrans))
		ai.unitData.GetCollider().Update(x, y, rotation)
		updatedUnitIDs[ai.unitData.GetID()] = true
	}

	return updatedUnitIDs
}
