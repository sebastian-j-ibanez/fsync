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
	enc := gob.NewEncoder(s.Con)
	err = enc.Encode(pktNum)
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

func Test3_AwaitSync(t *testing.T) {
	// Init directory manager
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	d, err := dir.NewDirManager(path + "/img-test")
	if err != nil {
		t.Fatal(err)
	}

	// Init client
	c := clt.Client{
		DirMan: *d,
	}

	// Await sync request
	err = c.AwaitSync()
	if err != nil {
		t.Fatal(err)
	}
}
