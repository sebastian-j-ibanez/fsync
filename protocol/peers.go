package protocol

import (
	"encoding/json"
	"os"

	"github.com/sebastian-j-ibanez/fsync/util"
)

const peerFile = "peer_data.json"

type Peer struct {
	IP   string
	Port string
}

func (p Peer) Addr() string {
	return p.IP + ":" + p.Port
}

// Append peer to peerFile
func RegisterPeer(p Peer) {
	peers := GetPeers()
	peers = append(peers, p)
	SavePeersToFile(peers)
}

// Return slice of peers from the peerFile
func GetPeers() []Peer {
	file, err := os.ReadFile(peerFile)
	util.CheckError(err)

	peers := []Peer{}
	err = json.Unmarshal([]byte(file), &peers)
	util.CheckError(err)

	return peers
}

func SavePeersToFile(peers []Peer) {
	jsonData, err := json.MarshalIndent(peers, "", "")
	util.CheckError(err)

	err = os.WriteFile(peerFile, jsonData, 0644)
	util.CheckError(err)
}
