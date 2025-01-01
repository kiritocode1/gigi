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

func NewCommit(treeHash string, parentHash []string) *Commit {
	return &Commit{
		TreeHash:   treeHash,
		ParentHash: parentHash,
		Author:     Sign{},
		Committer:  Sign{},
	}
}

func (c *Commit) String() string {
	author := c.Author.String()
	committer := c.Committer.String()
	message := c.Message
	treeHash := c.TreeHash
	parentHash := c.ParentHash
	return fmt.Sprintf("tree %s\nparent %s\nauthor %s\ncommitter %s\n\n%s", treeHash, parentHash, author, committer, message)
}

// the commit format is
// tree <tree hash>
// parent <parent hash>
// author <author>
// committer <committer>
//
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

func ParseCommit(data []byte) (Commit, error) {
	commit := NewCommit("", nil)

	parts := bytes.SplitN(data, []byte{'\n', '\n'}, 2)

	if len(parts) != 2 {
		return *commit, fmt.Errorf("invalid commit")
	}
	key := string(parts[0])
	value := string(parts[1])

	switch key {
	case "tree":
		commit.TreeHash = value
	case "parent":
		commit.ParentHash = append(commit.ParentHash, value)
	case "author":
		author, err := parseSign(value)
		if err != nil {
			return Commit{}, err
		}
		commit.Author = *author
	case "committer":
		committer, err := parseSign(value)
		if err != nil {
			return Commit{}, err
		}
		commit.Committer = *committer
	}

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
