package main

import (
	"encoding/binary"
)

// TimeSyncDoneMsg Payload for incoming commnication
type TimeSyncDoneMsg struct {
	MessageID byte
	PlayerID  uint64
}

// Encode transforms struct into byte array
func (m TimeSyncDoneMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.BigEndian.PutUint64(buf[1:], m.PlayerID)
	return buf
}

// DecodeTimeSyncDoneMsg transforms a byte array into a TimeSyncDoneMsg
func DecodeTimeSyncDoneMsg(buf []byte) TimeSyncDoneMsg {
	timesyncdonemsg := TimeSyncDoneMsg{
		MessageID: buf[0],
		PlayerID:  binary.BigEndian.Uint64(buf[1:])}
	return timesyncdonemsg
}
