package msg

import (
	"testing"
)

func TestTimeSyncDoneMsgEncodeDecode(t *testing.T) {

	timesyncdonemsg := TimeSyncDoneMsg{TimeSyncDoneMsgID, 10}
	encoded := timesyncdonemsg.Encode()
	decoded := DecodeTimeSyncDoneMsg(encoded)
	if timesyncdonemsg != decoded {
		t.Errorf("Encoded and Decoded TimeSyncDoneMsg not the same")
	}
}
