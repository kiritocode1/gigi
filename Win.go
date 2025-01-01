// +build windows

package main;




import (
    "os"
    "syscall"
)

func getFileMetadata(fileInfo os.FileInfo) (uint32, uint32, uint32, uint32, uint32) {
    stat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData)
    if !ok {
        return 0, 0, 0, 0, 0
    }
    return uint32(stat.FileAttributes), 0, 0, 0, uint32(fileInfo.Mode())
}

