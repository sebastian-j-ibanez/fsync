package client

import (
	"errors"
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

// Init sync with peers
func (c Client) InitSync() error {
	for _, peer := range c.Peers {
		// Connect to peer, init socket
		conn, err := net.Dial("tcp", peer.Addr())
		if err != nil {
			msg := "unable to establish connection: " + err.Error()
			return errors.New(msg)
		}
		c.Sock = prot.NewSocketHandler(conn)

		// Get peer file hashes
		var peerHashes []dir.FileHash
		err = c.Sock.ReceiveGenericData(&peerHashes)
		if err != nil {
			msg := "unable to receive file hashes: " + err.Error()
			return errors.New(msg)
		}

		// Get local file hashes
		localHashes, err := c.DirMan.HashDir()
		if err != nil {
			msg := "unable to hash directory: " + err.Error()
			return errors.New(msg)
		}

		// Send unique file hashes
		uniqueFiles := dir.GetUniqueHashes(localHashes, peerHashes)
		err = c.Sock.SendGenericData(*uniqueFiles)
		if err != nil {
			msg := "unable to send file hashes: " + err.Error()
			return errors.New(msg)
		}

		for _, file := range *uniqueFiles {
			// Send file packets
			path := c.DirMan.Path + "/" + file.Name
			err = c.Sock.UploadFile(path)
			if err != nil {
				msg := "unable to upload " + file.Name + ": " + err.Error()
				return errors.New(msg)
			}
		}
	}

	return nil
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

	// Accept peer connection
	conn, err := lis.Accept()
	if err != nil {
		return err
	}

	// Init socket handler
	c.Sock = prot.NewSocketHandler(conn)
	if c.Sock.Conn == nil {
		return errors.New("unable to establish connection")
	}

	// Get and send file hashes
	hashes, err := c.DirMan.HashDir()
	if err != nil {
		msg := "unable to hash directory: " + err.Error()
		return errors.New(msg)
	}
	err = c.Sock.SendGenericData(hashes)
	if err != nil {
		msg := "unable to send file hashes: " + err.Error()
		return errors.New(msg)
	}

	// Receive unique file hashes
	var uniqueHashes []dir.FileHash
	err = c.Sock.ReceiveGenericData(&uniqueHashes)
	if err != nil {
		msg := "unable to receive file hashes: " + err.Error()
		return errors.New(msg)
	}

	for _, file := range uniqueHashes {
		// Receive file packets
		err = c.Sock.DownloadFile(file.Name)
		if err != nil {
			msg := "unable to download " + file.Name + ": " + err.Error()
			return errors.New(msg)
		}
	}

	return nil
}
