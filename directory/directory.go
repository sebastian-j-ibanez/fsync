package directory

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

type DirManager struct {
	Path string
}

type FileHash struct {
	Name string
	Hash string
}

// Hash Directory
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
