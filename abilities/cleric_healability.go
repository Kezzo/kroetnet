package abilities

import (
	"kroetnet/collision"
	"kroetnet/units"
	"kroetnet/utility"
	"log"
	"math"
)

// ClericHealAbility ..
type ClericHealAbility struct {
	data      AbilityData
	triggered bool
	caster    units.Unit
}

var healActivationDelay = byte(20)
var healDuration = byte(0)
var healColliderLength int32 = 8000
var healColliderWidthDegress float64 = 35

// Init ...
func (a *ClericHealAbility) Init(data AbilityData, caster units.Unit) AbilityData {
	data.ActivationFrame = byte(math.Mod(float64(data.StartFrame+healActivationDelay), 255.))
	data.EndFrame = data.ActivationFrame
	a.data = data
	a.caster = caster

	//log.Println("Init with ActivationFrame: ", data.ActivationFrame, " EndFrame: ", data.EndFrame)

	return data
}

// Tick ...
func (a *ClericHealAbility) Tick(units []units.Unit, frame byte, updatedUnitIDs map[byte]bool) map[byte]bool {
	if !a.triggered && utility.IsFrameNowOrInPast(a.data.ActivationFrame, frame) {
		//log.Println("Did trigger at: ", frame, " ActivationFrame: ", a.data.ActivationFrame)

		x, y := a.caster.GetPosition()
		abilityCollider := collision.NewConeCollider(x, y, a.data.Rotation, healColliderLength, healColliderWidthDegress)

		for _, unitData := range units {
			// don't heal enemies
			if a.caster.GetTeam() != unitData.GetTeam() {
				continue
			}

			if abilityCollider.IsColliding(unitData.GetCollider()) {
				log.Println("ClericHealAbility collided with unit of team: ", unitData.GetTeam())

				unitData.AddHeal(50)
				updatedUnitIDs[unitData.GetID()] = true
			}
		}

		a.triggered = true
		return updatedUnitIDs
	}

	log.Println("Didnt trigger at: ", frame)
	return updatedUnitIDs
}

// IsActive ...
func (a *ClericHealAbility) IsActive(currentFrame byte) bool {
	return utility.IsFrameInFuture(a.data.EndFrame, currentFrame)
}
