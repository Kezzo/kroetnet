package main

import (
	"log"
	"net"
	"time"

	"github.com/google/gopacket"

	"github.com/google/gopacket/layers"
)

func handleNTPReq(pc net.PacketConn, addr net.Addr, req []byte, udp *layers.UDP) {

	log.Println("received ntp packet ", req)

	ntp := layers.NTP{}

	var udpserbuf gopacket.SerializeBuffer = gopacket.NewSerializeBuffer()
	serset := gopacket.SerializeOptions{}
	udp.SerializeTo(udpserbuf, serset)
	ntp.DecodeFromBytes(udpserbuf.Bytes(), nil)

	var ts layers.NTPTimestamp = layers.NTPTimestamp(uint64(time.Now().Unix()))
	ntp.ReceiveTimestamp = ts
	ntp.TransmitTimestamp = ts

	log.Println("created ntp packet ", ntp)

	var serbuf gopacket.SerializeBuffer = gopacket.NewSerializeBuffer()
	ntp.SerializeTo(serbuf, serset)

	resp := serbuf.Bytes()

	log.Println("sending response ntp packet ", resp)

	// send data
	if _, err := pc.WriteTo(resp, addr); err != nil {
		log.Fatalln("err sending data:", err)
	}
}

// getNTPSecs decompose current time as NTP seconds
func getNTPSeconds(t time.Time) (int64, int64) {
	// convert time to total # of secs since 1970
	// add NTP epoch offets as total #secs between 1900-1970
	secs := t.Unix() + int64(getNTPOffset())
	fracs := t.Nanosecond()
	return secs, int64(fracs)
}

// getNTPOffset returns the 70yrs between Unix epoch
// and NTP epoch (1970-1900) in seconds
func getNTPOffset() float64 {
	ntpEpoch := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	offset := unixEpoch.Sub(ntpEpoch).Seconds()
	return offset
}
