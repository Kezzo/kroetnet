package msg

// TimeSyncDoneMsg Payload for incoming commnication
type TimeSyncDoneMsg struct {
	MessageID,
	PlayerID byte
}

// Encode transforms struct into byte array
func (m TimeSyncDoneMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	return buf
}

// DecodeTimeSyncDoneMsg transforms a byte array into a TimeSyncDoneMsg
func DecodeTimeSyncDoneMsg(buf []byte) TimeSyncDoneMsg {
	timesyncdonemsg := TimeSyncDoneMsg{
		MessageID: buf[0],
		PlayerID:  buf[1]}
	return timesyncdonemsg
}
