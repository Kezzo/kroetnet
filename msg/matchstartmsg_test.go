package msg

import (
	"testing"
)

func TestMatchStartMsgEncodeDecode(t *testing.T) {

	matchstartmsg := MatchStartMsg{MatchStartMsgID, 15377262820688280}
	encoded := matchstartmsg.Encode()
	decoded := DecodeMatchStartMsg(encoded)
	if matchstartmsg != decoded {
		t.Errorf("Encoded and Decoded MatchStartMsg is not the same")
	}
}
