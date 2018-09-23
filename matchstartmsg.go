package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// MatchStartMsg Payload for incoming commnication
type MatchStartMsg struct {
	MessageID,
	MatchStartTimestamp interface{}
}

// Encode transforms struct into byte array
func (m MatchStartMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeMatchStartMsg transforms a byte array into a MatchStartMsg
func DecodeMatchStartMsg(buffer []byte) MatchStartMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg MatchStartMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
