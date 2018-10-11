package main

import (
	"bytes"
	"encoding/binary"
	"kroetnet/msg"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Game holds current game properties
type Game struct {
	players              []Player
	playerCount          int
	State                int
	Frame                byte
	StateChangeTimestamp int64
	recvCountMap         []bool
	pendingInputMsgs     []msg.InputMsg
	start                time.Time
	end                  time.Time
	playerStateQueue     []Queue
	network              Network
}

func newGame(playerCount, playerStateQueueCount int, port string) *Game {
	return &Game{
		players:              make([]Player, 0, playerCount),
		playerCount:          playerCount,
		playerStateQueue:     make([]Queue, playerStateQueueCount),
		StateChangeTimestamp: time.Now().Add(time.Second * 15).Unix(),
		network:              *newNetwork(port),
		recvCountMap:         make([]bool, playerCount),
		pendingInputMsgs:     make([]msg.InputMsg, 0, playerCount)}
}

func (g *Game) registerGameServer() {
	// test case
	port := g.network.Port
	count := strconv.Itoa(g.playerCount)
	jsonStr := []byte(`{"port":"` + port + `", "playerCount":` + count + `}`)
	log.Println("JSON: ", string(jsonStr))
	resp, err := http.Post(os.Getenv("WEBSERVER_ADDR")+"matchserver",
		"application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Register Game Server result: ", resp)
}

// Game server startup routines
func (g *Game) startServer() {
	go g.network.listenUDP()
	if os.Getenv("GO_ENV") == "DEV" {
		time.Sleep(2 * time.Second)
		g.registerGameServer()
	}
	go g.processMessages()
	go g.network.sendByteResponse()
	log.Println("Started match server")
	for {
		g.checkStateDuration()
	}
}

func (g *Game) checkStateDuration() {
	// if no ack is received for 5 seconds
	if time.Now().Unix()-g.StateChangeTimestamp > 5 {
		if g.State == 1 {
			// rollback to timesync state
			log.Println("ROLLBACK from State 1 to 0")
			g.State--
			// reset players joined
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
		for _, playerData := range g.players {
			lastState := g.playerStateQueue[playerData.id].nodes[len(g.playerStateQueue[playerData.id].nodes)-1]
			nextPosX, nextPosY := int32(0), int32(0)
			if lastState != nil {
				nextPosX, nextPosY = GetPosition(lastState.Xpos, lastState.Ypos, lastState.Xtrans, lastState.Ytrans)
			}

			g.playerStateQueue[playerData.id].Push(&PastState{byte(g.Frame), nextPosX, nextPosY, 127, 127})
		}

		// the input msgs need to be processed after the frame has been increased to be able to consider input msgs that arrived shortly before
		g.processPendingInputMsgs(g.network.connecton)
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

		//log.Println("Received buffer: ", buf)

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
	for i := 0; i < len(g.players); i++ {
		playerToCheck := g.players[i].ipAddr.String()
		incomingAddr := addr.String()
		if playerToCheck == incomingAddr {
			return -1
		}
	}

	// player not in match yet & match not full
	g.players = append(g.players, Player{id: len(g.players), ipAddr: addr})
	g.playerStateQueue[len(g.players)-1] = *NewQueue(15)
	return len(g.players) - 1
}

// CheckGameFull changes the gamestate when all players joined
func (g *Game) CheckGameFull(pc net.PacketConn, addr net.Addr) {
	if len(g.players) == g.playerCount {
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
		g.end = time.Now().Add(time.Second * 300)
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

func (g *Game) handleTimeSyncDone(pc net.PacketConn, addr net.Addr, buf []byte,
	playerID int) {
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
	//log.Println("Pkg Received: ", inputmsg)
	g.pendingInputMsgs = append(g.pendingInputMsgs, inputmsg)
}

func (g *Game) processPendingInputMsgs(pc net.PacketConn) {
	for _, inputmsg := range g.pendingInputMsgs {
		unitstatemsg := msg.UnitStateMsg{}
		for playerID, playerData := range g.players {
			if byte(playerData.id) == inputmsg.PlayerID {

				foundFrame := false
				for _, pastState := range g.playerStateQueue[playerData.id].nodes {
					if pastState != nil {
						// set translation to past frame
						if inputmsg.Frame == pastState.Frame {
							pastState.Xtrans = inputmsg.XTranslation
							pastState.Ytrans = inputmsg.YTranslation
							foundFrame = true
						}
					}
				}

				if !foundFrame {
					frames := make([]byte, 0)
					for _, pastState := range g.playerStateQueue[playerData.id].nodes {
						frames = append(frames, pastState.Frame)
					}

					log.Println("Did not find frame: ", inputmsg.Frame, " Frames: ", frames)
				}

				// calculate/validate all movements from predecessor
				for index, pastState := range g.playerStateQueue[playerData.id].nodes {
					if index > 0 {
						prevXPos := g.playerStateQueue[playerData.id].nodes[index-1].Xpos
						prevYPos := g.playerStateQueue[playerData.id].nodes[index-1].Ypos
						prevXTrans := g.playerStateQueue[playerData.id].nodes[index-1].Xtrans
						prevYTrans := g.playerStateQueue[playerData.id].nodes[index-1].Ytrans
						updXPos, updYPos := GetPosition(prevXPos, prevYPos, prevXTrans, prevYTrans)
						pastState.Xpos, pastState.Ypos = updXPos, updYPos
					}
				}

				// validate move
				newX, newY := GetPosition(g.playerStateQueue[playerData.id].nodes[len(g.playerStateQueue[playerData.id].nodes)-1].Xpos,
					g.playerStateQueue[playerData.id].nodes[len(g.playerStateQueue[playerData.id].nodes)-1].Ypos,
					g.playerStateQueue[playerData.id].nodes[len(g.playerStateQueue[playerData.id].nodes)-1].Xtrans,
					g.playerStateQueue[playerData.id].nodes[len(g.playerStateQueue[playerData.id].nodes)-1].Ytrans)
				playerData.rotation = inputmsg.Rotation
				// log.Println("PLAYER STATE", v.Y, v.X)
				g.players[playerID].X, g.players[playerID].Y = newX, newY

				// players state for all other clients
				unitstatemsg = msg.UnitStateMsg{
					MessageID: msg.UnitStateMsgID,
					UnitID:    byte(playerData.id),
					XPosition: newX,
					YPosition: newY,
					Rotation:  playerData.rotation,
					Frame:     g.Frame}

				if len(g.playerStateQueue[playerData.id].nodes) > 14 {
					oldState := g.playerStateQueue[playerData.id].nodes[0]
					oldXPos, oldYPos := GetPosition(oldState.Xpos, oldState.Ypos, oldState.Xtrans, oldState.Ytrans)

					resp := msg.PositionConfirmationMsg{
						MessageID: msg.PositionConfirmationMessageID,
						UnitID:    byte(playerData.id),
						XPosition: oldXPos,
						YPosition: oldYPos,
						Frame:     oldState.Frame}

					for _, v := range g.players {
						if v.id == int(inputmsg.PlayerID) {
							g.network.sendCh <- &OutPkt{pc, v.ipAddr, resp.Encode()}
						}
					}
				}

			}
		}
		// unitstate for all clients
		for _, v := range g.players {
			if v.id != int(inputmsg.PlayerID) {
				g.network.sendCh <- &OutPkt{pc, v.ipAddr, unitstatemsg.Encode()}
			}
		}
	}

	// clear pending input msgs
	g.pendingInputMsgs = make([]msg.InputMsg, 0, len(g.players))
}
