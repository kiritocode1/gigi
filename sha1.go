package main

import (
	"crypto/sha1"
	"fmt"
)

type ObjectType string

const (
	BlobObject   ObjectType = "blob"
	TreeObject   ObjectType = "tree"
	CommitObject ObjectType = "commit"
)

func HashContent(objectType ObjectType, content []byte) string {
	// Create header
	//~ "{type} {content_size}\0{content}"
	header := fmt.Sprintf("%s %d\x00", objectType, len(content))

	// Calculate hash
	h := sha1.New()
	h.Write([]byte(header))
	h.Write(content)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// ? uint8 because all the hashes are 32 bytes long
func HashFile(content []byte) string {
	return HashContent(BlobObject, content)
}

// ? ValidateHash checks if a string is a valid SHA-1 hash
func ValidateHash(hash string) bool {
	if len(hash) != 40 {
		return false
	}

	// Check if string is valid hexadecimal
	for _, r := range hash {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}

	return true
}

// HashObject is a generic function to hash any Git object
func HashObject(objectType ObjectType, data []byte) (string, []byte) {
	// Prepare the content with header
	header := fmt.Sprintf("%s %d\x00", objectType, len(data))
	content := append([]byte(header), data...)

	// Calculate hash
	h := sha1.New()
	h.Write(content)

	return fmt.Sprintf("%x", h.Sum(nil)), content
}
