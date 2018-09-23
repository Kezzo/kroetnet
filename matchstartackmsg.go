package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// MatchStartAckMsg Payload for incoming commnication
type MatchStartAckMsg struct {
	MessageID,
	PlayerID interface{}
}

// Encode transforms struct into byte array
func (m MatchStartAckMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeMatchStartAckMsg transforms a byte array into a MatchStartAckMsg
func DecodeMatchStartAckMsg(buffer []byte) MatchStartAckMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg MatchStartAckMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
