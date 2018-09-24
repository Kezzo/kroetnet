package main

import (
	"testing"
)

func TestMatchEndMsgEncodeDecode(t *testing.T) {

	matchendmsg := MatchEndMsg{matchEndMsgID}
	encoded := matchendmsg.Encode()
	decoded := DecodeMatchEndMsg(encoded)
	if matchendmsg != decoded {
		t.Errorf("Encoded and Decoded MatchEndMsg is not the same")
	}
}
