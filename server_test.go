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

const testImgPath = "/test-img/3d-geo.jpg"

// Send single packet
func Test1_SendBasicPacket(t *testing.T) {
	b := []byte{1,2,3}
	var p prot.Packet
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
func Test2_SendReadme(t *testing.T) {	
	var d dir.DirManager
	var err error
	d.Path, err = os.Getwd()
	util.CheckError(err)

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

func Test3_SendImageFile(t *testing.T) {
	home, err := os.UserHomeDir()
	util.CheckError(err)
	d := dir.DirManager {
		Path: home,
	}
	
	packets, err := d.PacketizeFile(testImgPath)
	util.CheckError(err)
	
	var s prot.SocketHandler
	lis, err := net.Listen("tcp", addr)
	util.CheckError(err)

	s.Con, err = lis.Accept()
	util.CheckError(err)

	s.SendPackets(packets)
}

func Test4_WriteFile(t *testing.T) {
	path, err := os.UserHomeDir()
	fullPath := path + "/test.txt"
	if err = os.WriteFile(fullPath, []byte("test"), 0666); err != nil {
		t.Error(err)
	}
}
