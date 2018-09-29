package msg

// InputMsg Payload for incoming commnication
type InputMsg struct {
	MessageID,
	PlayerID,
	XTranslation,
	YTranslation,
	Rotation,
	Frame byte
}

// Encode transforms struct into byte array
func (m InputMsg) Encode() []byte {
	buf := make([]byte, 6)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	buf[2] = m.XTranslation
	buf[3] = m.YTranslation
	buf[4] = m.Rotation
	buf[5] = m.Frame

	return buf
}

// DecodeInputMsg transforms a byte array into a InputMsg
func DecodeInputMsg(buf []byte) InputMsg {
	inputmsg := InputMsg{
		MessageID:    buf[0],
		PlayerID:     buf[1],
		XTranslation: buf[2],
		YTranslation: buf[3],
		Rotation:     buf[4],
		Frame:        buf[5]}
	return inputmsg
}
