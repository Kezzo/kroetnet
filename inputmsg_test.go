package main

import (
	"testing"
)

func TestInputMsgEncodeDecode(t *testing.T) {

	inputmsg := InputMsg{inputMsgID, 125, 200, 100, 10}
	encoded := inputmsg.Encode()
	decoded := DecodeInputMsg(encoded)
	if inputmsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
