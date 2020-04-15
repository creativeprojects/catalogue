// +build !windows,!darwin,!freebsd,!netbsd

package fs

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// extendedStat extracts info into an ExtendedFileInfo for unix based operating systems.
func extendedStat(fi os.FileInfo) ExtendedFileInfo {
	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		panic(fmt.Sprintf("conversion to syscall.Stat_t failed, type is %T", fi.Sys()))
	}

	extFI := ExtendedFileInfo{
		FileInfo: fi,
		DeviceID: uint64(s.Dev),
		Device:   uint64(s.Rdev),
		Size:     s.Size,

		AccessTime: time.Unix(s.Atim.Unix()),
		ModTime:    time.Unix(s.Mtim.Unix()),
		ChangeTime: time.Unix(s.Ctim.Unix()),
	}

	return extFI
}
