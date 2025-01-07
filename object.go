package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//Returns: (string, error)
//String is object hash;
//Purpose: Store compressed objects in .gg/objects
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

func isValidObjectType(objectType ObjectType) bool {
	switch objectType {
		case BlobObject , TreeObject , CommitObject:
			return true
		default:
			return false
	}
}




// Returns: (ObjectType, []byte, error)
// ObjectType is the type of the object;
// Purpose: Reads compressed objects from .gg/objects
func (repo *Repository) ReadObject(hash string) (ObjectType, []byte, error) {

	if !ValidateHash(hash) {
		return "", nil, fmt.Errorf("invalid hash")
	}

	// objectPath is the path to the object file .. first 2 characters of the hash for the 
	// directory and the rest of the hash for the file name that's how git works. 
	objectPath := filepath.Join(repo.path, ".gg/objects", hash[:2], hash[2:])

	// im reading the compressed data 
	compressedData, err := os.ReadFile(objectPath)
	if err != nil {

		if os.IsNotExist(err) {
			return "", nil, fmt.Errorf("object file not found")
		}
		return "", nil, fmt.Errorf("failed to read object file: %v", err)
	}

	const MaxSize = int64(100 * 1024 * 1024) // 100MB ... isse zyada nhi hona mangta

	limitReader := io.LimitReader(bytes.NewReader(compressedData), MaxSize)
	// this is a hack to prevent decompression bombs
	// the limitReader will be closed after the function returns
	raw, err := zlib.NewReader(limitReader)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read object file: %v", err)
	}

	//! defer raw.Close() to close the file until after the function returns
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

	ValidObjType := ObjectType(objtype)

	if !isValidObjectType(ValidObjType) {
		return "", nil, fmt.Errorf("invalid object type")
		
	}
	
	if len(content) != size {
		return "", nil, fmt.Errorf("object file mismatch")
	}

	return ValidObjType, content, nil

}
