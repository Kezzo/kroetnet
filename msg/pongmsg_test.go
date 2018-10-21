package msg

import (
	"testing"
)

func TestPongMsgEncodeDecode(t *testing.T) {

	pongmsg := PongMsg{PongMsgID, 15377262820688280}
	encoded := pongmsg.Encode()
	decoded := DecodePongMsg(encoded)
	if pongmsg != decoded {
		t.Errorf("Encoded and Decoded PongMsg is not the same")
	}
}
