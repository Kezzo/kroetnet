package ai

import (
	"kroetnet/units"
	"log"
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
				transX, transY := byte(127), byte(127)
				if targetX < x {
					transX = 0
				} else if targetX > x {
					transX = 170
				}
				if targetY < y {
					transY = 0
				} else if targetY > y {
					transY = 170
				}
				log.Println("SET POS", transX, transY)
				unit.SetPosition(x, y, transX, transY)
			}
		}
		updatedUnitIDs[ai.unitData.GetID()] = true

	}
	return updatedUnitIDs
}
