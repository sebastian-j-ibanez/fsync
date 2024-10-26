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
			hashes = append(hashes, FileHash{name, encodedHash})
		}
	}

	return hashes, err
}

// Return unique values in hashesA but not in hashesB
func GetUniqueHashes(hashesA []FileHash, hashesB []FileHash) *[]FileHash {
	sharedHashes := new([]FileHash)

	for _, hash := range hashesA {
		if !containsHash(hashesB, hash) {
			*sharedHashes = append(*sharedHashes, hash)
		}
	}

	return sharedHashes
}

func containsHash(hashes []FileHash, hash FileHash) bool {
	for _, h := range hashes {
		if h == hash {
			return true
		}
	}
	return false
}
