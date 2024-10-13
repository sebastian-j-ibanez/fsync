package client

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/mdns"
	prot "github.com/sebastian-j-ibanez/fsync/protocol"
)

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
