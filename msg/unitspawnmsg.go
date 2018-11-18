package msg

import "encoding/binary"

// UnitSpawnMsg Payload for outgoing commnication
type UnitSpawnMsg struct {
	MessageID     byte
	UnitID        byte
	TeamID        byte
	UnitType      byte
	XPosition     int32
	YPosition     int32
	Rotation      byte
	HealthPercent byte
	Frame         byte
}

// Encode transforms struct into byte array
func (m UnitSpawnMsg) Encode() []byte {
	buf := make([]byte, 15)
	buf[0] = m.MessageID
	buf[1] = m.UnitID
	buf[2] = m.TeamID
	buf[3] = m.UnitType
	binary.LittleEndian.PutUint32(buf[4:], uint32(m.XPosition))
	binary.LittleEndian.PutUint32(buf[8:], uint32(m.YPosition))
	buf[12] = m.Rotation
	buf[13] = m.HealthPercent
	buf[14] = m.Frame
	return buf
}

// DecodeUnitSpawnMsg transforms a byte array into a UnitSpawnMsg
func DecodeUnitSpawnMsg(buf []byte) UnitSpawnMsg {
	unitspawnmsg := UnitSpawnMsg{
		MessageID:     buf[0],
		UnitID:        buf[1],
		TeamID:        buf[2],
		UnitType:      buf[3],
		XPosition:     int32(binary.LittleEndian.Uint32(buf[4:8])),
		YPosition:     int32(binary.LittleEndian.Uint32(buf[8:12])),
		Rotation:      buf[12],
		HealthPercent: buf[13],
		Frame:         buf[14]}
	return unitspawnmsg
}
