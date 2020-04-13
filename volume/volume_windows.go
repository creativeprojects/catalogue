// +build windows

package volume

import (
	"errors"

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
		return err
	}
	volume.BytesTotal = totalNumberOfBytes
	volume.BytesFree = totalNumberOfFreeBytes
	return nil
}
