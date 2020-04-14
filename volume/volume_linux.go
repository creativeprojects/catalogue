// +build linux

package volume

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	fsdevice int = iota
	fsmount
	fstype
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

func getFilesystemInfo(volumePath string, vol *Volume) error {
	mounts, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return err
	}

	// This is kind of hacky (until I find a better way?)
	// => go through all the mounts and keep the longest path that prefixes the volumePath
	// It should be the volume where our path sits
	foundMount, foundDevice, foundFormat, foundLen := "", "", "", 0
	buffer := bytes.NewBuffer(mounts)
	for {
		line, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		entry := strings.SplitN(string(line), " ", 4)
		// Don't need to bother with these fs
		if entry[fstype] == "proc" || entry[fstype] == "cgroup" || entry[fstype] == "mqueue" || entry[fstype] == "devpts" || entry[fstype] == "sysfs" {
			continue
		}
		if strings.HasPrefix(volumePath, entry[fsmount]) && len(entry[fsmount]) > foundLen {
			foundLen = len(entry[fsmount])
			foundMount = entry[fsmount]
			foundDevice = entry[fsdevice]
			foundFormat = entry[fstype]
		}
	}
	vol.Path = foundMount
	vol.Device = foundDevice
	vol.Format = foundFormat

	return nil
}
