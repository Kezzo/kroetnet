package main

// TimeSyncDoneAckMsg Payload for incoming commnication
type TimeSyncDoneAckMsg struct {
	MessageID byte
}

// Encode transforms struct into byte array
func (m TimeSyncDoneAckMsg) Encode() []byte {
	buf := make([]byte, 1)
	buf[0] = m.MessageID
	return buf
}

// DecodeTimeSyncDoneAckMsg transforms a byte array into a TimeSyncDoneAckMsg
func DecodeTimeSyncDoneAckMsg(buf []byte) TimeSyncDoneAckMsg {
	timesyncdoneackmsg := TimeSyncDoneAckMsg{MessageID: buf[0]}
	return timesyncdoneackmsg
}
