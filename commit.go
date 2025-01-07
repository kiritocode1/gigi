package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

type Sign struct {
	Name  string
	Email string
	Time  time.Time
}

func (s *Sign) String() string {
	timey := s.Time.Format("-0700")
	//? kiritocode1 <kathawalearyan9@gmail.com> unix-timestamp timezone
	return fmt.Sprintf("%s <%s> %d %s", s.Name, s.Email, s.Time.Unix(), timey)
}

// ! to store the commmit author , committer
type Commit struct {
	TreeHash   string
	ParentHash []string
	Author     Sign
	Committer  Sign
	Message    string
}

/*
 * Returns: Commit
 * Purpose: Create a new commit object
*/
func NewCommit(treeHash string, parentHash []string , author Sign , committer Sign , message string) *Commit {
	return &Commit{
		TreeHash:   treeHash,
		ParentHash: parentHash,
		Author:     author,
		Committer:  committer,
		Message:   message,
	}
}

/*
 * Returns: string
 * Purpose: Returns the commit object as a string
 */
func (c *Commit) String() string {
    var parentHashes string
    for _, parent := range c.ParentHash {
        parentHashes += fmt.Sprintf("parent %s\n", parent)
    }
    return fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s", c.TreeHash, parentHashes, c.Author.String(), c.Committer.String(), c.Message)
}


// the commit format is
// tree <tree hash>
// parent <parent hash>
// author <author>
// committer <committer>
func (c *Commit) Serialize() []byte {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "tree %s\n", c.TreeHash)

	//& parent hash
	for _, parentHash := range c.ParentHash {
		fmt.Fprintf(&buffer, "parent %s\n", parentHash)
	}

	fmt.Fprintf(&buffer, "author %s\n", c.Author.String())
	fmt.Fprintf(&buffer, "committer %s\n", c.Committer.String())

	fmt.Fprintf(&buffer, "\n%s", c.Message)

	return buffer.Bytes()
}

/*
 * Returns: string
 * Purpose: Returns the commit object hash
 */
func (c *Commit) Hash() string {
	ctx := c.Serialize()
	headers := fmt.Sprintf("commit %d\x00", len(ctx))
	h := sha1.New()
	h.Write([]byte(headers))
	h.Write(ctx)
	return hex.EncodeToString(h.Sum(nil))
}




// made to parse the data into commit struct
func ParseCommit(data []byte) (Commit, error) {
	// Create a new commit object
	commit := NewCommit("", nil , Sign{} , Sign{} , "")
	// Split the data into two parts 
	//? " <Key> or object id \n message \n"
	parts := bytes.SplitN(data, []byte{'\n', '\n'}, 2)

	if len(parts) != 2 {
		return *commit, fmt.Errorf("invalid commit")
	}

	Key_or_obj_id := string(parts[0])
	Message := string(parts[1])//? message for the commit


	// Check the key or object id 
	// if the key is "tree", then the value is the tree hash
	// if the key is "parent", then the value is the parent hash
	// if the key is "author", then the value is the author sign
	// if the key is "committer", then the value is the committer sign
	switch Key_or_obj_id {
	case "tree":
		commit.TreeHash = Message
	case "parent":
		commit.ParentHash = append(commit.ParentHash, Message)
	case "author":
		author, err := parseSign(Message)
		if err != nil {
			return Commit{}, err
		}
		commit.Author = *author
	case "committer":
		committer, err := parseSign(Message);
		if err != nil {
			return Commit{}, err
		}
		commit.Committer = *committer
	}
	commit.Message = string(parts[1]); 
	return *commit, nil

}

/*
* 	Returns the Sign after parsing it ,
* 	because the sign is in the format of "name <email> time zone"
*/
func parseSign(sign string) (*Sign, error) {
	var name string
	var email string
	var timespace int64
	var timezone string

	//? Scan the sign string into name, email, timespace, and timezone
	n, err := fmt.Sscanf(sign, "%s <%s> %d %s", &name, &email, &timespace, &timezone)
	if err != nil || n < 4 {
		return nil, err
	}

	//! we need to remove the < and > from the email :)
	email = email[1 : len(email)-1]

	return &Sign{
		Name:  name,
		Email: email,
		Time:  time.Unix(timespace, 0),
	}, nil
	
}


/*
 * Returns: (string, error)
 * Purpose: Commits the current tree
 */
func (repo *Repository) Commit(message string, author Sign, committer Sign) (string, error) {


	// Get the current tree hash
	treeHash := repo.GetCurrentTreeHash(); 

	if treeHash == "" {
		return "", fmt.Errorf("failed to get current tree hash")
	}

	//! if there is no commit hash, then we need to get the parent hash from the HEAD file
	//! implementing this one is still pending
	parentHash := repo.GetCurrentCommitHash(); 

	commit := NewCommit(treeHash, []string{parentHash}, author, committer, message); 


	// Serialize the commit object
	// 
	commitData := commit.Serialize();
	commitHash := commit.Hash();

	// Write the commit object to the .gg/objects directory
	_ , err := repo.WriteObject(CommitObject, commitData);
	if err != nil {
		return "", fmt.Errorf("failed to write commit object: %v", err)
	}

	// Update the HEAD file with the new commit hash
	err = repo.UpdateHEAD(commitHash);
	if err != nil {
		return "", fmt.Errorf("failed to update HEAD: %v", err)
	}


	return commitHash, nil
};