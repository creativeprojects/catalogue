package volume

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanGetVolumeInformation(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Current working directory: %s", cwd)
	vol, err := NewVolumeFromPath(cwd)
	if err != nil {
		t.Skip(err)
	}
	assert.NotZero(t, vol.BytesTotal, "BytesTotal should not be zeri")
	assert.NotZero(t, vol.BytesFree, "BytesFree should not be zero")
	assert.True(t, vol.BytesTotal >= vol.BytesFree, "Total disk space should be greater or equal thant the free space")
	t.Logf("Total space: %d, Free space: %d", vol.BytesTotal, vol.BytesFree)
	assert.NotEmpty(t, vol.Format, "Format should not be empty")
	assert.NotEmpty(t, vol.Device, "Device should not be empty")
	assert.NotEmpty(t, vol.Path, "Path should not be empty")
	assert.NotEmpty(t, vol.VolumeType.String(), "VolumeType should not be empty")
	t.Logf("Device: %s Format: %s Mountpoint: %s", vol.Device, vol.Format, vol.Path)
}
