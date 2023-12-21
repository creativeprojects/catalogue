//go:build linux
// +build linux

package volume

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"golang.org/x/sys/unix"
)

const (
	fsdevice int = iota
	fsmount
	fstype
)

const deviceMajorShift = 8
const deviceMinorMask = (1 << deviceMajorShift) - 1

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

	stat, err := os.Stat(volumePath)
	if err != nil {
		return err
	}

	// stat.Dev contain the device IDs (major and minor) so we can search for it in /dev
	sys := stat.Sys().(*syscall.Stat_t)
	major := sys.Dev >> deviceMajorShift
	minor := sys.Dev & deviceMinorMask

	devices, err := ioutil.ReadDir("/dev")
	if err != nil {
		return err
	}
	partition := ""
	for _, fileInfo := range devices {
		fileSys := fileInfo.Sys().(*syscall.Stat_t)
		// For device files, the device IDs are available from Rdev
		if fileSys.Rdev == sys.Dev {
			partition = fileInfo.Name()
		}
	}
	if partition == "" {
		return fmt.Errorf("Cannot find device %d:%d in /dev", major, minor)
	}
	pterm.Debug.Printf("File is on device %d:%d (%s)\n", major, minor, partition)

	vol.Device = "/dev/" + partition
	vol.Path, vol.Format = getDriveMount(vol.Device)

	// Lookup for udevadm - typically containers don't include the tool
	if udevadm, _ := exec.LookPath("udevadm"); udevadm == "" {
		return nil
	}

	// Fill in more information from udevadm tool
	err = getDriveInfo(partition, vol)
	if err != nil {
		return err
	}

	return nil
}

func getDriveInfo(partition string, vol *Volume) error {
	var err error

	if partition == "" {
		return errors.New("Empty partition argument")
	}

	cmd := exec.Command("udevadm", "info", "--query=property", "--name="+partition)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Output(%v): %v", cmd.Args, err)
	}

	diskID := ""
	buffer := bytes.NewBuffer(output)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		keyValuePair := strings.SplitN(line, "=", 2)
		if len(keyValuePair) < 2 {
			continue
		}
		value := strings.TrimSpace(keyValuePair[1])
		switch keyValuePair[0] {
		case "DEVNAME":
			vol.Device = value

		case "ID_CDROM":
			if value == "1" {
				vol.VolumeType = DriveOptical
			}

		case "ID_FS_TYPE":
			vol.Format = value

		case "ID_FS_UUID":
			vol.VolumeID = value

		case "ID_PART_ENTRY_DISK":
			diskID = value
		}
	}

	if diskID != "" {
		pterm.Debug.Printf("Found partition from disk %q\n", diskID)
	}

	return nil
}

func getDriveMount(partition string) (string, string) {
	mounts, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return "", ""
	}
	mountPoints := make([]string, 0)
	filesystem := ""
	buffer := bytes.NewBuffer(mounts)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 3 {
			continue
		}
		if parts[fsdevice] == partition {
			mountPoints = append(mountPoints, parts[fsmount])
			if filesystem == "" {
				filesystem = parts[fstype]
			}
			if parts[fstype] != filesystem {
				// Not sure it would actually happen?
				filesystem += ", " + parts[fstype]
			}
		}
	}

	return strings.Join(mountPoints, ", "), filesystem
}
