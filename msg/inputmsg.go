package msg

import "encoding/binary"

// InputMsg Payload for incoming commnication
type InputMsg struct {
	MessageID byte
	PlayerID,
	Translation,
	Rotation,
	Frame uint64
}

// Encode transforms struct into byte array
func (m InputMsg) Encode() []byte {
	buf := make([]byte, 33)
	buf[0] = m.MessageID
	binary.BigEndian.PutUint64(buf[1:], m.PlayerID)
	binary.BigEndian.PutUint64(buf[9:], m.Translation)
	binary.BigEndian.PutUint64(buf[17:], m.Rotation)
	binary.BigEndian.PutUint64(buf[25:], m.Frame)
	return buf
}

// DecodeInputMsg transforms a byte array into a InputMsg
func DecodeInputMsg(buf []byte) InputMsg {
	inputmsg := InputMsg{
		MessageID:   buf[0],
		PlayerID:    binary.BigEndian.Uint64(buf[1:9]),
		Translation: binary.BigEndian.Uint64(buf[9:17]),
		Rotation:    binary.BigEndian.Uint64(buf[17:25]),
		Frame:       binary.BigEndian.Uint64(buf[25:])}
	return inputmsg
}
