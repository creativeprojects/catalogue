// +build windows

package volume

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

func getDiskSpace(volumePath string, volume *Volume) error {
	if volume == nil {
		return errors.New("Null argument volume")
	}
	directoryName, _ := windows.UTF16PtrFromString(volumePath)
	var freeBytesAvailableToCaller, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	err := windows.GetDiskFreeSpaceEx(directoryName, &freeBytesAvailableToCaller, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return fmt.Errorf("GetDiskFreeSpaceEx: %v", err)
	}
	volume.BytesTotal = totalNumberOfBytes
	volume.BytesFree = totalNumberOfFreeBytes
	return nil
}

func getFilesystemInfo(volumePath string, vol *Volume) error {
	var err error

	rootPath := filepath.VolumeName(volumePath)
	if rootPath == "" {
		// The path is not rooted, so we need to get the proper root path first
		fileName, _ := windows.UTF16PtrFromString(volumePath)
		var bufferLength uint32 = windows.MAX_LONG_PATH + 1
		volumePathName := make([]uint16, bufferLength)

		err = windows.GetVolumePathName(fileName, &volumePathName[0], bufferLength)
		if err != nil {
			return fmt.Errorf("GetVolumePathName: %v", err)
		}
		rootPath = uint16PtrToString(volumePathName, bufferLength)
	}

	// Validate rootPah with a trailing \
	if !strings.HasSuffix(rootPath, `\`) {
		rootPath += `\`
	}
	rootPathName, _ := windows.UTF16PtrFromString(rootPath)

	var volumeNameSize, volumeNameSerialNumber, maximumComponentLength, fileSystemFlags, fileSystemNameSize uint32

	volumeNameSize = windows.MAX_PATH + 1
	fileSystemNameSize = windows.MAX_PATH + 1

	volumeNameBuffer := make([]uint16, volumeNameSize)
	fileSystemNameBuffer := make([]uint16, fileSystemNameSize)

	err = windows.GetVolumeInformation(
		rootPathName,
		&volumeNameBuffer[0],
		volumeNameSize,
		&volumeNameSerialNumber,
		&maximumComponentLength,
		&fileSystemFlags,
		&fileSystemNameBuffer[0],
		fileSystemNameSize,
	)
	if err != nil {
		return fmt.Errorf("GetVolumeInformation('%s'): %v", rootPath, err)
	}
	vol.Name = uint16PtrToString(volumeNameBuffer, volumeNameSize)
	vol.Format = uint16PtrToString(fileSystemNameBuffer, fileSystemNameSize)
	vol.Path = rootPath

	// Get device GUID. It's ok if it's not available
	vol.Device, _ = getVolumeGUID(rootPath)

	return nil
}

func getVolumeGUID(mountPoint string) (string, error) {
	volumeMountPoint, _ := windows.UTF16PtrFromString(mountPoint)

	var bufferLength uint32 = 51
	volumeName := make([]uint16, bufferLength)

	err := windows.GetVolumeNameForVolumeMountPoint(volumeMountPoint, &volumeName[0], bufferLength)
	if err != nil {
		return "", err
	}
	return uint16PtrToString(volumeName, bufferLength), nil
}

func uint16PtrToString(ptr []uint16, size uint32) string {
	return strings.TrimRight(string(utf16.Decode(ptr[0:size])), " \r\n\x00")
}
