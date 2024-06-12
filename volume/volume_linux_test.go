package volume

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMountInfo(t *testing.T) {
	const testFile = "testdata/mountinfo"
	testCases := []struct {
		major    uint32
		minor    uint32
		expected Volume
	}{
		{major: 259, minor: 1, expected: Volume{VolumeID: "26", Path: "/", Format: "ext4", Device: "/dev/root"}},
		{major: 0, minor: 25, expected: Volume{VolumeID: "29", Path: "/sys", Format: "sysfs", Device: "sysfs"}},
		{major: 7, minor: 1, expected: Volume{VolumeID: "49", Path: "/snap/core18/2826", Format: "squashfs", Device: "/dev/loop1"}},
		{major: 259, minor: 3, expected: Volume{VolumeID: "55", Path: "/boot", Format: "ext4", Device: "/dev/nvme0n1p16"}},
		{major: 259, minor: 2, expected: Volume{VolumeID: "57", Path: "/boot/efi", Format: "vfat", Device: "/dev/nvme0n1p15"}},
		{major: 0, minor: 46, expected: Volume{VolumeID: "66", Path: "/run/user/1000", Format: "tmpfs", Device: "tmpfs"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%d:%d", testCase.major, testCase.minor), func(t *testing.T) {
			mounts, err := os.ReadFile(testFile)
			require.NoError(t, err)
			buffer := bytes.NewBuffer(mounts)
			volume := Volume{}

			err = addMountInfoFromBuffer(buffer, testCase.major, testCase.minor, &volume)
			require.NoError(t, err)
			assert.EqualValues(t, testCase.expected, volume)
		})
	}
}

func TestGetDriveInfo(t *testing.T) {
	testCases := []struct {
		filename string
		expected Volume
	}{
		{
			filename: "testdata/udevadm_boot",
			expected: Volume{
				Device:   "/dev/nvme0n1p16",
				Name:     "BOOT",
				Format:   "ext4",
				VolumeID: "284f1818-7a35-4dca-95ab-0f161f5cb1be",
			},
		},
		{
			filename: "testdata/udevadm_sr0",
			expected: Volume{
				Connection: "scsi",
				Device:     "/dev/sr0",
				Name:       "cidata",
				Format:     "iso9660",
				VolumeID:   "2024-06-11-14-16-06-76",
				VolumeType: DriveOptical,
			},
		},
		{
			filename: "testdata/udevadm_uefi",
			expected: Volume{
				Device:   "/dev/vda15",
				Name:     "UEFI",
				Format:   "vfat",
				VolumeID: "5EEB-8351",
			},
		},
		{
			filename: "testdata/udevadm_loop",
			expected: Volume{
				Device:     "/dev/loop1",
				Format:     "ext4",
				VolumeID:   "2e7ce6b3-5c74-449d-9201-800a370358b1",
				VolumeType: DriveLoopback,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.filename, func(t *testing.T) {
			mounts, err := os.ReadFile(testCase.filename)
			require.NoError(t, err)
			buffer := bytes.NewBuffer(mounts)
			volume := Volume{}

			err = getDriveInfoFromBuffer(buffer, &volume)
			require.NoError(t, err)
			assert.EqualValues(t, testCase.expected, volume)
		})
	}
}
