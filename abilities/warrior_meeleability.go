package abilities

import (
	"kroetnet/collision"
	"kroetnet/units"
	"kroetnet/utility"
	"log"
	"math"
)

// WarriorMeeleAbility ..
type WarriorMeeleAbility struct {
	data      AbilityData
	triggered bool
	caster    units.Unit
}

var meeleActivationDelay = byte(10)
var meeleDuration = byte(0)
var meeleColliderLength int32 = 4000
var meeleColliderWidthDegress float64 = 45

// Init ...
func (a *WarriorMeeleAbility) Init(data AbilityData, caster units.Unit) AbilityData {
	data.ActivationFrame = byte(math.Mod(float64(data.StartFrame+meeleActivationDelay), 255.))
	data.EndFrame = data.ActivationFrame
	a.data = data
	a.caster = caster

	//log.Println("Init with ActivationFrame: ", data.ActivationFrame, " EndFrame: ", data.EndFrame)

	return data
}

// Tick ...
func (a *WarriorMeeleAbility) Tick(units []units.Unit, frame byte, updatedUnitIDs map[byte]bool) map[byte]bool {
	if !a.triggered && utility.IsFrameNowOrInPast(a.data.ActivationFrame, frame) {
		//log.Println("Did trigger at: ", frame, " ActivationFrame: ", a.data.ActivationFrame)

		x, y := a.caster.GetPosition()
		abilityCollider := collision.NewConeCollider(x, y, a.data.Rotation, meeleColliderLength, meeleColliderWidthDegress)

		for _, unitData := range units {
			// TODO: Allow targeting team members (i.e. mfor heals)
			if a.caster.GetTeam() == unitData.GetTeam() {
				continue
			}

			if abilityCollider.IsColliding(unitData.GetCollider()) {
				log.Println("WarriorMeeleAbility collided with unit of team: ", unitData.GetTeam())

				unitData.AddDamage(30)
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
func (a *WarriorMeeleAbility) IsActive(currentFrame byte) bool {
	return utility.IsFrameInFuture(a.data.EndFrame, currentFrame)
}
