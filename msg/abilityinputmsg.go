package msg

// AbilityInputMsg Payload for incoming commnication
type AbilityInputMsg struct {
	MessageID,
	PlayerID,
	AbilityID,
	Rotation,
	StartFrame byte
}

// Encode transforms struct into byte array
func (m AbilityInputMsg) Encode() []byte {
	buf := make([]byte, 6)
	buf[0] = m.MessageID
	buf[1] = m.PlayerID
	buf[2] = m.AbilityID
	buf[3] = m.Rotation
	buf[4] = m.StartFrame

	return buf
}

// DecodeAbilityInputMsg transforms a byte array into a InputMsg
func DecodeAbilityInputMsg(buf []byte) AbilityInputMsg {
	Abilityinputmsg := AbilityInputMsg{
		MessageID:  buf[0],
		PlayerID:   buf[1],
		AbilityID:  buf[2],
		Rotation:   buf[3],
		StartFrame: buf[4]}
	return Abilityinputmsg
}
