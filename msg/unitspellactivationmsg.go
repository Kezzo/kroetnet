package msg

// UnitSpellActivationMsg Payload for outgoing commnication
type UnitSpellActivationMsg struct {
	MessageID,
	UnitID,
	SpellID,
	Rotation,
	StartFrame,
	ActivationFrame byte
}

// Encode transforms struct into byte array
func (m UnitSpellActivationMsg) Encode() []byte {
	buf := make([]byte, 12)
	buf[0] = m.MessageID
	buf[1] = m.UnitID
	buf[2] = m.SpellID
	buf[3] = m.Rotation
	buf[4] = m.StartFrame
	buf[5] = m.ActivationFrame
	return buf
}

// DecodeUnitSpellActivationMsg transforms a byte array into a UnitSpellActivationMsg
func DecodeUnitSpellActivationMsg(buf []byte) UnitSpellActivationMsg {
	UnitSpellActivationMsg := UnitSpellActivationMsg{
		MessageID:       buf[0],
		UnitID:          buf[1],
		SpellID:         buf[2],
		Rotation:        buf[3],
		StartFrame:      buf[4],
		ActivationFrame: buf[5]}
	return UnitSpellActivationMsg
}
