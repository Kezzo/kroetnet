package main

import (
	"encoding/binary"
	"kroetnet/msg"
	"log"
	"net"
	"sort"
	"time"
)

func reponseClient(pc net.PacketConn, addr net.Addr, buf []byte) {
	// log.Println("Reponse to send: ", buf, " to ", addr)
	if _, err := pc.WriteTo(buf, addr); err != nil {
		log.Fatalln("err sending data for msgID :", buf[0], err)
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
	log.Println("Pkg Received: ", inputmsg)
	resp := msg.UnitStateMsg{}
	for k, v := range game.players {
		if game.players[k].ipAddr == v.ipAddr {
			// send old state
			if game.statesMap[v.id].count > 14 {
				oldState := game.statesMap[v.id].Pop()
				oldUnitStateMsg := msg.UnitStateMsg{
					MessageID: msg.PositionConfirmationMessage,
					UnitID:    byte(v.id),
					XPosition: oldState.Xpos,
					YPosition: oldState.Ypos,
					Rotation:  0,
					Frame:     oldState.Frame}
				log.Println("POP Ele: ", &oldState)
				// log.Println("After POP QUEUE: ", game.statesMap[v.id].nodes)
				reponseClient(pc, addr, oldUnitStateMsg.Encode())
			}

			// validate move
			newX, newY := v.move(inputmsg)
			log.Println("NEW MOVE", newX, newY)
			resp = msg.UnitStateMsg{
				MessageID: msg.UnitStateMsgID,
				UnitID:    byte(v.id),
				XPosition: newX,
				YPosition: newY,
				Rotation:  v.rotation,
				Frame:     byte(game.Frame)}
			game.statesMap[v.id].Push(&PastState{byte(game.Frame), newX, newY,
				inputmsg.XTranslation, inputmsg.YTranslation})
			log.Println("After PUSH QUEUE: ", game.statesMap[v.id].nodes)

			// validate and update past moves
			// validateAllStates(v)
		}
	}
	// unitstate for all players
	for _, v := range game.players {
		reponseClient(pc, v.ipAddr, resp.Encode())
	}
}

func validateAllStates(v Player) {
	log.Println(game.statesMap[v.id])
	if len(game.statesMap[v.id].nodes) == 0 {
		return
	}
	sort.Slice(game.statesMap[v.id], func(i, j int) bool {
		return game.statesMap[v.id].nodes[i].Frame <
			game.statesMap[v.id].nodes[j].Frame
	})
	inpMsgArr := []msg.InputMsg{}
	for i := 0; i < len(game.statesMap[v.id].nodes)-1; i++ {
		ps := game.statesMap[v.id].nodes[i]
		inpMsgArr = append(inpMsgArr,
			msg.InputMsg{MessageID: 0,
				PlayerID: byte(v.id), XTranslation: ps.Xtans,
				YTranslation: ps.Ytrans, Rotation: 0, Frame: ps.Frame})

	}
	x, y := v.validateMoves(inpMsgArr[:len(inpMsgArr)-2])
	game.statesMap[v.id].nodes[13].Xpos = x
	game.statesMap[v.id].nodes[13].Ypos = y
}
