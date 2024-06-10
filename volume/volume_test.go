package volume

import (
	"os"
	"testing"

	"github.com/creativeprojects/catalogue/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanGetVolumeInformation(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	t.Logf("Current working directory: %s", cwd)
	vol, err := NewVolumeFromPath(cwd)
	require.NoError(t, err)

	assert.NotZero(t, vol.BytesTotal)
	assert.NotZero(t, vol.BytesFree)
	assert.GreaterOrEqual(t, vol.BytesTotal, vol.BytesFree)
	t.Logf("Total space: %d, Free space: %d", vol.BytesTotal, vol.BytesFree)
	assert.NotEmpty(t, vol.Format, "Format should not be empty")
	assert.NotEmpty(t, vol.Device, "Device should not be empty")
	assert.NotEmpty(t, vol.Mountpoint, "Mountpoint should not be empty")
	assert.NotEmpty(t, vol.VolumeType.String(), "VolumeType should not be empty")
	t.Logf("Device: %q Format: %q Path: %q Volume type: %q", vol.Device, vol.Format, vol.Path, vol.VolumeType.String())
	if !platform.IsWindows() {
		assert.NotEmpty(t, vol.DeviceID, "DeviceID should not be empty")
	}
}
