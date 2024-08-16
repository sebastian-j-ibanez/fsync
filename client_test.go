package main

import (
	"encoding/gob"
	dir "fsync/directory"
	prot "fsync/protocol"
	"fsync/util"
	"net"
	"os"
	"testing"
)

const addr = "127.0.0.1:2000"

// Receive packets and print their order num
func Test1_ReceiveBasicPacket(t *testing.T) {
	var s prot.SocketHandler
	var err error
	s.Con, err = net.Dial("tcp", addr)
	util.CheckError(err)
	
	packet, err := s.ReceivePacket()
	util.CheckError(err)
	t.Logf("\nPacket #: %d", len(packet.Body))
}

func Test2_WritePacketsToFile(t *testing.T) {
	var s prot.SocketHandler
	var err error

	s.Con, err = net.Dial("tcp", addr)
	util.CheckError(err)

	packet, err := s.ReceivePackets()
	util.CheckError(err)

	var d dir.DirManager
	d.Path, err = os.Getwd()
	util.CheckError(err)

	data := prot.GetPacketData(packet)
	util.CheckError(err)

	name := "3d-geo.jpg"
	d.WriteDataToFile(name, data)
}

func Test4_ReceiveImgPkt(t *testing.T) {
	var s prot.SocketHandler
	var err error
	s.Con, err = net.Dial("tcp", addr)
	util.CheckError(err)

	var imgPkt prot.Packet
	dec := gob.NewDecoder(s.Con)
	err = dec.Decode(&imgPkt)

	home, err := os.UserHomeDir()
	util.CheckError(err)

	file, err := os.Create(home + "/sunflower.png")
	util.CheckError(err)
	defer file.Close()

	_, err = file.Write(imgPkt.Body)
	util.CheckError(err)
}
