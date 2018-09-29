package msg

import (
	"testing"
)

func TestTimeSyncDoneAckMsgEncodeDecode(t *testing.T) {
	timesyncdoneackmsg := TimeSyncDoneAckMsg{TimeSyncDoneAckMsgID, 10}
	encoded := timesyncdoneackmsg.Encode()
	decoded := DecodeTimeSyncDoneAckMsg(encoded)
	if timesyncdoneackmsg != decoded {
		t.Errorf("Encoded and Decoded TimeSyncDoneAckMsg is not the same")
	}
}
