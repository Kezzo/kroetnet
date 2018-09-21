package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device            = "lo0"
	snapshotLen int32 = 1024
	promiscuous       = false
	err         error
	timeout     = 30 * time.Second
	handle      *pcap.Handle
)

func main() {

	port := ":2448"
	network := "udp"

	pc, err := net.ListenPacket(network, port)
	handleError(err)

	fmt.Printf("listening on (%s)%s\n", network, pc.LocalAddr())
	defer pc.Close()

	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	handleError(err)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	defer handle.Close()

	for packet := range packetSource.Packets() {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		fmt.Printf("\nBuffer Content: [ % x ] \n", buf[:n])

		if err != nil {
			log.Print("Error: ", err)
			continue
		}

		go sniffPackets(pc, addr, packet)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatalln("Error: ", err)
	}
}

func handleRequest(pc net.PacketConn, addr net.Addr, buf []byte, udp *layers.UDP) {
	log.Printf("received string: %s from: %s\n\n", string(buf), addr)
	log.Println("received buffer at 0: ", buf[0])
	switch buf[0] {
	case 35:
		handleNTPReq(pc, addr, buf, udp)
		break
	default:
		serve(pc, addr, buf)
	}

}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	response := []byte("Received your msg: " + string(buf))
	pc.WriteTo(response, addr)
}

func sniffPackets(pc net.PacketConn, addr net.Addr, packet gopacket.Packet) {
	//some details about the packet

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		log.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)
		log.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		log.Println("Protocol: ", ip.Protocol)
		log.Println()
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		log.Println("UDPLayer detected.")
		buf := udpLayer.LayerPayload()
		log.Printf("%s\n")
		udp, _ := udpLayer.(*layers.UDP)
		log.Println("Content: ", udp)

		handleRequest(pc, addr, buf, udp)
	}

	ntpLayer := packet.Layer(layers.LayerTypeNTP)
	if ntpLayer != nil {
		log.Println("NTP Layer detected.")
		log.Printf("%s\n", ntpLayer.LayerPayload())
		ntp, _ := ntpLayer.(*layers.NTP)
		log.Println("Content: ", ntp)
	}

	log.Println("Found Layers:")
	for _, layer := range packet.Layers() {
		log.Println("- ", layer.LayerType())
	}

}
