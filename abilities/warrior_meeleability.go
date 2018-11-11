package abilities

import (
	"kroetnet/player"
	"kroetnet/utility"
	"log"
	"math"
)

// WarriorMeeleAbility ..
type WarriorMeeleAbility struct {
	data      AbilityData
	triggered bool
}

var activationDelay = byte(10)
var duration = byte(0)

// Init ...
func (a *WarriorMeeleAbility) Init(data AbilityData) AbilityData {
	data.ActivationFrame = byte(math.Mod(float64(data.StartFrame+activationDelay), 255.))
	data.EndFrame = data.ActivationFrame
	a.data = data

	//log.Println("Init with ActivationFrame: ", data.ActivationFrame, " EndFrame: ", data.EndFrame)

	return data
}

// Tick ...
func (a *WarriorMeeleAbility) Tick(players []*player.Player, frame byte, updatedUnitIDs map[byte]bool) map[byte]bool {
	if !a.triggered && utility.IsFrameNowOrInPast(a.data.ActivationFrame, frame) {
		//log.Println("Did trigger at: ", frame, " ActivationFrame: ", a.data.ActivationFrame)
		// do damage
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
