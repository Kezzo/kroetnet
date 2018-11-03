package msg

import (
	"testing"
)

func TestAbilityInputMsgEncodeDecode(t *testing.T) {
	Abilityinputmsg := AbilityInputMsg{AbilityInputMsgID, 125, 200, 100, 50}
	encoded := Abilityinputmsg.Encode()
	decoded := DecodeAbilityInputMsg(encoded)
	if Abilityinputmsg != decoded {
		t.Errorf("Encoded and Decoded AbilityInputMsg is not the same")
	}
}
