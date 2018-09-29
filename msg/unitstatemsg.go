package msg

import "encoding/binary"

// UnitStateMsg Payload for incoming commnication
type UnitStateMsg struct {
	MessageID byte
	UnitID    byte
	XPosition int32
	YPosition int32
	Rotation  byte
	Frame     byte
}

// Encode transforms struct into byte array
func (m UnitStateMsg) Encode() []byte {
	buf := make([]byte, 19)
	buf[0] = m.MessageID
	buf[1] = m.UnitID
	binary.LittleEndian.PutUint32(buf[2:], uint32(m.XPosition))
	binary.LittleEndian.PutUint32(buf[9:], uint32(m.YPosition))
	buf[17] = m.Rotation
	buf[18] = m.Frame
	return buf
}

// DecodeUnitStateMsg transforms a byte array into a UnitStateMsg
func DecodeUnitStateMsg(buf []byte) UnitStateMsg {
	unitstatemsg := UnitStateMsg{
		MessageID: buf[0],
		UnitID:    buf[1],
		XPosition: int32(binary.LittleEndian.Uint32(buf[2:8])),
		YPosition: int32(binary.LittleEndian.Uint32(buf[9:16])),
		Rotation:  buf[17],
		Frame:     buf[18]}
	return unitstatemsg
}
