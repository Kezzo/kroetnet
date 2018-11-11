package abilities

import "kroetnet/units"

// Ability ...
type Ability interface {
	Init(data AbilityData, caster units.Unit) AbilityData
	Tick(players []units.Unit, frame byte, updatedUnitIDs map[byte]bool) map[byte]bool
	IsActive(currentFrame byte) bool
}
