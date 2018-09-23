package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// InputMsg Payload for incoming commnication
type InputMsg struct {
	MessageID,
	PlayerID,
	Translation,
	Rotation,
	Frame interface{}
}

// Encode transforms struct into byte array
func (m InputMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeInputMsg transforms a byte array into a InputMsg
func DecodeInputMsg(buffer []byte) InputMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg InputMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
