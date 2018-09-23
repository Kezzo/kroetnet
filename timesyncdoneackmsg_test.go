package main

import (
	"testing"
)

func TestTimeSyncDoneAckMsgEncodeDecode(t *testing.T) {

	timesyncdoneackmsg := TimeSyncDoneAckMsg{timeSyncDoneAckMsgID}
	encoded := timesyncdoneackmsg.Encode()
	decoded := DecodeTimeSyncDoneAckMsg(encoded)
	if timesyncdoneackmsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
