package main

import (
	"kroetnet/msg"
	"log"
	"math"
	"net"
	"time"
)

// RcvPkt is a received by ListenUDP
type RcvPkt struct {
	connection *net.UDPConn
	addr       net.Addr
	buffer     []byte
}

// OutPkt ...
type OutPkt struct {
	connection net.PacketConn
	addr       net.Addr
	buffer     []byte
}

// Game holds current game properties
type Game struct {
	players              []Player
	State                int
	Port                 string
	recvCh               chan *RcvPkt
	sendCh               chan *OutPkt
	Frame                byte
	StateChangeTimestamp int64
	start                time.Time
	end                  time.Time
	statesMap            []Queue
}

func newGame(playerCount, playerStateQueueCount int, port string) *Game {
	return &Game{
		State:     0,
		Port:      port,
		players:   make([]Player, playerCount),
		statesMap: make([]Queue, playerStateQueueCount),
		recvCh:    make(chan *RcvPkt),
		sendCh:    make(chan *OutPkt),
	}

}

// Game server startup routines
func (g *Game) startServer() {
	go g.receiveUDP()
	go g.processUDP()
	g.sendByteResponse()
	log.Println("Started game server")
}

func (g *Game) receiveUDP() {
	udpAddr, err := net.ResolveUDPAddr("udp", g.Port)
	if err != nil {
		panic(err)
	}
	network := "udp"
	pc, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("listening on (%s)%s\n", network, pc.LocalAddr())
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
		g.recvCh <- &RcvPkt{pc, addr, buf[:n]}

		g.checkStateDuration(pc, addr)
	}

}

func (g *Game) checkStateDuration(pc net.PacketConn, addr net.Addr) {
	// if no ack is received for 2 seconds
	if time.Now().Unix()-g.StateChangeTimestamp > 2 {
		if g.State == 1 {
			// rollback to timesync state
			g.State--
		} else if g.State == 3 {
			// rollback to input/game-end-reached state
			g.State--
		}
	}
	// if game.State == 2 && (time.Now().After(game.end)) {
	//   // game is over
	//   sendGameEnd(pc, addr)
	//   game.State = 3
	//   // game.end = time.Now().Unix() + 400
	// }
}

func (g *Game) incFrame(t time.Time) {
	if g.State == 2 {
		// log.Printf("Frame updated at %v", t)
		//fmt.Printf("Frame: %v at Time: %v \n", game.Frame, t.UnixNano()/1000000)
		// calculating the frame based on the match start protects from frame drift, when this function invoked slightly earlier or delayed.
		msSinceStart := time.Now().Sub(g.start).Nanoseconds() / 1000000
		g.Frame = byte(math.Mod(float64(msSinceStart/33), 255.))
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func (g *Game) processUDP() {
	for v := range g.recvCh {
		pc := v.connection
		addr := v.addr
		buf := v.buffer
		recvTime := time.Now()
		log.Println("received buffer", buf)
		msgID := buf[0]
		switch g.State {
		case 0:
			if msgID == msg.TimeReqMsgID {
				g.handleTimeReq(pc, addr, buf, recvTime)
			} else if msgID == msg.TimeSyncDoneMsgID {
				if playerID := g.AddPlayer(addr); playerID != -1 {
					g.handleTimeSyncDone(pc, addr, buf, playerID)
					g.CheckGameFull(pc, addr)
				}
			}
		case 1:
			if msgID == msg.MatchStartAckMsgID {
				// todo mark player start acked, of not send match start again
			}
		case 2:
			// handle inputs until game end
			if msgID == msg.InputMsgID {
				g.handleInputMsg(pc, addr, buf)
			}
		case 3:
			if msgID == msg.MatchEndAckMsgID {
				log.Println("GAME FINISHED")
			}
		}
	}
}

func (g *Game) sendByteResponse() {
	for v := range g.sendCh {
		reponseClient(v.connection, v.addr, v.buffer)
	}
}

var emptyPlayer = Player{}

// AddPlayer adds servers to the game
func (g *Game) AddPlayer(addr net.Addr) int {
	playerID := -1
	nextPlayerID := 0
	for i := 0; i < len(g.players); i++ {
		if g.players[i] != emptyPlayer {
			nextPlayerID = i + 1
			// find player with same addr from udp packet
			if g.players[i].ipAddr == addr {
				playerID = g.players[i].id
				break
			}
		}
	}
	// player no in match yet & match not full
	if playerID == -1 && g.players[len(g.players)-1] == emptyPlayer {
		g.players[nextPlayerID] = Player{id: nextPlayerID, ipAddr: addr}
		playerID = nextPlayerID
		g.statesMap[playerID] = *NewQueue(15)
		return playerID
	}
	return -1
}

// CheckGameFull changes the gamestate when all players joined
func (g *Game) CheckGameFull(pc net.PacketConn, addr net.Addr) {
	// last player joined and match is full
	if g.players[len(g.players)-1] != emptyPlayer {
		// wait for all players
		for _, v := range g.players {
			g.sendGameStart(pc, v.ipAddr)
		}
		// skip state 1 for now
		time.Sleep(time.Second)
		g.State = 2
		// frame tick every 33 ms
		go doEvery(33*time.Millisecond, g.incFrame)
		log.Println("GAME STARTED")

		g.start = time.Now()
		// game ends after 2 minutes
		g.end = time.Now().Add(time.Minute * 1)
		// game.State = 1
		// game.StateChangeTimestamp = time.Now().Unix()
	}
}

func (g *Game) handleInputMsg(pc net.PacketConn, addr net.Addr, buf []byte) {
	inputmsg := msg.DecodeInputMsg(buf)
	// log.Println("Pkg Received: ", inputmsg)
	unitstatemsg := msg.UnitStateMsg{}
	for k, v := range g.players {
		if byte(v.id) == inputmsg.PlayerID {

			// send old unitstate
			// if game.statesMap[v.id].count > 14 {
			//   oldState := game.statesMap[v.id].Pop()
			//   oldUnitStateMsg := msg.UnitStateMsg{
			//     MessageID: msg.PositionConfirmationMessage,
			//     UnitID:    byte(v.id),
			//     XPosition: oldState.Xpos,
			//     YPosition: oldState.Ypos,
			//     Rotation:  0,
			//     Frame:     oldState.Frame}
			//   log.Println("POP Ele: ", &oldState)
			//   // log.Println("After POP QUEUE: ", game.statesMap[v.id].nodes)
			//   reponseClient(pc, addr, oldUnitStateMsg.Encode())
			// }

			// validate move
			newX, newY := v.move(inputmsg)
			// log.Println("PLAYER STATE", v.Y, v.X)
			g.players[k].X, g.players[k].Y = newX, newY
			// confirmation for player
			resp := msg.PositionConfirmationMsg{
				MessageID: msg.PositionConfirmationMessageID,
				UnitID:    byte(v.id),
				XPosition: newX,
				YPosition: newY,
				Frame:     g.Frame}

			// players state for all other clients
			unitstatemsg = msg.UnitStateMsg{
				MessageID: msg.UnitStateMsgID,
				UnitID:    byte(v.id),
				XPosition: newX,
				YPosition: newY,
				Rotation:  v.rotation,
				Frame:     g.Frame}

			// g.statesMap[v.id].Push(&PastState{byte(game.Frame), newX, newY,
			//   inputmsg.XTranslation, inputmsg.YTranslation})
			reponseClient(pc, addr, resp.Encode())
			// log.Println("After PUSH QUEUE: ", game.statesMap[v.id].nodes)

			// validate and update past moves
			// validateAllStates(v)
		}
	}
	// unitstate for all players
	for _, v := range g.players {
		if v.ipAddr != addr {
			reponseClient(pc, v.ipAddr, unitstatemsg.Encode())
		}
	}
}

// func validateAllStates(v Player) {
//   log.Println(game.statesMap[v.id])
//   if len(game.statesMap[v.id].nodes) == 0 {
//     return
//   }
//   sort.Slice(game.statesMap[v.id], func(i, j int) bool {
//     return game.statesMap[v.id].nodes[i].Frame <
//       game.statesMap[v.id].nodes[j].Frame
//   })
//   inpMsgArr := []msg.InputMsg{}
//   for i := 0; i < len(game.statesMap[v.id].nodes)-1; i++ {
//     ps := game.statesMap[v.id].nodes[i]
//     inpMsgArr = append(inpMsgArr,
//       msg.InputMsg{MessageID: 0,
//         PlayerID: byte(v.id), XTranslation: ps.Xtans,
//         YTranslation: ps.Ytrans, Frame: ps.Frame})
//
//   }
//   x, y := v.validateMoves(inpMsgArr[:len(inpMsgArr)-2])
//   game.statesMap[v.id].nodes[13].Xpos = x
//   game.statesMap[v.id].nodes[13].Ypos = y
// }
