//go:build windows

package index

import "syscall"

func fileInfoSys(_ uint64) any {
	return &syscall.Win32FileAttributeData{}
}
