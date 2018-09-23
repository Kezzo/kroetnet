package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TimeSyncDoneAckMsg Payload for incoming commnication
type TimeSyncDoneAckMsg struct {
	MessageID interface{}
}

// Encode transforms struct into byte array
func (m TimeSyncDoneAckMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeTimeSyncDoneAckMsg transforms a byte array into a TimeSyncDoneAckMsg
func DecodeTimeSyncDoneAckMsg(buffer []byte) TimeSyncDoneAckMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg TimeSyncDoneAckMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
