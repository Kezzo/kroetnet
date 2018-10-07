package msg

// TimeSyncDoneAckMsg Payload for outgoing commnication
type TimeSyncDoneAckMsg struct {
	MessageID,
	PlayerID byte
}

// Encode transforms struct into byte array
func (m TimeSyncDoneAckMsg) Encode() []byte {
	buf := make([]byte, 2)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	return buf
}

// DecodeTimeSyncDoneAckMsg transforms a byte array into a TimeSyncDoneAckMsg
func DecodeTimeSyncDoneAckMsg(buf []byte) TimeSyncDoneAckMsg {
	timesyncdoneackmsg := TimeSyncDoneAckMsg{
		MessageID: buf[0],
		PlayerID:  buf[1]}
	return timesyncdoneackmsg
}
