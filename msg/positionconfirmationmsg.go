package msg

import "encoding/binary"

// PositionConfirmationMsg Payload for outgoing commnication
type PositionConfirmationMsg struct {
	MessageID byte
	UnitID    byte
	XPosition int32
	YPosition int32
	Frame     byte
}

// Encode transforms struct into byte array
func (m PositionConfirmationMsg) Encode() []byte {
	buf := make([]byte, 11)
	buf[0] = m.MessageID
	buf[1] = m.UnitID
	binary.LittleEndian.PutUint32(buf[2:], uint32(m.XPosition))
	binary.LittleEndian.PutUint32(buf[6:], uint32(m.YPosition))
	buf[10] = m.Frame
	return buf
}
