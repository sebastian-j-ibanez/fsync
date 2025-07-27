package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	dir "github.com/sebastian-j-ibanez/fsync/directory"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
)

const (
	defaultPort = 8080
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
	port := "0.0.0.0:" + strconv.Itoa(portNum)

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

	// Confirmation prompt
	conf, err := c.confirmDownload(uniqueHashes)
	if err != nil {
		return err
	}

	// Send confirmation
	var confPkt prot.Packet
	confPkt.SerializeToBody(conf, prot.Bool)
	c.Sock.SendEncryptedPacket(confPkt)

	if !conf {
		fmt.Println("Sync aborted...")
		return nil
	}

	return c.ReceiveUniqueFiles(uniqueHashes)
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
		fmt.Printf("Connection established with client (%s)\n", peer.Addr())

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

		// Receive confirmation packet
		var confPkt prot.Packet
		if err := c.Sock.ReceiveEncryptedPacket(&confPkt); err != nil {
			return fmt.Errorf("failed to receive confirmation: %w", err)
		}

		// Check confirmation
		var result bool
		err = confPkt.DeserializeBody(&result)
		if err != nil {
			msg := "unable to deserialize confirmation: " + err.Error()
			return errors.New(msg)
		}

		if result {
			if err := c.SendUniqueFiles(*uniqueFiles); err != nil {
				return err
			}
		} else {
			fmt.Println("Client rejected file transfer...")
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

func (c Client) confirmDownload(uniqueHashes []dir.FileHash) (bool, error) {
	for {
		// fmt.Println("\033[1mFiles\t\t\t\tSize\033[0m")
		totalSize := int64(0)
		for _, file := range uniqueHashes {
			// fmt.Printf("\033[1m%s\t\t\t\t%d\033[0m\n", file.Name, file.Size)
			totalSize += file.Size
		}

		fmt.Printf("\nTotal size: \033[1m%d\033[0m\n", totalSize)
		fmt.Print("Proceed with download? [y/n]: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		switch strings.TrimSpace(input) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		}
	}
}
