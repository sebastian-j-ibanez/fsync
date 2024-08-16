package main

import (
	"fmt"
	dir "fsync/directory"
	prot "fsync/protocol"
	"fsync/util"
	"net"
	"os"
	"testing"
)

// Send single packet
func Test1_SendBasicPacket(t *testing.T) {
	b := []byte{1,2,3}
	var p prot.Packet
	p.OrderNum = 1
	p.Body = b
	fmt.Println(p)
	
	var s prot.SocketHandler
	lis, err := net.Listen("tcp", addr)
	util.CheckError(err)

	s.Con, err = lis.Accept()
	util.CheckError(err)

	err = s.SendPacket(p)
	util.CheckError(err)
}

// Packetize and send README.md
func Test2_SendPacketizedFile(t *testing.T) {
	home, err := os.UserHomeDir()
	util.CheckError(err)
	
	var d dir.DirManager
	d.Path = home + "/Code/fsync-protocol"
	name := "README.md"
	packets, err := d.PacketizeFile(name)
	util.CheckError(err)

	var s prot.SocketHandler
	lis, err := net.Listen("tcp", addr)
	util.CheckError(err)

	s.Con, err = lis.Accept()
	util.CheckError(err)

	s.SendPackets(packets)
}

func Test3_SendFile(t *testing.T) {
	var d dir.DirManager
	home, err := os.UserHomeDir()
	util.CheckError(err)
	d.Path = home
	name := "3d-geo.jpg"
	packets, err := d.PacketizeFile(name)
	t.Log(len(packets))
	if err != nil {
		t.Error(err)
	}

	data := prot.GetPacketData(packets)
	t.Log(len(data))
	err = d.WriteDataToFile("itworked.jpg", data)
	if err != nil {
		t.Error(err)
	}
	
	// var s prot.SocketHandler
	// lis, err := net.Listen("tcp", addr)
	// util.CheckError(err)

	// s.Con, err = lis.Accept()
	// util.CheckError(err)

	// s.SendPackets(packets)
}

func Test4_WriteFile(t *testing.T) {
	// Open file
	path, err := os.UserHomeDir()
	fullPath := path + "/test.txt"
	if err = os.WriteFile(fullPath, []byte("test"), 0666); err != nil {
		t.Error(err)
	}
}

func Test5_SendImage(t *testing.T) {
	// Read file into packets
	var d dir.DirManager
	d.Path = "/mnt/crucial/Pictures/Wallpapers"
	packets, err := d.PacketizeFile("building-sky.png")

	// Listen for connection
	var s prot.SocketHandler
	lis, err := net.Listen("tcp", "192.168.2.27:2000")
	util.CheckError(err)
	s.Con, err = lis.Accept()
	util.CheckError(err)
	err = s.SendPackets(packets)
	util.CheckError(err)
}
