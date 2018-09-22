package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

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
		go digestRequest(pc, addr, buf[:n])
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatalln("Error: ", err)
	}
}

func digestRequest(pc net.PacketConn, addr net.Addr, buf []byte) {
	recvTime := time.Now()
	// log.Printf("received string: %s from: %s\n\n", string(buf), addr)
	log.Println("received buffer", buf)
	msgID := buf[0]
	switch msgID {
	case 1:
		handleNTPReq(pc, addr, buf, recvTime)
	case 2:
		handleKroetPkg(pc, addr, buf)
	default:
		serve(pc, addr, buf)
	}

}
