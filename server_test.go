package main

import (
	prot "fsync/protocol"
	"net"
	"testing"
)

const testImgPath = "./img/3d-geo.jpg"

// Send single packet
func Test1_SendPktNum(t *testing.T) {
	var s prot.SocketHandler
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	s.Con, err = lis.Accept()
	if err != nil {
		t.Fatal(err)
	}

	pktNum := int64(1)
	err = s.SendPktNum(pktNum)
	if err != nil {
		t.Fatal(err)
	}
}

func Test2_UploadFile(t *testing.T) {
	var s prot.SocketHandler
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	s.Con, err = lis.Accept()
	if err != nil {
		t.Fatal(err)
	}

	err = s.UploadFile(testImgPath)
	if err != nil {
		t.Fatal(err)
	}
}
