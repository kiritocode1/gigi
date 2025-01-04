package main;

import ( 
	"os"
	"path/filepath"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
)



func (repo *Repository) GetCurrentTreeHash() string{
	index , err := repo.ReadIndexFiles(); 
	if err != nil {
		fmt.Printf("Failed to read index: %v\n", err)
		return "";
	}

	tree := NewTree()
	for _, entry := range index.Entries {
		tree.AddEntry(entry.Mode, entry.Path, entry.Sha1)
	}

	// we need to serialize the tree and hash it , fir yeh return krr denge 
	treeData := tree.Serialize()
	treeHash := sha1.Sum(treeData)
	return hex.EncodeToString(treeHash[:])
}



func (repo *Repository) GetCurrentCommitHash() string{
	headPath := filepath.Join(repo.path, ".gg", "HEAD")
	headContent, err := os.ReadFile(headPath);
	if err != nil {
		return "";
	}
	//! wtf go :) 
	return string(headContent)
}



func (repo *Repository) UpdateHEAD(commitHash string) error{
	headPath := filepath.Join(repo.path, ".gg", "HEAD")
	return os.WriteFile(headPath, []byte(commitHash), 0644)
}