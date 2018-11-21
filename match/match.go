package match

import (
	"encoding/binary"
	"kroetnet/abilities"
	"kroetnet/ai"
	"kroetnet/collision"
	"kroetnet/msg"
	"kroetnet/network"
	"kroetnet/units"
	"log"
	"math"
	"net"
	"os"
	"time"
)

// Match ...
type Match struct {
	players                 []*units.Player
	npcUnits                []*units.NPCUnit
	allUnits                []units.Unit
	aiList                  []ai.FireBoss
	playerCount             int
	State                   int
	Frame                   byte
	StateChangeTimestamp    int64
	recvCountMap            []bool
	pendingInputMsgs        []msg.InputMsg
	pendingAbilityInputMsgs []msg.AbilityInputMsg
	start                   time.Time
	end                     time.Time
	playerStateQueue        []Queue
	abilities               map[byte]abilities.Ability
	network                 network.Network
	endMatch                bool
}

// NewMatch ...
func NewMatch(playerCount, playerStateQueueCount int, port string) *Match {
	return &Match{
		players:                 make([]*units.Player, 0, playerCount),
		npcUnits:                make([]*units.NPCUnit, 0, playerCount),
		allUnits:                make([]units.Unit, 0, playerCount),
		aiList:                  make([]ai.FireBoss, 0, 0),
		playerCount:             playerCount,
		playerStateQueue:        make([]Queue, playerStateQueueCount),
		abilities:               make(map[byte]abilities.Ability),
		StateChangeTimestamp:    time.Now().Add(time.Second * 15).Unix(),
		network:                 *network.NewNetwork(port),
		recvCountMap:            make([]bool, playerCount),
		pendingInputMsgs:        make([]msg.InputMsg, 0, playerCount),
		pendingAbilityInputMsgs: make([]msg.AbilityInputMsg, 0, playerCount)}
}

// StartServer Match server startup routines
func (m *Match) StartServer() {
	go m.network.ListenUDP()
	go m.network.SendByteResponse()
	go m.processMessages()

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
			if (time.Now().Unix() - playerData.LastPingTime.Unix()) > 5 {
				log.Println("Player with following id timed out: ", playerData.ID, " current time: ", time.Now().Unix(), " last ping time: ", playerData.LastPingTime.Unix())
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
			// kill match is easier than rollback
			m.endMatch = true
		}
	}
	if m.State == 2 && (time.Now().After(m.end)) {
		for _, v := range m.players {
			// send 3 times
			for i := 0; i < 3; i++ {
				log.Println("SEND MATCH END to Player ", v.ID)
				matchendmsg := msg.MatchEndMsg{MessageID: msg.MatchEndMsgID}
				m.network.SendCh <- &network.OutPkt{Connection: m.network.Connection,
					Addr: v.IPAddr, Buffer: matchendmsg.Encode()}
			}
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

			/*if len(m.pendingInputMsgs) > 0 {
				log.Println("Frame: ", currentFrame, " at Time: ", t.UnixNano()/1000000, " inputmsgs: ", len(m.pendingInputMsgs))
			}*/

			for {
				m.Frame = byte(math.Mod(float64(m.Frame+1), 255.))
				//log.Println("Next frame: ", m.Frame)
				for _, playerData := range m.players {
					lastState := m.playerStateQueue[playerData.ID].nodes[len(m.playerStateQueue[playerData.ID].nodes)-1]
					nextPosX, nextPosY := int32(0), int32(0)
					if lastState != nil {
						nextPosX, nextPosY = units.GetPlayerPosition(lastState.Xpos, lastState.Ypos, lastState.Xtrans, lastState.Ytrans)
					}

					m.playerStateQueue[playerData.ID].Push(&PastState{byte(m.Frame), nextPosX, nextPosY, 127, 127})
				}

				if m.Frame == currentFrame {
					// the input msgs need to be processed after the frame has been increased to be able to consider input msgs that arrived shortly before
					m.updateMatchState(m.network.Connection)
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
		if v.IPAddr == addr {
			m.recvCountMap[v.ID] = true
		}
	}
}

func (m *Match) processMessages() {
	for v := range m.network.RecvCh {
		pc := v.Connection
		addr := v.Addr
		buf := v.Buffer
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

			if msgID == msg.AbilityInputMsgID {
				m.handleAbilityInputMsg(pc, addr, buf)
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

// AddPlayer adds servers to the match
func (m *Match) AddPlayer(addr net.Addr) int {

	for i := 0; i < len(m.players); i++ {
		playerToCheck := m.players[i].IPAddr.String()
		incomingAddr := addr.String()
		if playerToCheck == incomingAddr {
			return -1
		}
	}

	// player not in match yet & match not full
	var playerID = byte(len(m.players))
	var xPosition = -2800 * int32(playerID)
	// put players in different teams to allow pvp
	// TODO: Remove this later
	var playerData = units.NewPlayer(playerID, 0, playerID, xPosition, 0, addr)
	m.allUnits = append(m.allUnits, playerData)

	m.players = append(m.players, playerData)
	m.playerStateQueue[len(m.players)-1] = *NewQueue(15)
	m.playerStateQueue[len(m.players)-1].Push(&PastState{byte(0), playerData.X, playerData.Y, 127, 127})
	return len(m.players) - 1
}

// CheckMatchFull changes the matchstate when all players joined
func (m *Match) CheckMatchFull(pc net.PacketConn, addr net.Addr) {
	if len(m.players) == m.playerCount {
		for _, v := range m.players {
			m.sendMatchStart(pc, v.IPAddr)
		}
		time.Sleep(time.Second)
		go doEvery(1*time.Millisecond, m.incFrame)
		m.State = 1
		m.StateChangeTimestamp = time.Now().Unix()
		m.recvCountMap = make([]bool, len(m.players))
		m.start = time.Now()
		m.end = time.Now().Add(time.Second * 300)
		log.Println("Server started match")

		m.addInitialNPCUnits()

		time.Sleep(time.Millisecond * 100)
		m.spawnInitialUnits(pc)
	}
}

// TODO: make this more generic
func (m *Match) addInitialNPCUnits() {
	for i := 0; i < 4; i++ {
		var x, y int32
		y = -18000

		switch i {
		case 0:
			x = 3000
			break
		case 1:
			x = -3000
			break
		case 2:
			x = 6000
			break
		case 3:
			x = -6000
			break
		}

		npcUnitData := &units.NPCUnit{
			ID:            byte(i + 3),
			Team:          1,
			UnitType:      3,
			X:             x,
			Y:             y,
			Rotation:      0,
			HealthPercent: 100,
			Collider:      &collision.CircleCollider{Xpos: x, Ypos: y, Radius: 1000}}

		if i == 1 {
			var aiData ai.FireBoss
			aiData.Init(npcUnitData)
			m.aiList = append(m.aiList, aiData)
		}

		m.npcUnits = append(m.npcUnits, npcUnitData)
		m.allUnits = append(m.allUnits, npcUnitData)
	}
}

func (m *Match) serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	response := []byte(" Alive!")
	m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: addr, Buffer: response}
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
	m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: addr, Buffer: rsp}
}

func (m *Match) handleTimeSyncDone(pc net.PacketConn, addr net.Addr, buf []byte,
	playerID int) {
	timesyncdoneackmsg := msg.TimeSyncDoneAckMsg{MessageID: msg.TimeSyncDoneAckMsgID, PlayerID: byte(playerID)}
	m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: addr, Buffer: timesyncdoneackmsg.Encode()}
}

func (m *Match) handlePing(pc net.PacketConn, addr net.Addr, buf []byte) {
	for playerID, playerData := range m.players {
		if playerData.IPAddr.String() == addr.String() {
			m.players[playerID].LastPingTime = time.Now()
			log.Println("Received ping msg from player:", playerID, " at time: ", time.Now().Unix())
		}
	}
	pongMsg := msg.PongMsg{
		MessageID:             msg.PongMsgID,
		TransmissionTimestamp: binary.LittleEndian.Uint64(buf[1:])}
	rsp := pongMsg.Encode()
	m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: addr, Buffer: rsp}
}

func (m *Match) sendMatchStart(pc net.PacketConn, addr net.Addr) {
	matchstart := msg.MatchStartMsg{MessageID: msg.MatchStartMsgID,
		MatchStartTimestamp: uint64(time.Now().UnixNano()/1000000 + 1000)}
	// ts is in ms and match start in now + 1 second
	m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: addr, Buffer: matchstart.Encode()}
}

func (m *Match) handleInputMsg(pc net.PacketConn, addr net.Addr, buf []byte) {
	inputmsg := msg.DecodeInputMsg(buf)
	//log.Println("Pkg Received: ", inputmsg)
	m.pendingInputMsgs = append(m.pendingInputMsgs, inputmsg)
}

func (m *Match) handleAbilityInputMsg(pc net.PacketConn, addr net.Addr, buf []byte) {
	Abilityinputmsg := msg.DecodeAbilityInputMsg(buf)
	//log.Println("Pkg Received: ", inputmsg)
	m.pendingAbilityInputMsgs = append(m.pendingAbilityInputMsgs, Abilityinputmsg)
}

func (m *Match) updateMatchState(pc net.PacketConn) {
	updatedUnitIDs := make(map[byte]bool)
	for _, inputmsg := range m.pendingInputMsgs {
		for _, playerData := range m.players {
			if byte(playerData.ID) == inputmsg.PlayerID {
				updatedUnitIDs = m.updatePlayerPosition(playerData, inputmsg, updatedUnitIDs)
			}
		}
	}

	// clear pending input msgs
	m.pendingInputMsgs = make([]msg.InputMsg, 0, len(m.players))

	for _, abilityinputmsg := range m.pendingAbilityInputMsgs {
		for _, playerData := range m.players {
			if byte(playerData.ID) == abilityinputmsg.PlayerID {
				m.processPendingAbilityInputMsg(pc, abilityinputmsg)
			}
		}
	}

	for _, ai := range m.aiList {
		updatedUnitIDs = ai.Tick(m.allUnits, updatedUnitIDs)
	}

	// clear pending input msgs
	m.pendingAbilityInputMsgs = make([]msg.AbilityInputMsg, 0, len(m.players))

	updatedUnitIDs = m.updateAbilityStates(updatedUnitIDs)

	m.sendMessagesForUpdatedUnits(pc, updatedUnitIDs)
}

func (m *Match) updateAbilityStates(updatedUnitIDs map[byte]bool) map[byte]bool {
	for playerID, ability := range m.abilities {
		if ability == nil {
			continue
		}

		updatedUnitIDs = ability.Tick(m.allUnits, m.Frame, updatedUnitIDs)

		if !ability.IsActive(m.Frame) {
			//log.Println("Removed at: ", m.Frame)
			m.abilities[playerID] = nil
		}
	}

	return updatedUnitIDs
}

func (m *Match) updatePlayerPosition(playerData *units.Player, inputmsg msg.InputMsg, updatedUnitIDs map[byte]bool) map[byte]bool {
	foundFrame := false
	for _, pastState := range m.playerStateQueue[playerData.ID].nodes {
		if pastState != nil {
			// set translation to past frame
			if inputmsg.Frame == pastState.Frame {
				pastState.Xtrans = inputmsg.XTranslation
				pastState.Ytrans = inputmsg.YTranslation
				foundFrame = true
				break
			}
		}
	}

	if !foundFrame {
		frames := make([]byte, 0)
		for _, pastState := range m.playerStateQueue[playerData.ID].nodes {
			frames = append(frames, pastState.Frame)
		}

		log.Println("Did not find frame: ", inputmsg.Frame, " Frames: ", frames)
	}

	// calculate/validate all movements from predecessor
	for index, pastState := range m.playerStateQueue[playerData.ID].nodes {
		if index > 0 {
			queue := m.playerStateQueue[playerData.ID]
			node := queue.nodes[index-1]

			// node hasn't been set yet.
			if node == nil {
				continue
			}

			prevXPos := node.Xpos
			prevYPos := node.Ypos
			prevXTrans := node.Xtrans
			prevYTrans := node.Ytrans

			updXPos, updYPos := units.GetPlayerPosition(prevXPos, prevYPos, prevXTrans, prevYTrans)
			pastState.Xpos, pastState.Ypos = updXPos, updYPos
		}
	}

	latestNode := m.playerStateQueue[playerData.ID].nodes[len(m.playerStateQueue[playerData.ID].nodes)-1]
	// validate move
	m.players[playerData.ID].X, m.players[playerData.ID].Y = units.GetPlayerPosition(
		latestNode.Xpos, latestNode.Ypos, latestNode.Xtrans, latestNode.Ytrans)
	m.players[playerData.ID].Rotation = inputmsg.Rotation

	m.players[playerData.ID].Collider.Update(m.players[playerData.ID].X, m.players[playerData.ID].Y, m.players[playerData.ID].Rotation)

	updatedUnitIDs[playerData.ID] = true

	return updatedUnitIDs
}

func (m *Match) sendMessagesForUpdatedUnits(pc net.PacketConn, updatedUnitIDs map[byte]bool) {
	for unitID, isUpdated := range updatedUnitIDs {
		if !isUpdated {
			continue
		}

		var unitData units.Unit
		unitFound := false

		for _, unitDataToCheck := range m.allUnits {
			if unitDataToCheck.GetID() == unitID {
				unitData = unitDataToCheck
				unitFound = true
				break
			}
		}

		if !unitFound {
			log.Println("Did not find unit with ID: ", unitID)
			continue
		}

		x, y := unitData.GetPosition()

		// players state for all other clients
		unitstatemsg := msg.UnitStateMsg{
			MessageID:     msg.UnitStateMsgID,
			UnitID:        unitData.GetID(),
			XPosition:     x,
			YPosition:     y,
			Rotation:      unitData.GetRotation(),
			HealthPercent: byte(unitData.GetHealthPercent()),
			Frame:         m.Frame}

		// unitstate for all clients
		for _, v := range m.players {
			// optimize this: split position/rotation mesage from general state messages?
			m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: v.IPAddr, Buffer: unitstatemsg.Encode()}
		}

		if unitData.IsPlayer() && m.playerStateQueue[unitData.GetID()].nodes[0] != nil {
			oldState := m.playerStateQueue[unitData.GetID()].nodes[0]
			oldXPos, oldYPos := units.GetPlayerPosition(oldState.Xpos, oldState.Ypos, oldState.Xtrans, oldState.Ytrans)

			resp := msg.PositionConfirmationMsg{
				MessageID: msg.PositionConfirmationMsgID,
				UnitID:    unitData.GetID(),
				XPosition: oldXPos,
				YPosition: oldYPos,
				Frame:     oldState.Frame}

			//log.Println("Sending pcm with frame: ", resp.Frame)

			for _, v := range m.players {
				if v.ID == unitData.GetID() {
					m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: v.IPAddr, Buffer: resp.Encode()}
				}
			}
		}
	}
}

func (m *Match) processPendingAbilityInputMsg(pc net.PacketConn, inputmsg msg.AbilityInputMsg) {
	playerData := &units.EmptyPlayer
	for _, v := range m.players {
		if v.ID == inputmsg.PlayerID {
			playerData = v
			break
		}
	}

	if playerData == &units.EmptyPlayer {
		return
	}

	ability, prs := m.abilities[inputmsg.PlayerID]

	if prs && ability != nil {
		// ability already active for player, a second one is not allowed
		return
	}

	// make this better
	// TODO: Check if player can use the abiity id

	switch inputmsg.AbilityID {
	case 0:
		m.abilities[inputmsg.PlayerID] = &abilities.WarriorMeeleAbility{}
	}

	var abilityData = abilities.AbilityData{
		AbilityID:    inputmsg.AbilityID,
		CasterUnitID: inputmsg.PlayerID,
		Rotation:     inputmsg.Rotation,
		StartFrame:   inputmsg.StartFrame}

	abilityData = m.abilities[inputmsg.PlayerID].Init(abilityData, playerData)

	m.players[playerData.ID].Rotation = inputmsg.Rotation

	abilityActMsg := msg.UnitAbilityActivationMsg{
		MessageID:       msg.UnitAbilityActivationMsgID,
		AbilityID:       inputmsg.AbilityID,
		UnitID:          inputmsg.PlayerID,
		Rotation:        inputmsg.Rotation,
		StartFrame:      inputmsg.StartFrame,
		ActivationFrame: abilityData.ActivationFrame,
		EndFrame:        abilityData.EndFrame}

	for _, v := range m.players {
		// don't send to sender of input msg
		if v.ID != inputmsg.PlayerID {
			m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: v.IPAddr, Buffer: abilityActMsg.Encode()}
		}
	}
}

func (m *Match) spawnInitialUnits(pc net.PacketConn) {
	for _, unitData := range m.allUnits {

		x, y := unitData.GetPosition()
		// players state for all other clients
		unitSpawnMsg := msg.UnitSpawnMsg{
			MessageID:     msg.UnitSpawnMsgID,
			UnitID:        unitData.GetID(),
			TeamID:        unitData.GetTeam(),
			UnitType:      unitData.GetUnitType(),
			XPosition:     x,
			YPosition:     y,
			Rotation:      unitData.GetRotation(),
			HealthPercent: byte(unitData.GetHealthPercent()),
			Frame:         0}

		// unitstate for all clients
		for _, v := range m.players {
			m.network.SendCh <- &network.OutPkt{Connection: pc, Addr: v.IPAddr, Buffer: unitSpawnMsg.Encode()}
		}
	}
}
