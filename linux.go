//go:build linux || darwin
// +build linux darwin

// go : +build linux darwin
package main

import (
	"os"
	"syscall"
)

func getFileMetadata(fileInfo os.FileInfo) (uint32, uint32, uint32, uint32, uint32) {
    stat, ok := fileInfo.Sys().(*syscall.Stat_t)
    if !ok {
        return 0, 0, 0, 0, 0
    }
    return uint32(stat.Dev), uint32(stat.Ino), uint32(stat.Uid), uint32(stat.Gid), uint32(fileInfo.Mode())
}