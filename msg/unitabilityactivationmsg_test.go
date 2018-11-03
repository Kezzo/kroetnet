package msg

import (
	"testing"
)

func TestUnitAbilityActivationMsgEncodeDecode(t *testing.T) {

	unitAbilityactivationmsg := UnitAbilityActivationMsg{UnitAbilityActivationMsgID, 2, 33, 230, 10, 20}
	encoded := unitAbilityactivationmsg.Encode()
	decoded := DecodeUnitAbilityActivationMsg(encoded)
	if unitAbilityactivationmsg != decoded {
		t.Errorf("Encoded and Decoded UnitAbilityActivationMsg is not the same")
	}
}
