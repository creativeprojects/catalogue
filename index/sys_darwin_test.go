//go:build darwin

package index

import "syscall"

func fileInfoSys(deviceID uint64) any {
	return &syscall.Stat_t{
		Dev: int32(deviceID),
	}
}
