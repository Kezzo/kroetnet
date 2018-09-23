package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TimeSyncRespMsg Payload for incoming commnication
type TimeSyncRespMsg struct {
	MessageID,
	TransmissionTimestamp,
	ServerReceptionTimestamp,
	ServerTransmissionTimestamp interface{}
}

// Encode transforms struct into byte array
func (m TimeSyncRespMsg) Encode() []byte {

	buf := make([]byte, 25)
	// results into incorrect byte array
	/*buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()*/

	// requires msg struct to have proper type
	/*buf[0] = m.MessageID
	binary.LittleEndian.PutUint64(buf[1:], int64(m.TransmissionTimestamp))
	binary.LittleEndian.PutUint64(buf[9:], int64(m.ServerReceptionTimestamp))
	binary.LittleEndian.PutUint64(buf[17:], int64(m.ServerTransmissionTimestamp))*/

	return buf
}

// DecodeTimeSyncRespMsg transforms a byte array into a TimeSyncRespMsg
func DecodeTimeSyncRespMsg(buffer []byte) TimeSyncRespMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg TimeSyncRespMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
