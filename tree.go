package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"io"
)

// ?  <mode> <name>\0<SHA-1 hash>
type TreeEntry struct {
	Mode uint32
	Name string
	Hash [20]byte
}

type Tree struct {
	entries []TreeEntry
}

func NewTree() *Tree {
	return &Tree{
		entries: make([]TreeEntry, 0),
	}
}
//Helper :  make a new tree entry
func (t *Tree) AddEntry(mode uint32, name string, hash [20]byte) {
	t.entries = append(t.entries, TreeEntry{
		Mode: mode,
		Name: name,
		Hash: hash,
	})
}

// Serialize the tree into a byte array
func (t *Tree) Serialize() []byte {
	sort.Slice(t.entries, func(i, j int) bool {
		return t.entries[i].Name < t.entries[j].Name
	})
	var buffer bytes.Buffer
	for _, entry := range t.entries {
		fmt.Fprintf(&buffer, "%o %s\x00", entry.Mode, entry.Name)
		buffer.Write(entry.Hash[:])
	}
	return buffer.Bytes()
}

func (t *Tree) Type() string {
	return "tree"
}

func (t *Tree) Hash() string {
	ctx := t.Serialize()
	headers := fmt.Sprintf("%s %d\x00", t.Type(), len(ctx))
	h := sha1.New()
	h.Write([]byte(headers))
	h.Write(ctx)

	return hex.EncodeToString(h.Sum(nil))
}
// we need to check the mode , if it is valid, because we need to convert it to uint32
func isValidMode(mode uint32) bool {
    validModes := []uint32{
        0100644, // regular file
        0100755, // executable file
        0040000, // directory
        0120000, // symbolic link
        0160000, // gglink (submodule)
    }
    
    for _, validMode := range validModes {
        if mode == validMode {
            return true
        }
    }
    return false
}




func ParseTree(data []byte) (Tree, error) {
	tree := NewTree()
	buffer := bytes.NewBuffer(data)


	for buffer.Len() > 0 {
		// line by line , we read mode and name
		line, err := buffer.ReadBytes(0)
		if err != nil {
			return *tree, err
		}
		if len(line) <= 1{ 
			return *tree, fmt.Errorf("invalid tree entry")
		}
		//?  remove the null byte in the end , then we can split the line
		line = line[:len(line)-1]
		partOfLine := bytes.SplitN(line, []byte{' '}, 2)
		if len(partOfLine) != 2 {
			return *tree, fmt.Errorf("invalid tree entry")
		}
		// convert the mode to uint32
		mode := uint32(0)
		if _, err := fmt.Sscanf(string(partOfLine[0]), "%o", &mode); err != nil {
			return *tree, fmt.Errorf("invalid tree entry")
		}

		if !isValidMode(mode) {
			return *tree, fmt.Errorf("invalid tree entry")
		}

		// we need to check if the mode  is in the correct format
		name := string(partOfLine[1]); if name=="" { return *tree, fmt.Errorf("invalid tree entry") };
		// hash is of value 20 bytes long, we need to take it and convert it to a byte array
		hash := make([]byte, 20)

		if n, err := io.ReadFull(buffer,hash); err != nil || n != 20 {
			return *tree, fmt.Errorf("invalid tree entry")
		}


		var HashEntry [20]byte
		copy(HashEntry[:], hash)

		tree.AddEntry(mode, name, HashEntry)
	}
	return *tree, nil
}
