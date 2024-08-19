package directory

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
)

type DirManager struct {
	Path string
}

type FileHash struct {
	Name string
	Hash string
}

func NewDirManager(path string) (*DirManager, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	d := &DirManager{
		Path: path,
	}
	
	return d, nil
}

// Get file names in directory
func (d DirManager) GetFileNames() ([]string, error) {
	var names []string
	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		names = append(names, entry.Name())
	}

	return names, nil
}

// Hash files in directory
func (d DirManager) HashDir() ([]FileHash, error) {
	var hashes []FileHash
	var err error

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
