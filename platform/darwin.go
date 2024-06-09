//go:build darwin

package platform

import (
	"os"
	"syscall"
)

const LineSeparator = "\n"

func IsDarwin() bool {
	return true
}

func IsWindows() bool {
	return false
}

func DeviceID(fileInfo os.FileInfo) uint64 {
	if fileInfo == nil {
		return 0
	}
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Dev)
	}
	return 0
}
