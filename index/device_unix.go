//go:build !windows

package index

import (
	"io/fs"
	"syscall"
)

func deviceID(fileInfo fs.FileInfo) uint64 {
	if fileInfo == nil {
		return 0
	}
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Dev)
	}
	return 0
}
