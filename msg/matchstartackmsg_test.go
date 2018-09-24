package msg

import (
	"testing"
)

func TestMatchStartAckMsgEncodeDecode(t *testing.T) {

	matchstartackmsg := MatchStartAckMsg{MatchStartAckMsgID, 10}
	encoded := matchstartackmsg.Encode()
	decoded := DecodeMatchStartAckMsg(encoded)
	if matchstartackmsg != decoded {
		t.Errorf("Encoded and Decoded MatchStartAckMsg is not the same")
	}
}
