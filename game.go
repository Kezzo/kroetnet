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
	recvCountMap         []bool
	start                time.Time
	end                  time.Time
	playerStateQueue     []Queue
	network              Network
}

func newGame(playerCount, playerStateQueueCount int, port string) *Game {
	return &Game{
		players:              make([]Player, playerCount),
		playerStateQueue:     make([]Queue, playerStateQueueCount),
		StateChangeTimestamp: time.Now().Add(time.Second * 15).Unix(),
		network:              *newNetwork(port),
		recvCountMap:         make([]bool, playerCount),
	}

}

// Game server startup routines
func (g *Game) startServer() {
	go g.network.listenUDP()
	go g.processMessages()
	go g.network.sendByteResponse()
	log.Println("Started match server")
	for {
		g.checkStateDuration()
	}
}

func (g *Game) checkStateDuration() {
	// if no ack is received for 2 seconds
	if time.Now().Unix()-g.StateChangeTimestamp > 2 {
		if g.State == 1 {
			// rollback to timesync state
			log.Println("ROLLBACK from State 1 to 0")
			g.State--
		}
	}
	if g.State == 2 && (time.Now().After(g.end)) {
		for _, v := range g.players {
			log.Println("SEND MATCH END to Player ", v.id)
			matchendmsg := msg.MatchEndMsg{MessageID: msg.MatchEndMsgID}
			g.network.sendCh <- &OutPkt{g.network.connecton,
				v.ipAddr, matchendmsg.Encode()}
		}
		g.State = 3
	}
}

func (g *Game) incFrame(t time.Time) {
	if g.State == 1 || g.State == 2 {
		// log.Printf("Frame updated at %v", t)
		//fmt.Printf("Frame: %v at Time: %v \n", game.Frame, t.UnixNano()/1000000)
		// calculating the frame based on the match start protects from frame drift, when this function invoked slightly earlier or delayed.
		msSinceStart := time.Now().Sub(g.start).Nanoseconds() / 1000000
		g.Frame = byte(math.Mod(float64(msSinceStart/33), 255.))
		for k, v := range g.players {
			g.playerStateQueue[k].Push(&PastState{byte(g.Frame), v.X, v.Y, 0, 0})
		}
		// for _, v := range g.playerStateQueue {
		//   for _, v := range v.nodes {
		//     log.Println("QUEUE", v)
		//   }
		// }
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func (g *Game) incAckCounter(addr net.Addr) {
	for _, v := range g.players {
		if v.ipAddr == addr {
			g.recvCountMap[v.id] = true
		}
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
				g.incAckCounter(addr)
				if len(g.recvCountMap) == len(g.players) {
					log.Println("All Clients send MatchStartAck")
					g.State = 2
				}
			}
		case 2:
			// handle inputs until game end
			if msgID == msg.InputMsgID {
				g.handleInputMsg(pc, addr, buf)
			}
		case 3:
			if msgID == msg.MatchEndAckMsgID {
				g.incAckCounter(addr)
				if len(g.recvCountMap) == len(g.players) {
					log.Println("GAME FINISHED, all clients send ACK")
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
	if g.players[len(g.players)-1] != emptyPlayer {
		for _, v := range g.players {
			g.sendGameStart(pc, v.ipAddr)
		}
		time.Sleep(time.Second)
		go doEvery(33*time.Millisecond, g.incFrame)
		g.State = 1
		g.StateChangeTimestamp = time.Now().Unix()
		g.recvCountMap = make([]bool, len(g.players))
		g.start = time.Now()
		// g.end = time.Now().Add(time.Minute * 1)
		g.end = time.Now().Add(time.Second * 30)
		log.Println("Server started match")
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

			// validate move
			newX, newY := v.move(inputmsg)
			g.players[k].X, g.players[k].Y = newX, newY

			// new position for actual frame and input
			for _, v := range g.playerStateQueue[v.id].nodes {
				if inputmsg.Frame == v.Frame {
					v.Xpos = newX
					v.Ypos = newY
				}
			}

			// add translation to the previous state in the queue
			for _, v := range g.playerStateQueue[v.id].nodes {
				if inputmsg.Frame-1 == v.Frame && v.Xtrans != 0 && v.Ytrans != 0 {
					v.Xtrans = inputmsg.XTranslation
					v.Ytrans = inputmsg.YTranslation
				}
			}

			// calculate all movements from the queued states
			for k, val := range g.playerStateQueue[v.id].nodes {
				tmpInput := msg.InputMsg{}
				tmpPlayer := Player{X: val.Xpos, Y: val.Ypos}
				if k > 0 {
					tmpInput = msg.InputMsg{
						XTranslation: g.playerStateQueue[v.id].nodes[k-1].Xtrans,
						YTranslation: g.playerStateQueue[v.id].nodes[k-1].Ytrans}
					if g.playerStateQueue[v.id].nodes[k-1].Xtrans != 0 &&
						g.playerStateQueue[v.id].nodes[k-1].Ytrans != 0 {
						X, Y := tmpPlayer.move(tmpInput)
						val.Xpos = X
						val.Xpos = Y
					}
				}
			}

			// players state for all other clients
			unitstatemsg = msg.UnitStateMsg{
				MessageID: msg.UnitStateMsgID,
				UnitID:    byte(v.id),
				XPosition: g.playerStateQueue[v.id].nodes[len(g.playerStateQueue[v.id].nodes)-1].Xpos,
				YPosition: g.playerStateQueue[v.id].nodes[len(g.playerStateQueue[v.id].nodes)-1].Ypos,
				Rotation:  v.rotation,
				Frame:     g.Frame}

			if len(g.playerStateQueue[v.id].nodes) > 14 {
				oldState := g.playerStateQueue[v.id].nodes[0]
				resp := msg.PositionConfirmationMsg{
					MessageID: msg.PositionConfirmationMessageID,
					UnitID:    byte(v.id),
					XPosition: oldState.Xpos,
					YPosition: oldState.Ypos,
					Frame:     oldState.Frame}
				g.network.sendCh <- &OutPkt{pc, addr, resp.Encode()}
			}

		}
	}
	// unitstate for all clients
	for _, v := range g.players {
		if v.ipAddr != addr {
			g.network.sendCh <- &OutPkt{pc, v.ipAddr, unitstatemsg.Encode()}
		}
	}
}
