package directory

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"slices"
)

type DirManager struct {
	Path string
}

type FileHash struct {
	Name string
	Hash string
	Size int64
}

// Init a new DirManager
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
	if d.Path == "" {
		return nil, errors.New("dir manager is uninitialized")
	}

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range dirEntries {
		names = append(names, entry.Name())
	}

	return names, nil
}

// Get file hashes
func (d DirManager) GetFileHashes(fileNames []string) ([]FileHash, error) {
	if len(fileNames) == 0 {
		return d.getAllFileHashes()
	}

	var hashes []FileHash
	for _, file := range fileNames {
		entry, err := d.findFileEntry(file)
		if err != nil {
			return nil, err
		}

		hash, err := d.hashFile(entry)
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, hash)
	}

	return hashes, nil
}

// Get hashes of all files in directory
func (d DirManager) getAllFileHashes() ([]FileHash, error) {
	var hashes []FileHash
	var err error

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			hash, err := d.hashFile(entry)
			if err != nil {
				return []FileHash{}, err
			}
			hashes = append(hashes, hash)
		}
	}

	return hashes, err
}

// Return the SHA256 hash of file
func (d DirManager) hashFile(entry os.DirEntry) (FileHash, error) {
	// Open + read file
	file, err := os.Open(d.Path + "/" + entry.Name())
	if err != nil {
		return FileHash{}, err
	}
	fileData := make([]byte, 0)
	file.Read(fileData)

	// Hash raw file data
	hash := sha256.Sum256(fileData)
	encodedHash := hex.EncodeToString(hash[:])

	info, err := entry.Info()
	if err != nil {
		return FileHash{}, err
	}

	result := FileHash{
		Name: entry.Name(),
		Hash: encodedHash,
		Size: info.Size(),
	}

	return result, nil
}

// Find file entry in DirManager path with matching name
func (d DirManager) findFileEntry(fileName string) (os.DirEntry, error) {
	if d.Path == "" {
		return nil, errors.New("dir manager is uninitialized")
	}

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	// Return true if file name match is found in cwd
	for _, currEntry := range dirEntries {
		if !currEntry.IsDir() && (currEntry.Name() == fileName) {
			return currEntry, nil
		}
	}

	return nil, errors.New("file not found")
}

// Return unique values in hashesA but not in hashesB
func GetUniqueHashes(hashesA []FileHash, hashesB []FileHash) *[]FileHash {
	sharedHashes := new([]FileHash)

	for _, hash := range hashesA {
		if !slices.Contains(hashesB, hash) {
			*sharedHashes = append(*sharedHashes, hash)
		}
	}

	return sharedHashes
}
