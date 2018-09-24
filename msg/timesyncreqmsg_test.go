package msg

import (
	"testing"
)

func TestTimeSyncReqMsgEncodeDecode(t *testing.T) {

	timesyncreqmsg := TimeSyncReqMsg{TimeSyncDoneMsgID, 15377262820688280}
	encoded := timesyncreqmsg.Encode()
	decoded := DecodeTimeSyncReqMsg(encoded)
	if timesyncreqmsg != decoded {
		t.Errorf("Encoded and Decoded TimeSyncReqMsg not the same")
	}
}
