package main

import (
	"encoding/gob"
	dir "fsync/directory"
	"fsync/protocol"
	"fsync/util"
	"net"
	"os"
	"testing"
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
	t.Logf("\nOrder num: %d\nBody size: %d\n", p.OrderNum, p.Body)
}

// Receive packets and print their order num
func Test2_PrintPacketData(t *testing.T) {
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

func Test3_WritePacketsToFile(t *testing.T) {
	var s protocol.SocketHandler
	var err error

	s.Con, err = net.Dial("tcp", addr)
	util.CheckError(err)

	packets, err := s.ReceivePackets()
	util.CheckError(err)

	var d dir.DirManager
	d.Path, err = os.Getwd()
	util.CheckError(err)

	data := protocol.GetPacketData(packets)
	util.CheckError(err)

	name := "3d-geo.jpg"
	d.WriteDataToFile(name, data)
}
