package msg

import "encoding/binary"

// UnitStateMsg Payload for incoming commnication
type UnitStateMsg struct {
	MessageID byte
	UnitID,
	XPosition,
	YPosition,
	Rotation,
	Frame uint64
}

// Encode transforms struct into byte array
func (m UnitStateMsg) Encode() []byte {
	buf := make([]byte, 41)
	buf[0] = m.MessageID
	binary.LittleEndian.PutUint64(buf[1:], m.UnitID)
	binary.LittleEndian.PutUint64(buf[9:], m.XPosition)
	binary.LittleEndian.PutUint64(buf[17:], m.YPosition)
	binary.LittleEndian.PutUint64(buf[25:], m.Rotation)
	binary.LittleEndian.PutUint64(buf[33:], m.Frame)
	return buf
}

// DecodeUnitStateMsg transforms a byte array into a UnitStateMsg
func DecodeUnitStateMsg(buf []byte) UnitStateMsg {
	unitstatemsg := UnitStateMsg{
		MessageID: buf[0],
		UnitID:    binary.LittleEndian.Uint64(buf[1:9]),
		XPosition: binary.LittleEndian.Uint64(buf[9:17]),
		YPosition: binary.LittleEndian.Uint64(buf[17:25]),
		Rotation:  binary.LittleEndian.Uint64(buf[25:33]),
		Frame:     binary.LittleEndian.Uint64(buf[33:])}
	return unitstatemsg
}
