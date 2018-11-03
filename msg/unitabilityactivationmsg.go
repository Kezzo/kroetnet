package msg

// UnitAbilityActivationMsg Payload for outgoing commnication
type UnitAbilityActivationMsg struct {
	MessageID,
	UnitID,
	AbilityID,
	Rotation,
	StartFrame,
	ActivationFrame byte
}

// Encode transforms struct into byte array
func (m UnitAbilityActivationMsg) Encode() []byte {
	buf := make([]byte, 12)
	buf[0] = m.MessageID
	buf[1] = m.UnitID
	buf[2] = m.AbilityID
	buf[3] = m.Rotation
	buf[4] = m.StartFrame
	buf[5] = m.ActivationFrame
	return buf
}

// DecodeUnitAbilityActivationMsg transforms a byte array into a UnitAbilityActivationMsg
func DecodeUnitAbilityActivationMsg(buf []byte) UnitAbilityActivationMsg {
	UnitAbilityActivationMsg := UnitAbilityActivationMsg{
		MessageID:       buf[0],
		UnitID:          buf[1],
		AbilityID:       buf[2],
		Rotation:        buf[3],
		StartFrame:      buf[4],
		ActivationFrame: buf[5]}
	return UnitAbilityActivationMsg
}
