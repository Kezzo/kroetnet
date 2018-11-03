package abilities

import "kroetnet/player"

// Ability ...
type Ability interface {
	Init(data AbilityData) AbilityData
	Tick(players []player.Player, frame byte, updatedUnitIDs map[byte]bool) map[byte]bool
	IsActive(currentFrame byte) bool
}
