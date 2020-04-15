package fs

import (
	"os"
	"time"
)

// ExtendedFileInfo is an extended stat_t, filled with attributes that are
// supported by most operating systems. The original FileInfo is embedded.
type ExtendedFileInfo struct {
	os.FileInfo

	DeviceID uint64 // ID of device containing the file
	Device   uint64 // Device ID (if this is a device file)
	Size     int64  // file size in byte

	AccessTime time.Time // last access time stamp
	ModTime    time.Time // last (content) modification time stamp
	ChangeTime time.Time // last status change time stamp
}

// ExtendedStat returns an ExtendedFileInfo constructed from the os.FileInfo.
func ExtendedStat(fi os.FileInfo) ExtendedFileInfo {
	if fi == nil {
		panic("os.FileInfo is nil")
	}

	return extendedStat(fi)
}
