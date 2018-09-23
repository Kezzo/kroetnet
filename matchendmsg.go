package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// MatchEndMsg Payload for incoming commnication
type MatchEndMsg struct {
	MessageID interface{}
}

// Encode transforms struct into byte array
func (m MatchEndMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeMatchEndMsg transforms a byte array into a MatchEndMsg
func DecodeMatchEndMsg(buffer []byte) MatchEndMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg MatchEndMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
