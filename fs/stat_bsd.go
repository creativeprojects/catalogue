// +build freebsd darwin netbsd

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

		AccessTime: time.Unix(s.Atimespec.Unix()),
		ModTime:    time.Unix(s.Mtimespec.Unix()),
		ChangeTime: time.Unix(s.Ctimespec.Unix()),
	}

	return extFI
}