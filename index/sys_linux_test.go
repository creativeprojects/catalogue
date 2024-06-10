//go:build linux

package index

import "syscall"

func fileInfoSys(deviceID uint64) any {
	return &syscall.Stat_t{
		Dev: deviceID,
	}
}
