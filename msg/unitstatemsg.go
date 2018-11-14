package msg

import "encoding/binary"

// UnitStateMsg Payload for outgoing commnication
type UnitStateMsg struct {
	MessageID     byte
	UnitID        byte
	XPosition     int32
	YPosition     int32
	Rotation      byte
	HealthPercent byte
	Frame         byte
}

// Encode transforms struct into byte array
func (m UnitStateMsg) Encode() []byte {
	buf := make([]byte, 13)
	buf[0] = m.MessageID
	buf[1] = m.UnitID
	binary.LittleEndian.PutUint32(buf[2:], uint32(m.XPosition))
	binary.LittleEndian.PutUint32(buf[6:], uint32(m.YPosition))
	buf[10] = m.Rotation
	buf[11] = m.HealthPercent
	buf[12] = m.Frame
	return buf
}

// DecodeUnitStateMsg transforms a byte array into a UnitStateMsg
func DecodeUnitStateMsg(buf []byte) UnitStateMsg {
	unitstatemsg := UnitStateMsg{
		MessageID:     buf[0],
		UnitID:        buf[1],
		XPosition:     int32(binary.LittleEndian.Uint32(buf[2:6])),
		YPosition:     int32(binary.LittleEndian.Uint32(buf[6:10])),
		Rotation:      buf[10],
		HealthPercent: buf[11],
		Frame:         buf[12]}
	return unitstatemsg
}
