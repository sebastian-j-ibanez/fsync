package protocol

import (
	"os"
	"encoding/json"
	"fsyncprotocol/util"
)

const peerFile = "peer_data.json"

type Peer struct{
	IP string
	Port string
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
