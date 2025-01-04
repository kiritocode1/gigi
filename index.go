package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

)

const (
	indexSignature = "DIRC"
	indexVersion   = 2
)

type TimeSpec struct {
	Seconds     int64
	Nanoseconds int64
}

type IndexEntry struct {
	Ctime TimeSpec
	Mtime TimeSpec
	Dev   uint32 
	Inode uint32
	Mode  uint32
	Uid   uint32 // for the user
	Gid   uint32
	Size  uint32
	Sha1  [20]byte

	// i need flags and path
	Flags uint16
	Path  string
}

type Index struct {
	Entries []IndexEntry 
}

func writeIndexEntries(entries *IndexEntry, buffer *bytes.Buffer) error {
	//! write metadata
	binary.Write(buffer, binary.BigEndian, entries.Ctime)
	binary.Write(buffer, binary.BigEndian, entries.Mtime)
	binary.Write(buffer, binary.BigEndian, entries.Dev)
	binary.Write(buffer, binary.BigEndian, entries.Inode)
	binary.Write(buffer, binary.BigEndian, entries.Mode)
	binary.Write(buffer, binary.BigEndian, entries.Uid)
	binary.Write(buffer, binary.BigEndian, entries.Gid)
	binary.Write(buffer, binary.BigEndian, entries.Size)
	

	//! write sha1
	buffer.Write(entries.Sha1[:])

	//! write flags and path
	binary.Write(buffer, binary.BigEndian, entries.Flags)
	buffer.WriteString(entries.Path)

	padding := 8 - (buffer.Len() % 8)
	// we do this because the index file is padded to 8 bytes
	//and we need to add padding to the end of the file by 8 bytes
	if padding < 8 {
		buffer.Write(make([]byte, padding))
	}

	return nil
}

func readIndexEntries(buffer *bytes.Buffer) (IndexEntry, error) {
	var entry IndexEntry

	binary.Read(buffer, binary.BigEndian, &entry.Ctime)
	binary.Read(buffer, binary.BigEndian, &entry.Mtime)

	// read metadata
	binary.Read(buffer, binary.BigEndian, &entry.Dev)
	binary.Read(buffer, binary.BigEndian, &entry.Inode)
	binary.Read(buffer, binary.BigEndian, &entry.Mode)
	binary.Read(buffer, binary.BigEndian, &entry.Uid)
	binary.Read(buffer, binary.BigEndian, &entry.Gid)
	binary.Read(buffer, binary.BigEndian, &entry.Size)

	// read sha
	buffer.Read(entry.Sha1[:])

	// read flags and path
	binary.Read(buffer, binary.BigEndian, &entry.Flags)

	pathLength := entry.Flags & 0x0fff
	pathBytes := make([]byte, pathLength)
	buffer.Read(pathBytes)
	entry.Path = string(pathBytes)

	padding := int(8 - ((62 + pathLength) % 8))
	if padding < 8 {
		buffer.Next(padding)
	}

	return entry, nil
}

func (repo *Repository) WriteIndex(entries *Index) error {

	buffer := new(bytes.Buffer)

	buffer.WriteString(indexSignature)
	binary.Write(buffer, binary.BigEndian, uint32(indexVersion))
	binary.Write(buffer, binary.BigEndian, uint32(len(entries.Entries)))

	for _, entry := range entries.Entries {
		if err := writeIndexEntries(&entry, buffer); err != nil {
			return fmt.Errorf("failed to write index entry: %v", err)
		}
	}

	// now im calculating the sha value
	hash := sha1.New()
	hash.Write(buffer.Bytes())
	checksum := hash.Sum(nil)
	buffer.Write(checksum)

	indexPath := filepath.Join(repo.path, ".gg", "index")

	return os.WriteFile(indexPath, buffer.Bytes(), 0644)
}

// reads the index from the disk. returns an error if the index file is not found or corrupted
func (repo *Repository) ReadIndexFiles() (*Index, error) {

	indexPath := filepath.Join(repo.path, ".gg", "index")

	data, err := os.ReadFile(indexPath)
	if err != nil {

		if os.IsNotExist(err) {
			return &Index{Entries: []IndexEntry{}}, nil
		}
		return nil, fmt.Errorf("failed to read index file: %v", err)
	}

	if len(data) < 12 {
		return nil, fmt.Errorf("invalid index file , too short")
	}

	if string(data[0:4]) != indexSignature {
		return nil, fmt.Errorf("invalid index file signature")
	}

	version := binary.BigEndian.Uint32(data[4:8])

	if version != indexVersion {
		return nil, fmt.Errorf("unsupported index version %d", version)
	}

	numEntries := binary.BigEndian.Uint32(data[8:12])
	contextLengthWorking := len(data) - 20

	h := sha1.New()
	h.Write(data[:contextLengthWorking])
	expectedChecksum := h.Sum(nil)
	actual := data[contextLengthWorking:]

	if !bytes.Equal(expectedChecksum, actual) {
		return nil, fmt.Errorf("invalid index file checksum")
	}

	buffer := bytes.NewBuffer(data[12:contextLengthWorking])

	index := &Index{Entries: make([]IndexEntry, 0, numEntries)}

	for i := 0; i < int(numEntries); i++ {
		entry, err := readIndexEntries(buffer)
		if err != nil {
			return nil, fmt.Errorf("failed to read index entry %d: %v", i, err)
		}
		index.Entries = append(index.Entries, entry)
	}
	return index, nil
}
