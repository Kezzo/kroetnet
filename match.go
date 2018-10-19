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

// Match ...
type Match struct {
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
	endMatch             bool
}

func newMatch(playerCount, playerStateQueueCount int, port string) *Match {
	return &Match{
		players:              make([]Player, 0, playerCount),
		playerCount:          playerCount,
		playerStateQueue:     make([]Queue, playerStateQueueCount),
		StateChangeTimestamp: time.Now().Add(time.Second * 15).Unix(),
		network:              *newNetwork(port),
		recvCountMap:         make([]bool, playerCount),
		pendingInputMsgs:     make([]msg.InputMsg, 0, playerCount)}
}

func (m *Match) registerMatchServer() {
	// test case
	port := m.network.Port
	count := strconv.Itoa(m.playerCount)
	jsonStr := []byte(`{"port":"` + port + `", "playerCount":` + count + `}`)
	log.Println("JSON: ", string(jsonStr))
	resp, err := http.Post(os.Getenv("WEBSERVER_ADDR")+"matchserver",
		"application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Register Match Server result: ", resp)
}

// Match server startup routines
func (m *Match) startServer() {
	go m.network.listenUDP()
	if os.Getenv("GO_ENV") == "DEV" {
		time.Sleep(2 * time.Second)
		m.registerMatchServer()
	}
	go m.processMessages()
	go m.network.sendByteResponse()
	log.Println("Started match server")

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for {
		<-ticker.C
		m.checkStateDuration()
		if m.endMatch {
			break
		}
	}
}

func (m *Match) checkStateDuration() {
	if m.State > 0 {
		pingCounter := 0
		for _, playerData := range m.players {
			if (time.Now().Unix() - playerData.lastMsg.Unix()) > 5 {
				// log.Println("Player with following id timed out: ", playerData.id)
				pingCounter++
			}
		}
		if pingCounter == m.playerCount {
			log.Println("All Players timed out -  EXIT")
			os.Exit(0)
		}
	}

	// if no ack is received for 5 seconds
	if time.Now().Unix()-m.StateChangeTimestamp > 5 {
		if m.State == 1 {
			// rollback to timesync state
			log.Println("ROLLBACK from State 1 to 0")
			m.State--
			// reset players joined
		}
	}
	if m.State == 2 && (time.Now().After(m.end)) {
		for _, v := range m.players {
			log.Println("SEND MATCH END to Player ", v.id)
			matchendmsg := msg.MatchEndMsg{MessageID: msg.MatchEndMsgID}
			m.network.sendCh <- &OutPkt{m.network.connecton,
				v.ipAddr, matchendmsg.Encode()}
		}
		m.State = 3
	}
}

func (m *Match) incFrame(t time.Time) {
	if m.State == 1 || m.State == 2 {
		// log.Printf("Frame updated at %v", t)
		//fmt.Printf("Frame: %v at Time: %v \n", m.Frame, t.UnixNano()/1000000)
		// calculating the frame based on the match start protects from frame drift, when this function invoked slightly earlier or delayed.
		msSinceStart := time.Now().Sub(m.start).Nanoseconds() / 1000000
		currentFrame := byte(math.Mod(float64(msSinceStart/33), 255.))

		if m.Frame != currentFrame {
			//log.Println("Frame: ", currentFrame, " at Time: ", t.UnixNano()/1000000)

			for {
				m.Frame = byte(math.Mod(float64(m.Frame+1), 255.))
				//log.Println("Next frame: ", m.Frame)
				for _, playerData := range m.players {
					lastState := m.playerStateQueue[playerData.id].nodes[len(m.playerStateQueue[playerData.id].nodes)-1]
					nextPosX, nextPosY := int32(0), int32(0)
					if lastState != nil {
						nextPosX, nextPosY = GetPosition(lastState.Xpos, lastState.Ypos, lastState.Xtrans, lastState.Ytrans)
					}

					m.playerStateQueue[playerData.id].Push(&PastState{byte(m.Frame), nextPosX, nextPosY, 127, 127})
				}

				if m.Frame == currentFrame {
					// the input msgs need to be processed after the frame has been increased to be able to consider input msgs that arrived shortly before
					m.processPendingInputMsgs(m.network.connecton)
					break
				}
			}
		}
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		t := <-ticker.C
		f(t)
	}
}

func (m *Match) incAckCounter(addr net.Addr) {
	for _, v := range m.players {
		if v.ipAddr == addr {
			m.recvCountMap[v.id] = true
		}
	}
}

func (m *Match) processMessages() {
	for v := range m.network.recvCh {
		pc := v.connection
		addr := v.addr
		buf := v.buffer
		recvTime := time.Now()
		msgID := buf[0]

		if msgID == msg.PingMsgID {
			m.handlePing(pc, addr, buf)
			continue
		}
		//log.Println("Received buffer: ", buf)
		switch m.State {
		case 0:
			if msgID == msg.TimeReqMsgID {
				m.handleTimeReq(pc, addr, buf, recvTime)
			} else if msgID == msg.TimeSyncDoneMsgID {
				if playerID := m.AddPlayer(addr); playerID != -1 {
					m.handleTimeSyncDone(pc, addr, buf, playerID)
					m.CheckMatchFull(pc, addr)
				}
			}
		case 1:
			if msgID == msg.MatchStartAckMsgID {
				m.incAckCounter(addr)
				if len(m.recvCountMap) == len(m.players) {
					log.Println("All Clients sent MatchStartAck")
					m.State = 2
				}
			}
		case 2:
			// handle inputs until match end
			if msgID == msg.InputMsgID {
				m.handleInputMsg(pc, addr, buf)
			}
		case 3:
			if msgID == msg.MatchEndAckMsgID {
				m.incAckCounter(addr)
				if len(m.recvCountMap) == len(m.players) {
					log.Println("MATCH FINISHED, all clients sent ACK")
					m.endMatch = true
				}
			}
		default:
			log.Println("Received invalid message :", buf, " from ", addr)
		}
	}
}

var emptyPlayer = Player{}

// AddPlayer adds servers to the match
func (m *Match) AddPlayer(addr net.Addr) int {
	for i := 0; i < len(m.players); i++ {
		playerToCheck := m.players[i].ipAddr.String()
		incomingAddr := addr.String()
		if playerToCheck == incomingAddr {
			return -1
		}
	}

	// player not in match yet & match not full
	m.players = append(m.players, Player{id: len(m.players), ipAddr: addr})
	m.playerStateQueue[len(m.players)-1] = *NewQueue(15)
	return len(m.players) - 1
}

// CheckMatchFull changes the matchstate when all players joined
func (m *Match) CheckMatchFull(pc net.PacketConn, addr net.Addr) {
	if len(m.players) == m.playerCount {
		for _, v := range m.players {
			m.sendMatchStart(pc, v.ipAddr)
		}
		time.Sleep(time.Second)
		go doEvery(33*time.Millisecond, m.incFrame)
		m.State = 1
		m.StateChangeTimestamp = time.Now().Unix()
		m.recvCountMap = make([]bool, len(m.players))
		m.start = time.Now()
		m.end = time.Now().Add(time.Second * 300)
		log.Println("Server started match")
	}
}

func (m *Match) serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	response := []byte(" Alive!")
	m.network.sendCh <- &OutPkt{pc, addr, response}
}

func (m *Match) handleTimeReq(pc net.PacketConn, addr net.Addr, buf []byte,
	recvTime time.Time) {
	timeResp := msg.TimeSyncRespMsg{
		MessageID:                   msg.TimeRespMsgID,
		TransmissionTimestamp:       binary.LittleEndian.Uint64(buf[1:]),
		ServerReceptionTimestamp:    uint64(recvTime.UnixNano() / 100),
		ServerTransmissionTimestamp: uint64(time.Now().UnixNano() / 100)}
	// nano seconcs / 100 == ticks
	rsp := timeResp.Encode()
	m.network.sendCh <- &OutPkt{pc, addr, rsp}
}

func (m *Match) handleTimeSyncDone(pc net.PacketConn, addr net.Addr, buf []byte,
	playerID int) {
	timesyncdoneackmsg := msg.TimeSyncDoneAckMsg{MessageID: msg.TimeSyncDoneAckMsgID, PlayerID: byte(playerID)}
	m.network.sendCh <- &OutPkt{pc, addr, timesyncdoneackmsg.Encode()}
}

func (m *Match) handlePing(pc net.PacketConn, addr net.Addr, buf []byte) {
	for playerID, playerData := range m.players {
		if playerData.ipAddr.String() == addr.String() {
			m.players[playerID].lastMsg = time.Now()
		}
	}
	pongMsg := msg.PongMsg{
		MessageID:             msg.PongMsgID,
		TransmissionTimestamp: binary.LittleEndian.Uint64(buf[1:])}
	rsp := pongMsg.Encode()
	m.network.sendCh <- &OutPkt{pc, addr, rsp}
}

func (m *Match) sendMatchStart(pc net.PacketConn, addr net.Addr) {
	matchstart := msg.MatchStartMsg{MessageID: msg.MatchStartMsgID,
		MatchStartTimestamp: uint64(time.Now().UnixNano()/1000000 + 1000)}
	// ts is in ms and match start in now + 1 second
	m.network.sendCh <- &OutPkt{pc, addr, matchstart.Encode()}
}

func (m *Match) handleInputMsg(pc net.PacketConn, addr net.Addr, buf []byte) {
	inputmsg := msg.DecodeInputMsg(buf)
	//log.Println("Pkg Received: ", inputmsg)
	m.pendingInputMsgs = append(m.pendingInputMsgs, inputmsg)
}

func (m *Match) processPendingInputMsgs(pc net.PacketConn) {
	updatedPlayerIDs := make([]int, 0, 2)
	for _, inputmsg := range m.pendingInputMsgs {
		for playerID, playerData := range m.players {
			if byte(playerData.id) == inputmsg.PlayerID {

				foundFrame := false
				for _, pastState := range m.playerStateQueue[playerData.id].nodes {
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
					for _, pastState := range m.playerStateQueue[playerData.id].nodes {
						frames = append(frames, pastState.Frame)
					}

					log.Println("Did not find frame: ", inputmsg.Frame, " Frames: ", frames)
				}

				// calculate/validate all movements from predecessor
				for index, pastState := range m.playerStateQueue[playerData.id].nodes {
					if index > 0 {
						queue := m.playerStateQueue[playerData.id]
						node := queue.nodes[index-1]

						// node hasn't been set yet.
						if node == nil {
							continue
						}

						prevXPos := node.Xpos
						prevYPos := node.Ypos
						prevXTrans := node.Xtrans
						prevYTrans := node.Ytrans

						updXPos, updYPos := GetPosition(prevXPos, prevYPos, prevXTrans, prevYTrans)
						pastState.Xpos, pastState.Ypos = updXPos, updYPos
					}
				}

				latestNode := m.playerStateQueue[playerData.id].nodes[len(m.playerStateQueue[playerData.id].nodes)-1]
				// validate move
				m.players[playerID].X, m.players[playerID].Y = GetPosition(
					latestNode.Xpos, latestNode.Ypos, latestNode.Xtrans, latestNode.Ytrans)
				m.players[playerID].rotation = inputmsg.Rotation

				addPlayerID := true
				for _, playerIDEntry := range updatedPlayerIDs {
					if playerIDEntry == playerID {
						addPlayerID = false
						break
					}
				}

				if addPlayerID {
					updatedPlayerIDs = append(updatedPlayerIDs, playerID)
				}
			}
		}
	}

	for _, playerID := range updatedPlayerIDs {
		playerData := m.players[playerID]
		// players state for all other clients
		unitstatemsg := msg.UnitStateMsg{
			MessageID: msg.UnitStateMsgID,
			UnitID:    byte(playerData.id),
			XPosition: playerData.X,
			YPosition: playerData.Y,
			Rotation:  playerData.rotation,
			Frame:     m.Frame}

		// unitstate for all clients
		for _, v := range m.players {
			if v.id != playerID {
				m.network.sendCh <- &OutPkt{pc, v.ipAddr, unitstatemsg.Encode()}
			}
		}

		if m.playerStateQueue[playerData.id].nodes[0] != nil {
			oldState := m.playerStateQueue[playerData.id].nodes[0]
			oldXPos, oldYPos := GetPosition(oldState.Xpos, oldState.Ypos, oldState.Xtrans, oldState.Ytrans)

			resp := msg.PositionConfirmationMsg{
				MessageID: msg.PositionConfirmationMessageID,
				UnitID:    byte(playerData.id),
				XPosition: oldXPos,
				YPosition: oldYPos,
				Frame:     oldState.Frame}

			//log.Println("Sending pcm with frame: ", resp.Frame)

			for _, v := range m.players {
				if v.id == playerID {
					m.network.sendCh <- &OutPkt{pc, v.ipAddr, resp.Encode()}
				}
			}
		}

		// clear pending input msgs
		m.pendingInputMsgs = make([]msg.InputMsg, 0, len(m.players))
	}
}
