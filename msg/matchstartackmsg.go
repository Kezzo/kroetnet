package msg

import "encoding/binary"

// MatchStartAckMsg Payload for incoming commnication
type MatchStartAckMsg struct {
	MessageID byte
	PlayerID  uint64
}

// Encode transforms struct into byte array
func (m MatchStartAckMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.BigEndian.PutUint64(buf[1:], m.PlayerID)
	return buf
}

// DecodeMatchStartAckMsg transforms a byte array into a MatchStartAckMsg
func DecodeMatchStartAckMsg(buf []byte) MatchStartAckMsg {
	matchstartackmsg := MatchStartAckMsg{
		MessageID: buf[0],
		PlayerID:  binary.BigEndian.Uint64(buf[1:])}
	return matchstartackmsg
}
