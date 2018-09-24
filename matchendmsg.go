package main

// MatchEndMsg Payload for incoming commnication
type MatchEndMsg struct {
	MessageID byte
}

// Encode transforms struct into byte array
func (m MatchEndMsg) Encode() []byte {
	buf := make([]byte, 1)
	buf[0] = m.MessageID
	return buf
}

// DecodeMatchEndMsg transforms a byte array into a MatchEndMsg
func DecodeMatchEndMsg(buf []byte) MatchEndMsg {
	matchendmsg := MatchEndMsg{MessageID: buf[0]}
	return matchendmsg
}
