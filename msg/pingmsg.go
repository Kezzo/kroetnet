package msg

import "encoding/binary"

// PingMsg Payload for incoming commnication
type PingMsg struct {
	MessageID             byte
	TransmissionTimestamp uint64
}

// Encode transforms struct into byte array
func (m PingMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.LittleEndian.PutUint64(buf[1:], m.TransmissionTimestamp)

	return buf
}

// DecodePingMsg transforms a byte array into a PingMsg
func DecodePingMsg(buf []byte) PingMsg {
	PingMsg := PingMsg{
		MessageID:             buf[0],
		TransmissionTimestamp: binary.LittleEndian.Uint64(buf[1:9])}
	return PingMsg
}
