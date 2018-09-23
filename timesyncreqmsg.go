package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TimeSyncReqMsg Payload for incoming commnication
type TimeSyncReqMsg struct {
	MessageID,
	TransmissionTimestamp interface{}
}

// Encode transforms struct into byte array
func (m TimeSyncReqMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeTimeSyncReqMsg transforms a byte array into a TimeSyncReqMsg
func DecodeTimeSyncReqMsg(buffer []byte) TimeSyncReqMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg TimeSyncReqMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
