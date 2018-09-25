package msg

import "encoding/binary"

// MatchEndAckMsg Payload for incoming commnication
type MatchEndAckMsg struct {
	MessageID byte
	PlayerID  uint64
}

// Encode transforms struct into byte array
func (m MatchEndAckMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.LittleEndian.PutUint64(buf[1:], m.PlayerID)
	return buf
}

// DecodeMatchEndAckMsg transforms a byte array into a MatchEndAckMsg
func DecodeMatchEndAckMsg(buf []byte) MatchEndAckMsg {
	matchendackmsg := MatchEndAckMsg{
		MessageID: buf[0],
		PlayerID:  binary.LittleEndian.Uint64(buf[1:9])}
	return matchendackmsg
}
