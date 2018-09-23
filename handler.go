package main

import (
	"encoding/binary"
	"log"
	"net"
	"time"
)

func reponseClient(pc net.PacketConn, addr net.Addr, buf []byte) {
	log.Println("Reponse to send: ", buf, " to ", addr)
	if _, err := pc.WriteTo(buf, addr); err != nil {
		log.Fatalln("err sending data:", err)
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	response := []byte("Received your msg: " + string(buf))
	reponseClient(pc, addr, response)
}

func handleTimeSynchReq(pc net.PacketConn, addr net.Addr, buf []byte, recvTime time.Time) {
	rsp := make([]byte, 25)
	rsp[0] = MSG_TIME_RESP
	copy(rsp[1:], buf[1:])
	// nano seconcs / 100 == ticks
	binary.LittleEndian.PutUint64(rsp[9:], uint64(recvTime.UnixNano()/100))
	sendTime := time.Now().UnixNano() / 100
	binary.LittleEndian.PutUint64(rsp[17:], uint64(sendTime))
	// send data
	reponseClient(pc, addr, rsp)
}

func handleKroetPkg(pc net.PacketConn, addr net.Addr, buf []byte) {
	kroetpkg := KroetPkg{
		MessageID:   binary.LittleEndian.Uint64(buf[0:1]),
		PlayerID:    binary.LittleEndian.Uint64(buf[1:2]),
		Translation: binary.LittleEndian.Uint64(buf[2:6]),
		Rotation:    binary.LittleEndian.Uint64(buf[6:10]),
		Frame:       binary.LittleEndian.Uint64(buf[10:11])}

	// validate moves
	log.Println("Pkg Received: ", kroetpkg)
	rsp := kroetpkg.Encode()
	// response with position (X,Y)
	reponseClient(pc, addr, rsp)
}
