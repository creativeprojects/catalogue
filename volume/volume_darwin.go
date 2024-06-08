//go:build darwin
// +build darwin

package volume

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/sys/unix"
	"howett.net/plist"
)

type diskutilPartitioInfo struct {
	BusProtocol                    string
	Ejectable                      bool
	Internal                       bool
	Removable                      bool
	RemovableMedia                 bool
	RemovableMediaOrExternalDevice bool
	VolumeName                     string
	VolumeUUID                     string
}

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
	var err error

	if vol == nil {
		return errors.New("Null argument vol")
	}
	stat := &unix.Statfs_t{}
	err = unix.Statfs(volumePath, stat)
	if err != nil {
		return err
	}
	vol.Format = ByteArrayToString(stat.Fstypename[:])
	vol.Path = ByteArrayToString(stat.Mntonname[:])
	vol.Device = ByteArrayToString(stat.Mntfromname[:])

	// If vol.Device starts with // it means it's a network drive
	if strings.HasPrefix(vol.Device, "//") {
		vol.VolumeType = DriveNetwork
		vol.Name = path.Base(vol.Device)
		return nil
	}

	// Fill in more information from diskutil tool
	err = getDriveInfo(vol.Device, vol)
	if err != nil {
		return fmt.Errorf("getDriveInfo %q: %w", vol.Device, err)
	}

	return nil
}

func getDriveInfo(partition string, vol *Volume) error {
	var err error

	if partition == "" {
		return errors.New("Empty partition argument")
	}

	cmd := exec.Command("diskutil", "info", "-plist", partition)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Output(%v): %v", cmd.Args, err)
	}

	buffer := bytes.NewReader(output)
	data := diskutilPartitioInfo{}
	decoder := plist.NewDecoder(buffer)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	vol.Name = data.VolumeName
	vol.VolumeID = data.VolumeUUID
	vol.Connection = data.BusProtocol

	if data.Removable {
		vol.VolumeType = DriveRemovable
		if data.BusProtocol == "Disk Image" {
			vol.VolumeType = DriveLoopback
		}
	} else if data.Internal {
		vol.VolumeType = DriveFixed
	} else if data.Ejectable {
		vol.VolumeType = DriveRemovable
	}

	return nil
}

func ByteArrayToString(byteArray []byte) string {
	n := bytes.Index(byteArray, []byte{0})
	return string(byteArray[:n])
}
