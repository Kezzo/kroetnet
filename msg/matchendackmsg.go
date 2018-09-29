package msg

// MatchEndAckMsg Payload for incoming commnication
type MatchEndAckMsg struct {
	MessageID,
	PlayerID byte
}

// Encode transforms struct into byte array
func (m MatchEndAckMsg) Encode() []byte {
	buf := make([]byte, 2)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	return buf
}

// DecodeMatchEndAckMsg transforms a byte array into a MatchEndAckMsg
func DecodeMatchEndAckMsg(buf []byte) MatchEndAckMsg {
	matchendackmsg := MatchEndAckMsg{
		MessageID: buf[0],
		PlayerID:  buf[1]}
	return matchendackmsg
}
