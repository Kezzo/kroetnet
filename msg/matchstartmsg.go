package msg

import "encoding/binary"

// MatchStartMsg Payload for incoming commnication
type MatchStartMsg struct {
	MessageID           byte
	MatchStartTimestamp uint64
}

// Encode transforms struct into byte array
func (m MatchStartMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.LittleEndian.PutUint64(buf[1:], m.MatchStartTimestamp)
	return buf
}

// DecodeMatchStartMsg transforms a byte array into a MatchStartMsg
func DecodeMatchStartMsg(buf []byte) MatchStartMsg {
	matchstartmsg := MatchStartMsg{
		MessageID:           buf[0],
		MatchStartTimestamp: binary.LittleEndian.Uint64(buf[1:])}
	return matchstartmsg
}
