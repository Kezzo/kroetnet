package main

import (
	"fmt"
	"kroetnet/msg"
	"log"
	"net"
	"time"
)

var game = Game{State: 0}

func main() {

	port := ":2448"
	network := "udp"
	pc, err := net.ListenPacket(network, port)
	handleError(err)
	fmt.Printf("listening on (%s)%s\n", network, pc.LocalAddr())
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		// fmt.Printf("\nBuffer Content: [ % x ] \n", buf[:n])

		if err != nil {
			log.Print("Error: ", err)
			continue
		}
		checkStateDuration(pc, addr)
		go digestPacket(pc, addr, buf[:n])
	}
}

func checkStateDuration(pc net.PacketConn, addr net.Addr) {
	// if no ack is received for 2 seconds
	if time.Now().Unix()-game.StateChangeTimestamp > 2 {
		if game.State == 1 {
			// rollback to timesync state
			game.State--
		} else if game.State == 2 && game.end > time.Now().Unix() {
			// game is over
			sendGameEnd(pc, addr)
			game.State = 3
			game.end = time.Now().Unix() + 40
		} else if game.State == 3 {
			// rollback to input/game-end-reached state
			game.State--
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatalln("Error: ", err)
	}
}

func digestPacket(pc net.PacketConn, addr net.Addr, buf []byte) {
	recvTime := time.Now()
	log.Println("received buffer", buf)

	msgID := buf[0]
	switch game.State {
	case 0:
		if msgID == msg.TimeReqMsgID {
			handleTimeReq(pc, addr, buf, recvTime)
		} else if msgID == msg.TimeSyncDoneMsgID {
			handleTimeSyncDone(pc, addr, buf)
			// wait for all players
			sendGameStart(pc, addr)
			game.State = 1
			game.StateChangeTimestamp = time.Now().Unix()
		}
	case 1:
		if msgID == msg.MatchStartAckMsgID {
			// check if every ack was received
			game.State = 2
			game.end = time.Now().Unix() + 30
		}
	case 2:
		// hande inputs until game end
		if msgID == msg.InputMsgID {
			handleInputMsg(pc, addr, buf)
		}
	case 3:
		if msgID == msg.MatchEndAckMsgID {
			log.Println("GAME FINISHED")
		}
	}

}
