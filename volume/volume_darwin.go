// +build darwin

package volume

import (
	"errors"

	"golang.org/x/sys/unix"
)

func getDiskSpace(volumePath string, volume *Volume) error {
	if volume == nil {
		return errors.New("Null argument volume")
	}
	stat := &unix.Statfs_t{}
	err := unix.Statfs(volumePath, stat)
	if err != nil {
		return err
	}
	volume.BytesTotal = uint64(stat.Bsize) * uint64(stat.Blocks)
	volume.BytesFree = uint64(stat.Bsize) * uint64(stat.Bavail) // Bavail is the space available to a non super user
	return nil
}
