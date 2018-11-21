package ai

import (
	"kroetnet/units"
)

// AI ...
type AI interface {
	Init(unit units.Unit)
	Tick(units []units.Unit, updatedUnitIDs map[byte]bool) map[byte]bool
}
