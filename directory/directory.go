package directory

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	prot "fsync/protocol"
	"os"
)

type DirManager struct {
	Path string
}

type FileHash struct {
	Name string
	Hash string
}

// Convert file to packets
func (d DirManager) PacketizeFile(name string) ([]prot.Packet, error) {
	// Check if dirmanager has valid path
	if err := d.checkPath(); err != nil {
		return nil, err
	}

	// Read from file
	fullPath := d.Path + "/" + name
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
 	
	// Calculate number of packets
	pktNum := int(1)
	fileSize := len(fileData)
	if fileSize >= prot.MaxBodySize {
		pktNum = int(fileSize) / int(prot.MaxBodySize)
	}

	// Convert file data to consistent-sized chunks 
	data, err := chunkBuffer(fileData, prot.MaxBodySize)
	if err != nil {
		return nil, err
	}
	
	// Create and append packets
	var packets []prot.Packet
	for i := range pktNum {
		tmpPacket := prot.Packet{
			OrderNum: i,
			Body: data[i],
		}
		packets = append(packets, tmpPacket)
	}
	
	return packets, nil
}

// Write data to file
func (d DirManager) WriteDataToFile(name string, data []byte) error {
	// Check if dirmanager has valid path
	if err := d.checkPath(); err != nil {
		return err
	}

	// Open file
	fullPath := d.Path + "/" + name
	if err := os.WriteFile(fullPath, data, 0666); err != nil {
		return err
	}
	
	return nil
}

// Hash files in directory
func (d DirManager) HashDir() ([]FileHash, error) {
	var hashes []FileHash
	var err error

	err = d.checkPath()
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			// Open + read file
			name := entry.Name()
			file, err := os.Open(d.Path + "/" + name)
			if err != nil {
				return nil, err
			}
			var fileData []byte
			file.Read(fileData)

			// Hash raw file data
			hash := sha256.Sum256(fileData)
			encodedHash := hex.EncodeToString(hash[:])
			hashes = append(hashes, FileHash{ name, encodedHash })
		}
	}
	
	return hashes, err
}

// Return error if DirManager does not have a valid path
func (d DirManager) checkPath() error {
	_, err := os.Stat(d.Path)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	return nil
}

// Convert byte slice into smaller evenly-sized byte slices (taking into account leftover any bytes)
func chunkBuffer(buffer []byte, chunkSize int) ([][]byte, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must be greater than 0")
	}

	var chunks [][]byte
	dataLen := len(buffer)

	// Iterate over the data and split it into chunks
	for start := 0; start < dataLen; start += chunkSize {
		end := start + chunkSize
		if end > dataLen {
			end = dataLen
		}
		chunks = append(chunks, buffer[start:end])
	}
	
	return chunks, nil
}
