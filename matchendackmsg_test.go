package main

import (
	"testing"
)

func TestMatchEndAckMsgEncodeDecode(t *testing.T) {

	matchendackmsg := MatchEndAckMsg{matchStartAckMsgID, 10}
	encoded := matchendackmsg.Encode()
	decoded := DecodeMatchEndAckMsg(encoded)
	if matchendackmsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
