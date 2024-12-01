package directory

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"

	"github.com/gobwas/glob"
)

type DirManager struct {
	Path string
}

type FileHash struct {
	Name string
	Hash string
}

type fileMatch func(os.DirEntry) bool

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

// Get file hashes from directory.
// Get all file hashes, unless glob pattern is provided.
func (d DirManager) GetFileHashes(filePattern string) ([]FileHash, error) {
	var hashes []FileHash
	var err error

	if isGlobPattern(filePattern) {
		hashes, err = d.GetFileHashesWithGob(filePattern)
		if err != nil {
			return nil, err
		}
	} else if filePattern {

	} else {
		hashes, err = d.GetAllFileHashes()
		if err != nil {
			return nil, err
		}
	}

	return hashes, nil
}

// Get hashes of files that math glob pattern in directory
func (d DirManager) GetFileHashesWithGob(filePattern string) ([]FileHash, error) {
	var hashes []FileHash
	var err error

	var g glob.Glob
	glob.MustCompile(filePattern)

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		entryname := entry.Name()
		if !entry.IsDir() && g.Match(entryname) {
			// Open + read file
			file, err := os.Open(d.Path + "/" + entryname)
			if err != nil {
				return nil, err
			}
			var fileData []byte
			file.Read(fileData)

			// Hash raw file data
			hash := sha256.Sum256(fileData)
			encodedHash := hex.EncodeToString(hash[:])
			hashes = append(hashes, FileHash{entryname, encodedHash})
		}
	}

	return hashes, err
}

// Get hashes of all files in directory
func (d DirManager) GetAllFileHashes() ([]FileHash, error) {
	var hashes []FileHash
	var err error

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			hash, err := d.HashFile(&entry)
			if err != nil {
				return []FileHash{}, err
			}
			hashes = append(hashes, FileHash{entry.Name(), hash})
		}
	}

	return hashes, err
}

func (d DirManager) GetFileHash(fileName string) (FileHash, error) {
	var hash FileHash

	dirEntries, err := os.ReadDir(d.Path)
	if err != nil {
		return hash, nil
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() && entry.Name() == fileName {

		}
	}

	return hash, nil
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

func (d DirManager) HashFile(entry *os.DirEntry) (string, error) {
	// Open + read file
	file, err := os.Open(d.Path + "/" + (*entry).Name())
	if err != nil {
		return "", err
	}
	fileData := make([]byte, 0)
	file.Read(fileData)

	// Hash raw file data
	hash := sha256.Sum256(fileData)
	encodedHash := hex.EncodeToString(hash[:])

	return encodedHash, nil
}

// Check if a hash is found in a slice of hashes
func containsHash(hashes []FileHash, hash FileHash) bool {
	for _, h := range hashes {
		if h == hash {
			return true
		}
	}
	return false
}

func isGlobPattern(s string) bool {
	return strings.ContainsAny(s, "*?[]")
}

func matchSpecificFile(entry os.DirEntry) bool {

}
