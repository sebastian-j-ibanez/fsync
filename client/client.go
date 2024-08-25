package client

import (
	"errors"
	"net"

	dir "github.com/sebastian-j-ibanez/fsync/directory"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
)

type Client struct {
	DirMan dir.DirManager
	Sock   prot.SocketHandler
	Peers  []prot.Peer
}

// Init sync with peers
func (c Client) InitSync() error {
	var err error
	cltHashes, err := c.DirMan.HashDir()
	if err != nil {
		msg := "unable to hash directory: " + err.Error()
		return errors.New(msg)
	}

	for _, peer := range c.Peers {
		// Connect to peer
		c.Sock.Con, err = net.Dial("tcp", peer.Addr())
		if err != nil {
			msg := "unable to connect to peer: " + err.Error()
			return errors.New(msg)
		}

		// Get peer file hashes
		peerHashes, err := c.Sock.ReceiveFileHashes()
		if err != nil {
			msg := "unable to receive file hashes: " + err.Error()
			return errors.New(msg)
		}

		// Upload files that peer does not have
		uniqueFiles := dir.GetUniqueHashes(cltHashes, peerHashes)

		// Send unique file hashes to server
		err = c.Sock.SendFileHashes(*uniqueFiles)
		if err != nil {
			msg := "unable to send file hashes: " + err.Error()
			return errors.New(msg)
		}

		for _, file := range *uniqueFiles {
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

// Init sync with peers
func (c Client) InitSyncWithPeer(peer prot.Peer) error {
	var err error
	cltHashes, err := c.DirMan.HashDir()
	if err != nil {
		msg := "unable to hash directory: " + err.Error()
		return errors.New(msg)
	}

	// Connect to peer
	c.Sock.Con, err = net.Dial("tcp", peer.Addr())
	if err != nil {
		msg := "unable to connect to peer: " + err.Error()
		return errors.New(msg)
	}

	// Get peer file hashes
	peerHashes, err := c.Sock.ReceiveFileHashes()
	if err != nil {
		msg := "unable to receive file hashes: " + err.Error()
		return errors.New(msg)
	}

	// Upload files that peer does not have
	uniqueFiles := dir.GetUniqueHashes(cltHashes, peerHashes)

	// Send unique file hashes to server
	err = c.Sock.SendFileHashes(*uniqueFiles)
	if err != nil {
		msg := "unable to send file hashes: " + err.Error()
		return errors.New(msg)
	}

	for _, file := range *uniqueFiles {
		path := c.DirMan.Path + "/" + file.Name
		err = c.Sock.UploadFile(path)
		if err != nil {
			msg := "unable to upload " + file.Name + ": " + err.Error()
			return errors.New(msg)
		}
	}

	return nil
}

// Await sync from peer over default port
func (c Client) AwaitSync() error {
	var err error
	lis, err := net.Listen("tcp", ":2000")
	if err != nil {
		return err
	}

	// Accept peer connection
	c.Sock.Con, err = lis.Accept()
	if err != nil {
		return err
	}

	// Get and send file hashes
	hashes, err := c.DirMan.HashDir()
	if err != nil {
		msg := "unable to hash directory: " + err.Error()
		return errors.New(msg)
	}

	err = c.Sock.SendFileHashes(hashes)
	if err != nil {
		msg := "unable to send file hashes: " + err.Error()
		return errors.New(msg)
	}

	// Receive unique file hashes
	uniqueHashes, err := c.Sock.ReceiveFileHashes()
	if err != nil {
		msg := "unable to receive file hashes: " + err.Error()
		return errors.New(msg)
	}

	for _, file := range uniqueHashes {
		err := c.Sock.DownloadFile(file.Name)
		if err != nil {
			msg := "unable to download " + file.Name + ": " + err.Error()
			return errors.New(msg)
		}
	}

	return nil
}
