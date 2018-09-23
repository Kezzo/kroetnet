package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// TimeSyncDoneMsg Payload for incoming commnication
type TimeSyncDoneMsg struct {
	MessageID,
	PlayerID interface{}
}

// Encode transforms struct into byte array
func (m TimeSyncDoneMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeTimeSyncDoneMsg transforms a byte array into a TimeSyncDoneMsg
func DecodeTimeSyncDoneMsg(buffer []byte) TimeSyncDoneMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg TimeSyncDoneMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
