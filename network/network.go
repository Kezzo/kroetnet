package network

import (
	"log"
	"net"
)

// RcvPkt is received by ListenUDP
type RcvPkt struct {
	Connection *net.UDPConn
	Addr       net.Addr
	Buffer     []byte
}

// OutPkt ...
type OutPkt struct {
	Connection net.PacketConn
	Addr       net.Addr
	Buffer     []byte
}

// Network ...
type Network struct {
	RecvCh     chan *RcvPkt
	SendCh     chan *OutPkt
	Port       string
	Connection *net.UDPConn
}

// NewNetwork ...
func NewNetwork(port string) *Network {
	return &Network{
		RecvCh: make(chan *RcvPkt),
		SendCh: make(chan *OutPkt),
		Port:   port,
	}
}

// SendByteResponse ...
func (n *Network) SendByteResponse() {
	for v := range n.SendCh {
		/*if v.Buffer[0] != 1 {
			log.Println("Sending buffer ", v.Buffer, " to ", v.Addr)
		}*/
		if _, err := v.Connection.WriteTo(v.Buffer, v.Addr); err != nil {
			log.Fatalln("err sending data for msgID :", v.Buffer[0], err)
		}
	}
}

// ListenUDP ...
func (n *Network) ListenUDP() {
	udpAddr, err := net.ResolveUDPAddr("udp", n.Port)
	if err != nil {
		panic(err)
	}
	network := "udp"
	pc, err := net.ListenUDP(network, udpAddr)
	// temp solution
	n.Connection = pc
	if err != nil {
		panic(err)
	}
	log.Printf("listening on (%s)%s\n", network, pc.LocalAddr())

	n.Port = pc.LocalAddr().String()
	defer pc.Close()
	buf := make([]byte, 1024)
	for {

		num, addr, err := pc.ReadFrom(buf)

		if buf[0] != 0 {
			log.Println("Received buffer ", buf[:num], " from ", addr)
		}

		if err != nil {
			log.Print("Error: ", err)
			continue
		}

		if num > 0 {
			n.RecvCh <- &RcvPkt{pc, addr, buf[:num]}
			buf = make([]byte, 1024)
		}
	}
}
