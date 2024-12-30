package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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



    
    return &Repository{path: path}, nil
}
    
func main() {
    // Example usage
    repo, err := InitRepository("./example-repo")
    if err != nil {
        fmt.Printf("Failed to initialize repository: %v\n", err)
        return
    }
    fmt.Println("Repository initialized successfully" , repo)
}
