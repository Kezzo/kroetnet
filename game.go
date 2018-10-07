package main

import (
	"encoding/binary"
	"kroetnet/msg"
	"log"
	"math"
	"net"
	"os"
	"time"
)

// Game holds current game properties
type Game struct {
	players              []Player
	State                int
	Frame                byte
	StateChangeTimestamp int64
	recvCount            int
	start                time.Time
	end                  time.Time
	playerStateQueue     []Queue
	network              Network
}

func newGame(playerCount, playerStateQueueCount int, port string) *Game {
	return &Game{
		players:          make([]Player, playerCount),
		playerStateQueue: make([]Queue, playerStateQueueCount),
		network:          *newNetwork(port),
	}

}

// Game server startup routines
func (g *Game) startServer() {
	go g.network.listenUDP()
	go g.processMessages()
	go g.network.sendByteResponse()
	log.Println("Started game server")
	for {
		g.checkStateDuration()
	}
}

func (g *Game) checkStateDuration() {
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
	if g.State == 2 && (time.Now().After(g.end)) {
		g.State = 3
		for _, v := range g.players {
			matchendmsg := msg.MatchEndMsg{MessageID: msg.MatchEndMsgID}
			g.network.sendCh <- &OutPkt{g.network.connecton,
				v.ipAddr, matchendmsg.Encode()}
		}
	}
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

func (g *Game) processMessages() {
	for v := range g.network.recvCh {
		pc := v.connection
		addr := v.addr
		buf := v.buffer
		recvTime := time.Now()
		msgID := buf[0]

		if msgID == msg.PingMsgID {
			g.handlePing(pc, addr, buf)
			continue
		}

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
				g.recvCount++
				if g.recvCount == len(g.players) {
					// todo mark player start acked, if not send match start again
				}
			}
		case 2:
			// handle inputs until game end
			if msgID == msg.InputMsgID {
				g.handleInputMsg(pc, addr, buf)
			}
		case 3:
			if msgID == msg.MatchEndAckMsgID {
				g.recvCount++
				if g.recvCount == len(g.players) {
					log.Println("GAME FINISHED")
					os.Exit(0)
				}
			}
		default:
			log.Println("Received invalid message :", buf, " from ", addr)
		}
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
	// player not in match yet & match not full
	if playerID == -1 && g.players[len(g.players)-1] == emptyPlayer {
		g.players[nextPlayerID] = Player{id: nextPlayerID, ipAddr: addr}
		playerID = nextPlayerID
		g.playerStateQueue[playerID] = *NewQueue(15)
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
		g.recvCount = 0
		g.start = time.Now()
		g.end = time.Now().Add(time.Minute * 1)
	}
}

func (g *Game) serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	response := []byte(" Alive!")
	g.network.sendCh <- &OutPkt{pc, addr, response}
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
	g.network.sendCh <- &OutPkt{pc, addr, rsp}
}

func (g *Game) handleTimeSyncDone(pc net.PacketConn, addr net.Addr, buf []byte, playerID int) {
	timesyncdoneackmsg := msg.TimeSyncDoneAckMsg{MessageID: msg.TimeSyncDoneAckMsgID, PlayerID: byte(playerID)}
	g.network.sendCh <- &OutPkt{pc, addr, timesyncdoneackmsg.Encode()}
}

func (g *Game) handlePing(pc net.PacketConn, addr net.Addr, buf []byte) {
	pongMsg := msg.PongMsg{
		MessageID:             msg.PongMsgID,
		TransmissionTimestamp: binary.LittleEndian.Uint64(buf[1:])}
	rsp := pongMsg.Encode()
	g.network.sendCh <- &OutPkt{pc, addr, rsp}
}

func (g *Game) sendGameStart(pc net.PacketConn, addr net.Addr) {
	matchstart := msg.MatchStartMsg{MessageID: msg.MatchStartMsgID,
		MatchStartTimestamp: uint64(time.Now().UnixNano()/1000000 + 1000)}
	// ts is in ms and match start in now + 1 second
	g.network.sendCh <- &OutPkt{pc, addr, matchstart.Encode()}
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
