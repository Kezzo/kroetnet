package msg

import "encoding/binary"

// TimeSyncReqMsg Payload for incoming commnication
type TimeSyncReqMsg struct {
	MessageID             byte
	TransmissionTimestamp uint64
}

// Encode transforms struct into byte array
func (m TimeSyncReqMsg) Encode() []byte {
	buf := make([]byte, 9)
	buf[0] = m.MessageID
	binary.BigEndian.PutUint64(buf[1:], m.TransmissionTimestamp)

	return buf
}

// DecodeTimeSyncReqMsg transforms a byte array into a TimeSyncReqMsg
func DecodeTimeSyncReqMsg(buf []byte) TimeSyncReqMsg {
	timesyncreqmsg := TimeSyncReqMsg{
		MessageID:             buf[0],
		TransmissionTimestamp: binary.BigEndian.Uint64(buf[1:9])}
	return timesyncreqmsg
}
