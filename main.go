package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Object interface {
	Type() string
	Hash() string
	Serialize() []byte
}

type Blob struct {
	content []byte
}

func (b *Blob) Type() string {
	return "blob"
}

func (b *Blob) Hash() string {
	h := sha1.New()
	data := append([]byte(fmt.Sprintf("blob %d\x00", len(b.content))), b.content...)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func (b *Blob) Serialize() []byte {
	return b.content
}

type Repository struct {
	path string
}
func InitRepository(path string) (*Repository, error) {
	gitPath := filepath.Join(path, ".gg")

	// Create basic Git directory structure
	dirs := []string{
		gitPath,
		filepath.Join(gitPath, "objects"),
		filepath.Join(gitPath, "refs/heads"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create HEAD file with ref: refs/heads/main
	headFileContent := []byte("ref: refs/heads/main\n")

	if err := os.WriteFile(filepath.Join(gitPath, "HEAD"), headFileContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write HEAD file: %v", err)
	}

	return &Repository{path: path}, nil
}

func main() {
	// Example usage
	repo, err := InitRepository("./example-repo")
	if err != nil {
		fmt.Printf("Failed to initialize repository: %v\n", err)
		return
	}
	fmt.Println("Repository initialized successfully", repo)

	// test hash
	blob := []byte("Hello, World!")
	fileHash := HashFile(blob)
	fmt.Println("Hash of file: ", fileHash)

	isValid := ValidateHash(fileHash)
	fmt.Println("Is valid hash: ", isValid)

	isNotValid := ValidateHash("invalid-hash")
	fmt.Println("Is not valid hash: ", isNotValid)


}

// ! Add life cycle :
//? git add <file.txt>
//& creates a Blob object  using the HASHFILE() and HASHObject();
//& update the Index ( using Index and IndexEntry);

// ! Commit life cycle :
//? git commit -m "message"
// if there is no index file, create one
//& create a new commit object using the Commit() function
//& create a new tree object using the Tree() function
//& Update the HEAD to point to the new commit
