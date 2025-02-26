package client

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	dir "github.com/sebastian-j-ibanez/fsync/directory"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
)

const (
	defaultPort = 2000
	serviceName = "_fsync._tcp"
	timeout     = 10 * time.Second
)

type Client struct {
	DirMan dir.DirManager
	Sock   prot.SocketHandler
	Peers  []prot.Peer
}

// Await sync from peer over default port
func (c Client) AwaitSync(portNum int) error {
	// Set port
	if portNum == -1 {
		portNum = defaultPort
	}
	port := ": " + strconv.Itoa(portNum)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	fmt.Printf("Listening over port %d...\n", portNum)

	// Accept peer connection
	conn, err := lis.Accept()
	if err != nil {
		return err
	}
	fmt.Printf("Connection established with client (%s)\n", conn.RemoteAddr().String())

	c.Sock, err = prot.NewSocketHandler(conn, true)
	if err != nil {
		return errors.New("unable to establish connection: " + err.Error())
	}

	// Send local hashes
	err = c.SendUniqueHashes(nil)
	if err != nil {
		msg := "unable to send file hashes: " + err.Error()
		return errors.New(msg)
	}

	// Receive file hashes
	var uniqueHashes []dir.FileHash
	uniqueHashes, err = c.ReceiveUniqueHashes()
	if err != nil {
		msg := "unable to receive file hashes: " + err.Error()
		return errors.New(msg)
	}

	err = c.ReceiveUniqueFiles(uniqueHashes)
	if err != nil {
		return err
	}

	return nil
}

// Init sync with peers
func (c Client) InitSync(filePattern []string) error {
	// Get local file hashes
	localHashes, err := c.DirMan.GetFileHashes(filePattern)
	if err != nil {
		msg := "unable to hash directory: " + err.Error()
		return errors.New(msg)
	}

	for _, peer := range c.Peers {
		// Connect to peer, init socket
		conn, err := net.Dial("tcp", peer.Addr())
		if err != nil {
			msg := "unable to establish connection: " + err.Error()
			return errors.New(msg)
		}

		c.Sock, err = prot.NewSocketHandler(conn, false)
		if err != nil {
			return errors.New("unable to initialize socket handler: " + err.Error())
		}

		// Get peer file hashes
		peerHashes, err := c.ReceiveUniqueHashes()
		if err != nil {
			msg := "unable to receive file hashes: " + err.Error()
			return errors.New(msg)
		}

		// Send unique file hashes
		uniqueFiles := dir.GetUniqueHashes(localHashes, peerHashes)
		err = c.SendUniqueHashes(*uniqueFiles)
		if err != nil {
			msg := "unable to send file hashes: " + err.Error()
			return errors.New(msg)
		}

		err = c.SendUniqueFiles(*uniqueFiles)
		if err != nil {
			return err
		}
	}

	return nil
}

// Receive file hashes from socket
func (c Client) ReceiveUniqueHashes() ([]dir.FileHash, error) {
	uniqueHashes := []dir.FileHash{}

	err := c.Sock.ReceiveEncryptedData(&uniqueHashes, prot.FileHashes)
	if err != nil {
		return nil, err
	}

	return uniqueHashes, nil
}

// Sends unique hashes over client socket
// Default to file hashes of every file in the directory
func (c Client) SendUniqueHashes(uniqueHashes []dir.FileHash) error {
	var err error
	if uniqueHashes == nil {
		uniqueHashes, err = c.DirMan.GetFileHashes(nil) // Empty slice will default to all files in directory
		if err != nil {
			return err
		}
	}

	var pkt prot.Packet
	err = pkt.SerializeToBody(uniqueHashes, prot.FileHashes)
	if err != nil {
		return err
	}

	err = c.Sock.SendEncryptedPacket(pkt)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) SendUniqueFiles(uniqueFiles []dir.FileHash) error {
	var err error
	for _, file := range uniqueFiles {
		path := c.DirMan.Path + "/" + file.Name
		err = c.Sock.UploadFile(path)
		if err != nil {
			msg := "unable to upload " + file.Name + ": " + err.Error()
			return errors.New(msg)
		}
	}

	return nil
}

func (c Client) ReceiveUniqueFiles(uniqueHashes []dir.FileHash) error {
	var err error
	for _, file := range uniqueHashes {
		err = c.Sock.DownloadFile(file.Name)
		if err != nil {
			msg := "unable to download " + file.Name + ": " + err.Error()
			return errors.New(msg)
		}
	}

	return nil
}
