// +build darwin

package volume

import (
	"errors"

	"golang.org/x/sys/unix"
)

func getDiskSpace(volumePath string, vol *Volume) error {
	if vol == nil {
		return errors.New("Null argument vol")
	}
	stat := &unix.Statfs_t{}
	err := unix.Statfs(volumePath, stat)
	if err != nil {
		return err
	}
	vol.BytesTotal = uint64(stat.Bsize) * uint64(stat.Blocks)
	vol.BytesFree = uint64(stat.Bsize) * uint64(stat.Bavail) // Bavail is the space available to a non super user
	return nil
}

func getFilesystemInfo(volumePath string, vol *Volume) error {
	if vol == nil {
		return errors.New("Null argument vol")
	}
	stat := &unix.Statfs_t{}
	err := unix.Statfs(volumePath, stat)
	if err != nil {
		return err
	}
	vol.Format = string(int8ToBytes(stat.Fstypename[:]))
	vol.Path = string(int8ToBytes(stat.Mntonname[:]))
	vol.Device = string(int8ToBytes(stat.Mntfromname[:]))
	return nil
}

func int8ToBytes(value []int8) []byte {
	output := make([]byte, len(value))
	for index, char := range value {
		output[index] = byte(char)
	}
	return output
}
