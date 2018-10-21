package msg

import (
	"testing"
)

func TestUnitSpellActivationMsgEncodeDecode(t *testing.T) {

	unitspellactivationmsg := UnitSpellActivationMsg{UnitSpellActivationMsgID, 2, 33, 230, 10, 20}
	encoded := unitspellactivationmsg.Encode()
	decoded := DecodeUnitSpellActivationMsg(encoded)
	if unitspellactivationmsg != decoded {
		t.Errorf("Encoded and Decoded UnitSpellActivationMsg is not the same")
	}
}
