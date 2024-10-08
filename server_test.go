package main

import (
	"encoding/gob"
	"net"
	"os"
	"testing"

	clt "github.com/sebastian-j-ibanez/fsync/client"
	dir "github.com/sebastian-j-ibanez/fsync/directory"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
)

const testImgPath = "./img/3d-geo.jpg"

// Send single packet
func Test1_SendPktNum(t *testing.T) {
	var s prot.SocketHandler
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	s.Conn, err = lis.Accept()
	if err != nil {
		t.Fatal(err)
	}

	pktNum := int64(1)
	enc := gob.NewEncoder(s.Conn)
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

	s.Conn, err = lis.Accept()
	if err != nil {
		t.Fatal(err)
	}

	err = s.UploadFile(testImgPath)
	if err != nil {
		t.Fatal(err)
	}
}

func Test3_AwaitSync(t *testing.T) {
	const port = 2000

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
	err = c.AwaitSync(port)
	if err != nil {
		t.Fatal(err)
	}
}
