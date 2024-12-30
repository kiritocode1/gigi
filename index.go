package main;

import ( 
	"bytes"
	"encoding/binary"

)


const ( 
	indexSignature = "DIRC"
	indexVersion = 2
)


type TimeSpec struct {
    Seconds int32
    Nanoseconds int32
}


type IndexEntry struct {
	Ctime  TimeSpec
	Mtime  TimeSpec
	Dev    uint32
	Inode  uint32
	Mode   uint32
	Uid    uint32 // for the user
	Gid    uint32
	Size   uint32	
	Sha1   [20]byte

	// i need flags and path 
	Flags uint16
	Path  string

} 

// ! HOW DOES IT KNOWWW W


type Index struct {
	Entries []IndexEntry // ? HOW TF DOES IT KNOWWW 
} 

func writeIndexEntries(entries *IndexEntry, buffer *bytes.Buffer)  error {
	//! write metadata 
	binary.Write(buffer, binary.BigEndian, entries.Ctime)
	binary.Write(buffer, binary.BigEndian, entries.Mtime)
	binary.Write(buffer, binary.BigEndian, entries.Dev)
	binary.Write(buffer, binary.BigEndian, entries.Inode)
	binary.Write(buffer, binary.BigEndian, entries.Mode)
	binary.Write(buffer, binary.BigEndian, entries.Uid)
	binary.Write(buffer, binary.BigEndian, entries.Gid)
	binary.Write(buffer, binary.BigEndian, entries.Size)


	//! write times 
	binary.Write(buffer, binary.BigEndian, entries.Ctime)
	binary.Write(buffer, binary.BigEndian, entries.Mtime)

	//! write sha1 
	buffer.Write(entries.Sha1[:])

	//! write flags and path
	binary.Write(buffer, binary.BigEndian, entries.Flags)
	buffer.WriteString(entries.Path)


    padding := 8 - (buffer.Len() % 8) 
	// we do this because the index file is padded to 8 bytes 
	//and we need to add padding to the end of the file by 8 bytes
    if padding < 8 {
        buffer.Write(make([]byte, padding))
    }

	return nil
};


func ReadIndex(buffer *bytes.Buffer) (IndexEntry, error) { 
	var entry IndexEntry; 

	binary.Read(buffer, binary.BigEndian, &entry.Ctime)
	binary.Read(buffer, binary.BigEndian, &entry.Mtime)

	// read metadata
	binary.Read(buffer, binary.BigEndian, &entry.Dev)
	 binary.Read(buffer, binary.BigEndian, &entry.Inode)
    binary.Read(buffer, binary.BigEndian, &entry.Mode)
    binary.Read(buffer, binary.BigEndian, &entry.Uid)
    binary.Read(buffer, binary.BigEndian, &entry.Gid)
    binary.Read(buffer, binary.BigEndian, &entry.Size)

	// read sha
	buffer.Read(entry.Sha1[:])

	// read flags and path
	binary.Read(buffer, binary.BigEndian, &entry.Flags)

	pathLength := entry.Flags & 0x0fff ; 
	pathBytes := make([]byte, pathLength)
	buffer.Read(pathBytes)
	entry.Path = string(pathBytes)

    padding := 8 - ((62 + pathLength) % 8)
    if padding < 8 {
        buffer.Next(int(padding))
    }

	return entry, nil
}

