package msg

import (
	"bytes"
	"testing"
)

func TestPosConfEncode(t *testing.T) {
	posmsg := PositionConfirmationMsg{PositionConfirmationMessageID,
		0, 1300, 1560, 14}
	encoded := posmsg.Encode()
	testBuf := []byte{12, 0, 20, 5, 0, 0, 24, 6, 0, 0, 14}
	if !bytes.Equal(encoded, testBuf) {
		t.Errorf("Encoded PosConf test failed")
	}
}
