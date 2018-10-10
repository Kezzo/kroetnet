package main

import (
	"log"
	"net"
)

// RcvPkt is received by ListenUDP
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

// Network ...
type Network struct {
	recvCh    chan *RcvPkt
	sendCh    chan *OutPkt
	Port      string
	connecton *net.UDPConn
}

func newNetwork(port string) *Network {
	return &Network{
		recvCh: make(chan *RcvPkt),
		sendCh: make(chan *OutPkt),
		Port:   port,
	}
}

func (n *Network) sendByteResponse() {
	for v := range n.sendCh {
		/* if v.buffer[0] != 1 {
			log.Println("Reponse to send: ", v.buffer, " to ", v.addr)
		}*/

		if _, err := v.connection.WriteTo(v.buffer, v.addr); err != nil {
			log.Fatalln("err sending data for msgID :", v.buffer[0], err)
		}
	}
}

func (n *Network) listenUDP() {
	udpAddr, err := net.ResolveUDPAddr("udp", n.Port)
	if err != nil {
		panic(err)
	}
	network := "udp"
	pc, err := net.ListenUDP(network, udpAddr)
	// temp solution
	n.connecton = pc
	if err != nil {
		panic(err)
	}
	log.Printf("listening on (%s)%s\n", network, pc.LocalAddr())
	defer pc.Close()
	for {
		buf := make([]byte, 1024)
		num, addr, err := pc.ReadFrom(buf)
		// log.Println("Received buffer ", buf[:num], " from ", addr)
		if err != nil {
			log.Print("Error: ", err)
			continue
		}
		n.recvCh <- &RcvPkt{pc, addr, buf[:num]}
	}

}
