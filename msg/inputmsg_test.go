package msg

import (
	"testing"
)

func TestInputMsgEncodeDecode(t *testing.T) {

	inputmsg := InputMsg{InputMsgID, 125, 200, 100, 10}
	encoded := inputmsg.Encode()
	decoded := DecodeInputMsg(encoded)
	if inputmsg != decoded {
		t.Errorf("Encoded and Decoded InputMsg is not the same")
	}
}
