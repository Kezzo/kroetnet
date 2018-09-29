package msg

import (
	"testing"
)

func TestUnitStateMsgEncodeDecode(t *testing.T) {

	unitstatemsg := UnitStateMsg{UnitStateMsgID, 125, 43857, -1345300, 125, 10}
	encoded := unitstatemsg.Encode()
	decoded := DecodeUnitStateMsg(encoded)
	if unitstatemsg != decoded {
		t.Errorf("Encoded and Decoded UnitStateMsg is not the same")
	}
}
