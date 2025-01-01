package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//? write object to disk
func (repo *Repository) WriteObject(object ObjectType, content []byte) (string, error) {

	hash, data := HashObject(object, content)

	objectPath := filepath.Join(repo.path, ".gg", "objects", hash[:2], hash[2:])

	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", filepath.Dir(objectPath), err)
	}

	var compressedData bytes.Buffer

	zw := zlib.NewWriter(&compressedData)

	// Write the header
	if _, err := zw.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data: %v", err)
	}
	zw.Close()
	// Write the compressed data
	if err := os.WriteFile(objectPath, compressedData.Bytes(), 0444); err != nil {
		return "", fmt.Errorf("failed to write object file: %v", err)
	}

	return hash, nil
}

func (repo *Repository) ReadObject(hash string) (ObjectType, []byte, error) {

	if !ValidateHash(hash) {
		return "", nil, fmt.Errorf("invalid hash")
	}

	objectPath := filepath.Join(repo.path, ".gg/objects", hash[:2], hash[2:])

	compressedData, err := os.ReadFile(objectPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read object file: %v", err)
	}

	raw, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", nil, fmt.Errorf("failed to read object file: %v", err)
	}

	defer raw.Close()
	var Decompressed bytes.Buffer
	if _, err := io.Copy(&Decompressed, raw); err != nil {
		return "", nil, fmt.Errorf("failed to read object file: %v", err)
	}

	data := Decompressed.Bytes()

	split := bytes.SplitN(data, []byte{0}, 2)

	if len(split) != 2 {
		return "", nil, fmt.Errorf("invalid object file")
	}

	header := string(split[0])
	content := split[1]

	var objtype string
	var size int

	_, err = fmt.Sscanf(header, "%s %d\x00", &objtype, &size)
	if err != nil {
		return "", nil, fmt.Errorf("invalid object file")
	}

	if len(content) != size {
		return "", nil, fmt.Errorf("object file mismatch")
	}

	return ObjectType(objtype), content, nil

}
