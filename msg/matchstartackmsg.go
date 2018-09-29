package msg

// MatchStartAckMsg Payload for incoming commnication
type MatchStartAckMsg struct {
	MessageID,
	PlayerID byte
}

// Encode transforms struct into byte array
func (m MatchStartAckMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	return buf
}

// DecodeMatchStartAckMsg transforms a byte array into a MatchStartAckMsg
func DecodeMatchStartAckMsg(buf []byte) MatchStartAckMsg {
	matchstartackmsg := MatchStartAckMsg{
		MessageID: buf[0],
		PlayerID:  buf[1]}
	return matchstartackmsg
}
