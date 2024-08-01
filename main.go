package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const peerFile = "peer_data.json"

type Peer struct{
	IP string
	Port string
}

func main() {
	// var p []Peer
	// p = append(p, Peer{
	// 	IP: "127.0.0.1",
	// 	Port: "2000",
	// })
	// savePeersToFile(p)

	// p := getPeers()
	// for _, peer := range p {
	// 	fmt.Println(peer)
	// }
}

// Append peer to peerFile
func registerPeer(p Peer) {
	peers := getPeers()
	peers = append(peers, p)
	savePeersToFile(peers)
}

// Return slice of peers from the peerFile
func getPeers() []Peer {
	file, err := os.ReadFile(peerFile)
	checkError(err)
	
	peers := []Peer{}
	err = json.Unmarshal([]byte(file), &peers)
	checkError(err)

	return peers
}

func savePeersToFile(peers []Peer) {
	jsonData, err := json.MarshalIndent(peers, "", "")
	checkError(err)

	err = os.WriteFile(peerFile, jsonData, 0644)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
