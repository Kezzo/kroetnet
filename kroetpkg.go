package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// KroetPkg Payload for commnication
type KroetPkg struct {
	MessageID,
	PlayerID,
	Translation,
	Rotation,
	Frame string
}

// KroetPkgs Interface
type KroetPkgs interface {
	Encode() []byte
}

// Encode transforms struct into byte array
func (p KroetPkg) Encode() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal("Encoding Pkg error: ", err)
	}
	return buf.Bytes()
}

// DecodeKroetPkg transforms a byte array into a KroetPkg
func DecodeKroetPkg(buffer []byte) KroetPkg {
	buf := &bytes.Buffer{}
	buf.Write(buffer)
	var pkg KroetPkg
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&pkg)
	if err != nil {
		log.Fatal("Decoding Pkg error: ", err)
	}
	return pkg
}
