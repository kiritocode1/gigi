package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)




func Encrypt(filePath string) (string, error) {
	File, err := os.Open(filePath) // open the file
	if err != nil {
		return "", err
	} // if there is an error
	defer File.Close() // close the file
	h := sha1.New()
	if _, err := io.Copy(h, File); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
