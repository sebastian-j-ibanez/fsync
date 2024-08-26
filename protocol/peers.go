package protocol

import (
	"encoding/json"
	"os"
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
func RegisterPeer(p Peer) error {
	peers, err := GetPeers()
	if err != nil {
		return err
	}

	peers = append(peers, p)
	err = SavePeersToFile(peers)
	if err != nil {
		return err
	}

	return nil
}

// Return slice of peers from the peerFile
func GetPeers() ([]Peer, error) {
	peers := []Peer{}
	currFile, err := os.ReadFile(peerFile)

	if os.IsNotExist(err) {
		newFile, err := os.Create(peerFile)
		if err != nil {
			return nil, err
		}

		jsonData, err := json.Marshal(peers)
		if err != nil {
			return nil, err
		}

		_, err = newFile.Write(jsonData)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Ensure json file at least has [], so it can be unmarshalled
		if len(currFile) <= 0 || currFile[0] != '[' {
			firstBrace := []byte{'['}
			currFile = append(firstBrace, currFile...)
			currFile = append(currFile, ']')
		}

		err = json.Unmarshal([]byte(currFile), &peers)
		if err != nil {
			return nil, err
		}
	}

	return peers, nil
}

func SavePeersToFile(peers []Peer) error {
	jsonData, err := json.MarshalIndent(peers, "", "")
	if err != nil {
		return err
	}

	err = os.WriteFile(peerFile, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
