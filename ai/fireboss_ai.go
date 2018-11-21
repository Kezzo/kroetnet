package ai

import (
	"kroetnet/units"
	"math"
)

// FireBoss ...
type FireBoss struct {
	unitData  units.Unit
	triggered bool
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
func (ai *FireBoss) Tick(units []units.Unit, updatedUnitIDs map[byte]bool) map[byte]bool {
	if !ai.triggered {
		possibleTargets := []Target{}
		for _, unitData := range units {
			aiX, aiY := ai.unitData.GetPosition()
			uX, uY := unitData.GetPosition()
			dist := math.Sqrt(math.Pow(float64(aiX*uX), 2) + math.Pow(float64(aiY*uY), 2))
			possibleTargets = append(possibleTargets, Target{unitData, dist})
		}
		targetUnit := ai.unitData
		for _, target := range possibleTargets {
			if target.unit.GetID() == 0 {
				// if target.distance < min && target.unit.IsPlayer() == true {
				targetUnit = target.unit
			}
		}
		for _, unit := range units {
			if unit.GetID() == ai.unitData.GetID() {
				x, y := unit.GetPosition()
				targetX, targetY := targetUnit.GetPosition()
				xDiff, yDiff := float64(x-targetX), float64(y-targetY)
				xTrans, yTrans := 127, 127
				// dont run inside of player ( most of the time )
				if xDiff > 400 || xDiff < -400 && yDiff > 400 || yDiff < -400 {
					if xDiff > 0 && yDiff > 0 {
						xTrans, yTrans = 80, 80
					} else if xDiff < 0 && yDiff < 0 {
						xTrans, yTrans = 205, 205
					} else if xDiff < 0 && yDiff > 0 {
						xTrans, yTrans = 205, 80
					} else if xDiff > 0 && yDiff < 0 {
						xTrans, yTrans = 80, 205
					}
				}
				unit.SetPosition(x, y, byte(xTrans), byte(yTrans))
			}
		}
		updatedUnitIDs[ai.unitData.GetID()] = true

	}
	return updatedUnitIDs
}
