package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// UnitStateMsg Payload for incoming commnication
type UnitStateMsg struct {
	MessageID,
	UnitID,
	XPosition,
	YPosition,
	Rotation,
	Frame interface{}
}

// Encode transforms struct into byte array
func (m UnitStateMsg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeUnitStateMsg transforms a byte array into a UnitStateMsg
func DecodeUnitStateMsg(buffer []byte) UnitStateMsg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var msg UnitStateMsg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return msg
}
