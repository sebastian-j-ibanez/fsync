package directory

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

// Get file hashes. Can be:
// 1. All files
// 2. Files that match glob pattern
// 3. File that matches name
func (d DirManager) GetFileHashes(filePattern string) ([]FileHash, error) {
	if filePattern == "" {
		return d.getAllFileHashes()
	} else if isGlobPattern(filePattern) {
		return d.getGlobPatternFileHashes(filePattern)
	} else {
		return d.getSpecificFileHash(filePattern)
	}
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

// Get hashes of files that math glob pattern in directory
func (d DirManager) getGlobPatternFileHashes(filePattern string) ([]FileHash, error) {
	var hashes []FileHash

	if strings.Contains(filePattern, "..") {
		err := errors.New("invalid glob pattern: '..' is not allowed")
		return nil, err
	}

	g := glob.MustCompile(filePattern)

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
			hashBytes := sha256.Sum256(fileData)
			hash := hex.EncodeToString(hashBytes[:])

			info, err := entry.Info()
			if err != nil {
				return nil, err
			}

			fh := FileHash{
				Name: entryname,
				Hash: hash,
				Size: info.Size(),
			}

			hashes = append(hashes, fh)
		}
	}

	return hashes, err
}

// Get file hash slice for single file entry
func (d DirManager) getSpecificFileHash(filePattern string) ([]FileHash, error) {
	entry, err := d.findFileEntry(filePattern)
	if err != nil {
		return nil, err
	}

	hash, err := d.hashFile(entry)
	if err != nil {
		return nil, err
	}

	return []FileHash{hash}, nil
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

	result := FileHash{
		Name: entry.Name(),
		Hash: encodedHash,
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
		if !containsHash(hashesB, hash) {
			*sharedHashes = append(*sharedHashes, hash)
		}
	}

	return sharedHashes
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

// Check if string is a file glob pattern
func isGlobPattern(s string) bool {
	return strings.ContainsAny(s, "*?[]")
}
