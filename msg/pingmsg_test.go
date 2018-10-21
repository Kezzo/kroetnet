package msg

import (
	"testing"
)

func TestPingMsgEncodeDecode(t *testing.T) {

	pingmsg := PingMsg{PingMsgID, 15377262820688280}
	encoded := pingmsg.Encode()
	decoded := DecodePingMsg(encoded)
	if pingmsg != decoded {
		t.Errorf("Encoded and Decoded PingMsg is not the same")
	}
}
