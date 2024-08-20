package main

import (
	"encoding/gob"
	clt "fsync/client"
	dir "fsync/directory"
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

	
	dec := gob.NewDecoder(s.Con)
	var pktNum int64
	err = dec.Decode(pktNum)
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

func Test3_InitSync(t *testing.T) {
	// Init directory manager
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	d, err := dir.NewDirManager(path + "/img")

	// Mock peer
	peer := prot.Peer{
		IP: "127.0.0.1",
		Port: "2000",
	}

	// Init client
	c := clt.Client{
		DirMan: *d,
		Peers: []prot.Peer { peer },
	}

	// Init sync with peer
	err = c.InitSync()
	if err != nil {
		t.Fatal(err)
	}
}

