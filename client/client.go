package client

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/mdns"
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
	var err error
	cltHashes, err := c.DirMan.HashDir()
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
		c.Sock = prot.NewSocketHandler(conn)

		// Get peer file hashes
		peerHashes, err := c.Sock.ReceiveFileHashes()
		if err != nil {
			msg := "unable to receive file hashes: " + err.Error()
			return errors.New(msg)
		}

		// Send unique file hashes
		uniqueFiles := dir.GetUniqueHashes(cltHashes, peerHashes)
		err = c.Sock.SendFileHashes(*uniqueFiles)
		if err != nil {
			msg := "unable to send file hashes: " + err.Error()
			return errors.New(msg)
		}

		// Upload unique files
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

// Await sync from peer over default port
func (c Client) AwaitSync(portNum int) error {
	// Set port
	if portNum == 0 {
		portNum = defaultPort
	}
	port := ": " + strconv.Itoa(portNum)

	var err error
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

	// Download unique files
	for _, file := range uniqueHashes {
		err := c.Sock.DownloadFile(file.Name)
		if err != nil {
			msg := "unable to download " + file.Name + ": " + err.Error()
			return errors.New(msg)
		}
	}

	return nil
}

// Discover mDNS fsync service
func DiscoverMDNSService() (prot.Peer, error) {
	entryCh := make(chan *mdns.ServiceEntry, 4)
	defer close(entryCh)

	go func() {
		mdns.Lookup(serviceName, entryCh)
	}()

	select {
	case service := <-entryCh:
		if service != nil {
			ip := service.AddrV4.String()
			port := strconv.Itoa(service.Port)
			peer := prot.Peer{
				IP:   ip,
				Port: port,
			}

			return peer, nil
		} else {
			return prot.Peer{}, errors.New("unable to receive peer info")
		}
	case <-time.After(timeout):
		return prot.Peer{}, errors.New("timed out searching for peer")
	}
}

// Start mDNS service for other peers to connect to
func BroadcastMDNSService(port int, endBroadcast <-chan bool) error {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"peer-to-peer file syncing"}
	service, err := mdns.NewMDNSService(host, serviceName, "", "", port, nil, info)
	if err != nil {
		return fmt.Errorf("unable to create mDNS service: %w", err)
	}

	// Create the mDNS server
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return fmt.Errorf("failed to create mDNS server: %w", err)
	}
	defer server.Shutdown()

	// Check for end broadcast flag
	for {
		select {
		case <-endBroadcast:
			return nil
		default:
			time.Sleep(2 * time.Second)
		}
	}
}
