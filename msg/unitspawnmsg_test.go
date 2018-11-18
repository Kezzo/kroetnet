package msg

import (
	"testing"
)

func TestUnitSpawnMsgEncodeDecode(t *testing.T) {
	unitSpawnmsg := UnitSpawnMsg{UnitSpawnMsgID, 125, 0, 0, 43857, -1345300, 230, 50, 10}
	encoded := unitSpawnmsg.Encode()
	decoded := DecodeUnitSpawnMsg(encoded)
	if unitSpawnmsg != decoded {
		t.Errorf("Encoded and Decoded UnitSpawnMsg is not the same")
	}
}
