package msg

import (
	"testing"
)

func TestMatchEndAckMsgEncodeDecode(t *testing.T) {

	matchendackmsg := MatchEndAckMsg{MatchStartAckMsgID, 10}
	encoded := matchendackmsg.Encode()
	decoded := DecodeMatchEndAckMsg(encoded)
	if matchendackmsg != decoded {
		t.Errorf("Encoded and Decoded MatchEndAckMsg is not the same")
	}
}
