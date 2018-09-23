package main

import (
	"testing"
)

func TestTimeSyncRespMsgEncodeDecode(t *testing.T) {

	timesyncrespmsg := TimeSyncRespMsg{timeRespMsgID, 15377262820688281, 15377262820688282, 15377262820688283}
	encoded := timesyncrespmsg.Encode()
	decoded := DecodeTimeSyncRespMsg(encoded)
	if timesyncrespmsg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
