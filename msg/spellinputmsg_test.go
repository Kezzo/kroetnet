package msg

import (
	"testing"
)

func TestSpellInputMsgEncodeDecode(t *testing.T) {
	spellinputmsg := SpellInputMsg{SpellInputMsgID, 125, 200, 100, 50}
	encoded := spellinputmsg.Encode()
	decoded := DecodeSpellInputMsg(encoded)
	if spellinputmsg != decoded {
		t.Errorf("Encoded and Decoded SpellInputMsg is not the same")
	}
}
