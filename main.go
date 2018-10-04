package main

import (
	"fmt"
	"kroetnet/msg"
	"log"
	"math"
	"net"
	"time"
)

var game = Game{
	State:     0,
	players:   make([]Player, 2),
	statesMap: make([]Queue, 5)}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":2448")
	handleError(err)

	network := "udp"
	pc, err := net.ListenUDP(network, udpAddr)
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
		// fmt.Println("State is ", game.State)
		// fmt.Println("Frame is ", game.Frame)

		checkStateDuration(pc, addr)
		digestPacket(pc, addr, buf[:n])
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func incFrame(t time.Time) {
	if game.State == 2 {
		// log.Printf("Frame updated at %v", t)
		//fmt.Printf("Frame: %v at Time: %v \n", game.Frame, t.UnixNano()/1000000)
		game.Frame = byte(math.Mod(float64(game.Frame)+1., 255.))
	}
}

func checkStateDuration(pc net.PacketConn, addr net.Addr) {
	// if no ack is received for 2 seconds
	if time.Now().Unix()-game.StateChangeTimestamp > 2 {
		if game.State == 1 {
			// rollback to timesync state
			game.State--
		} else if game.State == 3 {
			// rollback to input/game-end-reached state
			game.State--
		}
	}
	// if game.State == 2 && (time.Now().After(game.end)) {
	//   // game is over
	//   sendGameEnd(pc, addr)
	//   game.State = 3
	//   // game.end = time.Now().Unix() + 400
	// }
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
			if playerID := AddPlayer(addr); playerID != -1 {
				handleTimeSyncDone(pc, addr, buf, playerID)
				CheckGameFull(pc, addr)
			}
		}
	case 1:
		if msgID == msg.MatchStartAckMsgID {
			// check if every ack was received
			game.State = 2
			log.Println("GAME STARTED")
			// game ends after 2 minutes
			game.end = time.Now().Add(time.Minute * 1)
		}
	case 2:
		// handle inputs until game end
		if msgID == msg.InputMsgID {
			handleInputMsg(pc, addr, buf)
		}
	case 3:
		if msgID == msg.MatchEndAckMsgID {
			log.Println("GAME FINISHED")
		}
	}

}
