package main

import (
	"encoding/binary"
	"kroetnet/msg"
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
	response := []byte(" Alive!")
	reponseClient(pc, addr, response)
}

func handleTimeReq(pc net.PacketConn, addr net.Addr, buf []byte,
	recvTime time.Time) {
	timeResp := msg.TimeSyncRespMsg{
		MessageID:                   msg.TimeRespMsgID,
		TransmissionTimestamp:       binary.LittleEndian.Uint64(buf[1:]),
		ServerReceptionTimestamp:    uint64(recvTime.UnixNano() / 100),
		ServerTransmissionTimestamp: uint64(time.Now().UnixNano() / 100)}

	// nano seconcs / 100 == ticks
	rsp := timeResp.Encode()
	// send data
	reponseClient(pc, addr, rsp)
}

func handleTimeSyncDone(pc net.PacketConn, addr net.Addr, buf []byte, playerID int) {
	timesyncdoneackmsg := msg.TimeSyncDoneAckMsg{MessageID: msg.TimeSyncDoneAckMsgID, PlayerID: byte(playerID)}
	reponseClient(pc, addr, timesyncdoneackmsg.Encode())
}

func sendGameStart(pc net.PacketConn, addr net.Addr) {
	matchstart := msg.MatchStartMsg{MessageID: msg.MatchStartMsgID,
		MatchStartTimestamp: uint64(time.Now().UnixNano()/1000000 + 1000)}
	// ts is in ms and match start in now + 1 second
	reponseClient(pc, addr, matchstart.Encode())
}

func sendGameEnd(pc net.PacketConn, addr net.Addr) {
	matchendmsg := msg.MatchEndMsg{MessageID: msg.MatchEndMsgID}
	reponseClient(pc, addr, matchendmsg.Encode())
}

func handleInputMsg(pc net.PacketConn, addr net.Addr, buf []byte) {
	inputmsg := msg.DecodeInputMsg(buf)
	// validate moves
	log.Println("Pkg Received: ", inputmsg)
	rsp := inputmsg.Encode()
	// response with position (X,Y)
	reponseClient(pc, addr, rsp)
}
