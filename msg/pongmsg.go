package msg

import "encoding/binary"

// PongMsg Payload for outgoing commnication
type PongMsg struct {
	MessageID             byte
	TransmissionTimestamp uint64
}

// Encode transforms struct into byte array
func (m PongMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.LittleEndian.PutUint64(buf[1:], m.TransmissionTimestamp)

	return buf
}

// DecodePongMsg transforms a byte array into a PongMsg
func DecodePongMsg(buf []byte) PongMsg {
	PongMsg := PongMsg{
		MessageID:             buf[0],
		TransmissionTimestamp: binary.LittleEndian.Uint64(buf[1:9])}
	return PongMsg
}
