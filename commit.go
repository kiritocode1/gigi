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

func NewCommit(treeHash string, parentHash []string , author Sign , committer Sign , message string) *Commit {
	return &Commit{
		TreeHash:   treeHash,
		ParentHash: parentHash,
		Author:     author,
		Committer:  committer,
		Message:   message,
	}
}

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
	commit := NewCommit("", nil , Sign{} , Sign{} , "")
	parts := bytes.SplitN(data, []byte{'\n', '\n'}, 2)

	if len(parts) != 2 {
		return *commit, fmt.Errorf("invalid commit")
	}
	Key_or_obj_id := string(parts[0])
	Message := string(parts[1])

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
		committer, err := parseSign(Message);if err != nil {
			return Commit{}, err
		}
		commit.Committer = *committer
	}
	commit.Message = string(parts[1]); 
	return *commit, nil

}

func parseSign(sign string) (*Sign, error) {
	var name string
	var email string
	var timespace int64
	var timezone string

	n, err := fmt.Sscanf(sign, "%s <%s> %d %s", &name, &email, &timespace, &timezone)
	if err != nil || n < 4 {
		return nil, err
	}

	email = email[1 : len(email)-1]

	return &Sign{
		Name:  name,
		Email: email,
		Time:  time.Unix(timespace, 0),
	}, nil
	
}



func (repo *Repository) Commit(message string, author Sign, committer Sign) (string, error) {


	treeHash := repo.GetCurrentTreeHash(); 

	if treeHash == "" {
		return "", fmt.Errorf("failed to get current tree hash")
	}

	parentHash := repo.GetCurrentCommitHash(); 

	commit := NewCommit(treeHash, []string{parentHash}, author, committer, message); 


	commitData := commit.Serialize();
	commitHash := commit.Hash();


	_ , err := repo.WriteObject(CommitObject, commitData);
	if err != nil {
		return "", fmt.Errorf("failed to write commit object: %v", err)
	}


	err = repo.UpdateHEAD(commitHash);
	if err != nil {
		return "", fmt.Errorf("failed to update HEAD: %v", err)
	}


	return commitHash, nil
};