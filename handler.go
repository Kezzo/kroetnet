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
	timeResp := TimeSyncRespMsg{
		MessageID:                   timeRespMsgID,
		TransmissionTimestamp:       binary.LittleEndian.Uint64(buf[1:]),
		ServerReceptionTimestamp:    recvTime.UnixNano() / 100,
		ServerTransmissionTimestamp: time.Now().UnixNano() / 100}

	// nano seconcs / 100 == ticks

	rsp := timeResp.Encode()

	// send data
	reponseClient(pc, addr, rsp)
}

func handleInputMsg(pc net.PacketConn, addr net.Addr, buf []byte) {
	inputmsg := InputMsg{
		MessageID:   binary.LittleEndian.Uint64(buf[0:1]),
		PlayerID:    binary.LittleEndian.Uint64(buf[1:2]),
		Translation: binary.LittleEndian.Uint64(buf[2:6]),
		Rotation:    binary.LittleEndian.Uint64(buf[6:10]),
		Frame:       binary.LittleEndian.Uint64(buf[10:11])}

	// validate moves
	log.Println("Pkg Received: ", inputmsg)
	rsp := inputmsg.Encode()
	// response with position (X,Y)
	reponseClient(pc, addr, rsp)
}
