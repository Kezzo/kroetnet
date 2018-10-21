package msg

// SpellInputMsg Payload for incoming commnication
type SpellInputMsg struct {
	MessageID,
	PlayerID,
	SpellID,
	Rotation,
	StartFrame byte
}

// Encode transforms struct into byte array
func (m SpellInputMsg) Encode() []byte {
	buf := make([]byte, 6)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	buf[2] = m.SpellID
	buf[3] = m.Rotation
	buf[4] = m.StartFrame

	return buf
}

// DecodeSpellInputMsg transforms a byte array into a InputMsg
func DecodeSpellInputMsg(buf []byte) SpellInputMsg {
	spellinputmsg := SpellInputMsg{
		MessageID:  buf[0],
		PlayerID:   buf[1],
		SpellID:    buf[2],
		Rotation:   buf[3],
		StartFrame: buf[4]}
	return spellinputmsg
}
