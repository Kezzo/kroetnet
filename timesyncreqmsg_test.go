package main

import (
	"testing"
)

func TestTimeSyncReqMsgEncodeDecode(t *testing.T) {

	timesyncreqmsg := TimeSyncReqMsg{timeSyncDoneMsgID, 15377262820688280}
	encoded := timesyncreqmsg.Encode()
	decoded := DecodeTimeSyncReqMsg(encoded)
	if timesyncreqmsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
