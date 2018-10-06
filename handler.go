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
		log.Fatalln("err sending data for msgID :", buf[0], err)
	}
}

func (g *Game) serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	response := []byte(" Alive!")
	g.sendCh <- &OutPkt{pc, addr, response}
}

func (g *Game) handleTimeReq(pc net.PacketConn, addr net.Addr, buf []byte,
	recvTime time.Time) {
	timeResp := msg.TimeSyncRespMsg{
		MessageID:                   msg.TimeRespMsgID,
		TransmissionTimestamp:       binary.LittleEndian.Uint64(buf[1:]),
		ServerReceptionTimestamp:    uint64(recvTime.UnixNano() / 100),
		ServerTransmissionTimestamp: uint64(time.Now().UnixNano() / 100)}

	// nano seconcs / 100 == ticks
	rsp := timeResp.Encode()
	g.sendCh <- &OutPkt{pc, addr, rsp}
}

func (g *Game) handleTimeSyncDone(pc net.PacketConn, addr net.Addr, buf []byte, playerID int) {
	timesyncdoneackmsg := msg.TimeSyncDoneAckMsg{MessageID: msg.TimeSyncDoneAckMsgID, PlayerID: byte(playerID)}
	g.sendCh <- &OutPkt{pc, addr, timesyncdoneackmsg.Encode()}
}

func (g *Game) sendGameStart(pc net.PacketConn, addr net.Addr) {
	matchstart := msg.MatchStartMsg{MessageID: msg.MatchStartMsgID,
		MatchStartTimestamp: uint64(time.Now().UnixNano()/1000000 + 1000)}
	// ts is in ms and match start in now + 1 second
	g.sendCh <- &OutPkt{pc, addr, matchstart.Encode()}
}

func (g *Game) sendGameEnd(pc net.PacketConn, addr net.Addr) {
	matchendmsg := msg.MatchEndMsg{MessageID: msg.MatchEndMsgID}
	g.sendCh <- &OutPkt{pc, addr, matchendmsg.Encode()}
}
