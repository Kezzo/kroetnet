package msg

import "encoding/binary"

// TimeSyncRespMsg Payload for incoming commnication
type TimeSyncRespMsg struct {
	MessageID byte
	TransmissionTimestamp,
	ServerReceptionTimestamp,
	ServerTransmissionTimestamp uint64
}

// Encode transforms struct into byte array
func (m TimeSyncRespMsg) Encode() []byte {
	buf := make([]byte, 25)
	buf[0] = m.MessageID
	binary.BigEndian.PutUint64(buf[1:], m.TransmissionTimestamp)
	binary.BigEndian.PutUint64(buf[9:], m.ServerReceptionTimestamp)
	binary.BigEndian.PutUint64(buf[17:], m.ServerTransmissionTimestamp)

	return buf
}

// DecodeTimeSyncRespMsg transforms a byte array into a TimeSyncRespMsg
func DecodeTimeSyncRespMsg(buf []byte) TimeSyncRespMsg {
	timesyncrespmsg := TimeSyncRespMsg{
		MessageID:                   buf[0],
		TransmissionTimestamp:       binary.BigEndian.Uint64(buf[1:9]),
		ServerReceptionTimestamp:    binary.BigEndian.Uint64(buf[9:17]),
		ServerTransmissionTimestamp: binary.BigEndian.Uint64(buf[17:])}
	return timesyncrespmsg
}
