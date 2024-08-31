package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"testing"

	clt "github.com/sebastian-j-ibanez/fsync/client"
	dir "github.com/sebastian-j-ibanez/fsync/directory"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
)

const addr = "127.0.0.1:2000"

// Receive packets and print their order num
func Test1_ReceivePktNum(t *testing.T) {
	var s prot.SocketHandler
	var err error
	s.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	dec := gob.NewDecoder(s.Conn)
	var pktNum int64
	err = dec.Decode(&pktNum)
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
	s.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	destination, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = s.DownloadFile(destination + "/img-test/3d-geo.jpg")
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
	if err != nil {
		t.Fatal(err)
	}

	// Mock peer
	peer := prot.Peer{
		IP:   "127.0.0.1",
		Port: "2000",
	}

	// Init client
	c := clt.Client{
		DirMan: *d,
		Peers:  []prot.Peer{peer},
	}

	// Init sync with peer
	err = c.InitSync()
	if err != nil {
		t.Fatal(err)
	}
}

func Test4_RegisterPeer(t *testing.T) {
	// Init file
	_, err := os.Create("peer_data.json")
	if err != nil {
		t.Fatal(err)
	}

	// Register mock peer
	p := prot.Peer{
		IP:   "127.0.0.1",
		Port: "2000",
	}
	prot.RegisterPeer(p)
	fmt.Println("registered default peer.")
}
