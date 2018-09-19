package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := ":2448"
	fmt.Println("Start listening on port " + port)

	pc, err := net.ListenPacket("udp", port)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)

		if err != nil {
			log.Print("Error: ", err)
			continue
		}

		log.Println("Received: ", buf)

		go serve(pc, addr, buf[:n])
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	buf[2] |= 0x80
	pc.WriteTo(buf, addr)
}
