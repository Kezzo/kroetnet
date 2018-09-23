package main

import (
	"testing"
)

func TestUnitStateMsgEncodeDecode(t *testing.T) {

	unitstatemsg := UnitStateMsg{unitStateMsgID, 125, 43857, 1345300, 125, 10}
	encoded := unitstatemsg.Encode()
	decoded := DecodeUnitStateMsg(encoded)
	if unitstatemsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
