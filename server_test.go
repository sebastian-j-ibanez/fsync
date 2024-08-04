package main

import (
	"fmt"
	"net"
	"testing"	
	"fsyncprotocol/protocol"
	"fsyncprotocol/util"
)

// Send single packet
func Test1_SendBasicPacket(t *testing.T) {
	b := []byte{1,2,3}
	var p protocol.Packet
	p.OrderNum = 1
	p.BodySize = 3
	p.Body = b
	fmt.Println(p)
	
	var s protocol.SocketHandler
	lis, err := net.Listen("tcp", addr)
	util.CheckError(err)

	s.Con, err = lis.Accept()
	util.CheckError(err)

	err = s.SendPacket(p)
	util.CheckError(err)
}

// Packetize and send README.md
func Test2_SendPacketizedFile(t *testing.T) {
	path := "/home/sebas/Code/fsync-protocol"
	name := "README.md"
	packets, err := protocol.PacketizeFile(path, name)
	util.CheckError(err)

	var s protocol.SocketHandler
	lis, err := net.Listen("tcp", addr)
	util.CheckError(err)

	s.Con, err = lis.Accept()
	util.CheckError(err)

	s.SendPackets(packets)
}
