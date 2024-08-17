package main

import (
	prot "fsync/protocol"
	"net"
	"os"
	"testing"
)

const addr = "127.0.0.1:2000"

// Receive packets and print their order num
func Test1_ReceivePktNum(t *testing.T) {
	var s prot.SocketHandler
	var err error
	s.Con, err = net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	
	pktNum, err := s.ReceivePktNum()
	if err != nil {
		t.Fatal(err)
	}

	expected := int64(1)
	if pktNum != expected {
		t.Fatalf("expected: %d\treceived: %d", expected, pktNum)
	}
}

func Test2_DownloadFile(t *testing.T) {
	var s prot.SocketHandler
	var err error
	s.Con, err = net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	destination, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	
	err = s.DownloadFile(destination + "/3d-geo.jpg")
	if err != nil {
		t.Fatal(err)
	}
}


