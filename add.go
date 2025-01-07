package main

import (
	"encoding/hex"
	"fmt"
	"os"
)


//  the path of a file --> adds it in the  
func (repo *Repository) AddFile(path string) error {
	// 1. Read file contents
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 2. Create blob and get hash
	hash := HashFile(contents)

	// 3. Write object to .gg/objects
	_, err = repo.WriteObject(BlobObject, contents)
	if err != nil {
		return fmt.Errorf("failed to write object: %v", err)
	}

	// 4. Read current index
	index, err := repo.ReadIndexFiles()
	if err != nil {
		return fmt.Errorf("failed to read index: %v", err)
	}

	// 5. Get file metadata
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to read file info: %v", err)
	}

	dev, inode, uid, gid, mode := getFileMetadata(fileInfo)

	// 6. Create new index entry
	entry := IndexEntry{
		Ctime: TimeSpec{
			Seconds:     fileInfo.ModTime().Unix(),
			Nanoseconds: fileInfo.ModTime().UnixNano(),
		},
		Mtime: TimeSpec{
			Seconds:     fileInfo.ModTime().Unix(),
			Nanoseconds: fileInfo.ModTime().UnixNano(),
		},
		Dev:   dev,
		Inode: inode,
		Mode:  mode,
		Uid:   uid,
		Gid:   gid,
		Size:  uint32(fileInfo.Size()),
		Sha1:  [20]byte{},        // Need to convert hash string to [20]byte
		Flags: uint16(len(path)), // Flags should include path length
		Path:  path,
	}

	// Convert hash string to [20]byte
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return fmt.Errorf("failed to decode hash: %v", err)
	}
	copy(entry.Sha1[:], hashBytes)

	// 7. Add entry to index and write it
	index.Entries = append(index.Entries, entry)
	return repo.WriteIndex(index)
}