//go:build windows

package platform

import "os"

const LineSeparator = "\r\n"

func IsDarwin() bool {
	return false
}

func IsWindows() bool {
	return true
}

func DeviceID(fileInfo os.FileInfo) uint64 {
	return 0
}
