package main

import (
	"testing"
)

func TestMatchStartAckMsgEncodeDecode(t *testing.T) {

	matchstartackmsg := MatchStartAckMsg{matchStartAckMsgID, 10}
	encoded := matchstartackmsg.Encode()
	decoded := DecodeMatchStartAckMsg(encoded)
	if matchstartackmsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
