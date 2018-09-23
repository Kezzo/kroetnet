package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// MatchEndAckMsg Payload for incoming commnication
type MatchEndAckMsg struct {
	MessageID,
	PlayerID interface{}
}

// Encode transforms struct into byte array
func (m MatchEndAckMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeMatchEndAckMsg transforms a byte array into a MatchEndAckMsg
func DecodeMatchEndAckMsg(buffer []byte) MatchEndAckMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg MatchEndAckMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
