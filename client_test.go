package main

import (
	"encoding/gob"
	"net"
	"testing"
	"fsync/protocol"
	"fsync/util"
)

const addr = "127.0.0.1:2000"

// Receive and print single packet
func Test1_ReceivePacket(t *testing.T) {
	var s protocol.SocketHandler
	var err error
	s.Con, err = net.Dial("tcp", addr)
	util.CheckError(err)
	gob.Register(protocol.Packet{})
	dec := gob.NewDecoder(s.Con)
	p := protocol.Packet{}
	err = dec.Decode(&p)
	util.CheckError(err)
	t.Logf("\nOrder num: %d\nBody size: %d\nBody: %s", p.OrderNum, p.BodySize, p.Body)
}

// Receive packets and print their order num
func Test2_ReceivePacketizedFile(t *testing.T) {
	var s protocol.SocketHandler
	var err error
	s.Con, err = net.Dial("tcp", addr)
	util.CheckError(err)
	
	packets, err := s.ReceivePackets()
	util.CheckError(err)

	for _, packet := range packets {
		t.Logf("\nPacket #: %d", packet.OrderNum)
	}
}
