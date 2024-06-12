package volume

import (
	"bytes"
	"errors"
	"fmt"
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

	err = addMountInfo(uint32(major), uint32(minor), vol)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(vol.Device, "/dev/") {
		vol.VolumeType = DriveVirtual
		return nil
	}

	// Lookup for udevadm - typically containers don't include the tool
	if udevadm, _ := exec.LookPath("udevadm"); udevadm == "" {
		return nil
	}

	// Fill in more information from udevadm tool
	err = getDriveInfo(vol.Device, vol)
	if err != nil {
		return err
	}

	return nil
}

func addMountInfo(major, minor uint32, vol *Volume) error {
	mounts, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(mounts)
	return addMountInfoFromBuffer(buffer, major, minor, vol)
}

func addMountInfoFromBuffer(buffer *bytes.Buffer, major, minor uint32, vol *Volume) error {
	/*
	   See http://man7.org/linux/man-pages/man5/proc.5.html

	   36 35 98:0 /mnt1 /mnt2 rw,noatime master:1 - ext3 /dev/root rw,errors=continue
	   (1)(2)(3)   (4)   (5)      (6)      (7)   (8) (9)   (10)         (11)

	   (1) mount ID:  unique identifier of the mount (may be reused after umount)
	   (2) parent ID:  ID of parent (or of self for the top of the mount tree)
	   (3) major:minor:  value of st_dev for files on filesystem
	   (4) root:  root of the mount within the filesystem
	   (5) mount point:  mount point relative to the process's root
	   (6) mount options:  per mount options
	   (7) optional fields:  zero or more fields of the form "tag[:value]"
	   (8) separator:  marks the end of the optional fields
	   (9) filesystem type:  name of filesystem of the form "type[.subtype]"
	   (10) mount source:  filesystem specific information or "none"
	   (11) super options:  per super block options
	*/
	const (
		keyMountID = iota
		keyParentID
		keyDevice
		keyRootMountPoint
		keyMountPoint
		keyMountOptions
		keyOptionalFields
	)
	const (
		keyFilesystemType = iota
		keyMountSource
		keySuperOptions
	)

	device := fmt.Sprintf("%d:%d", major, minor)

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, " - ")
		if len(parts) != 2 {
			continue
		}

		left := strings.Split(parts[0], " ")
		if len(left) < keyOptionalFields+1 {
			continue
		}
		if left[keyDevice] != device {
			continue
		}
		vol.VolumeID = left[keyMountID]
		vol.Path = left[keyMountPoint]

		right := strings.Split(parts[1], " ")
		if len(right) < keySuperOptions+1 {
			continue
		}
		vol.Format = right[keyFilesystemType]
		vol.Device = right[keyMountSource]
		// we stop at the first match (bind-mounted filesystems may have multiple entries)
		break
	}

	return nil
}

func getDriveInfo(device string, vol *Volume) error {
	var err error

	if device == "" {
		return errors.New("empty device argument")
	}

	cmd := exec.Command("udevadm", "info", "--query=property", "--name="+device)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Output(%v): %v", cmd.Args, err)
	}

	buffer := bytes.NewBuffer(output)
	return getDriveInfoFromBuffer(buffer, vol)
}

func getDriveInfoFromBuffer(buffer *bytes.Buffer, vol *Volume) error {
	diskID := ""
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

		case "ID_FS_LABEL":
			vol.Name = value

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

func getDeviceID(volumePath string, vol *Volume) error {
	if vol == nil {
		return errors.New("Null argument vol")
	}
	stat := &unix.Stat_t{}
	err := unix.Stat(volumePath, stat)
	if err != nil {
		return err
	}
	vol.DeviceID = uint64(stat.Dev)
	return nil
}
