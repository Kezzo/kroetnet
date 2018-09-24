package msg

import (
	"testing"
)

func TestTimeSyncDoneAckMsgEncodeDecode(t *testing.T) {

	timesyncdoneackmsg := TimeSyncDoneAckMsg{TimeSyncDoneAckMsgID}
	encoded := timesyncdoneackmsg.Encode()
	decoded := DecodeTimeSyncDoneAckMsg(encoded)
	if timesyncdoneackmsg != decoded {
		t.Errorf("Encoded and Decoded TimeSyncDoneAckMsg is not the same")
	}
}
