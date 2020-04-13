package volume

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanGetDiskSpace(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Current working directory: %s", cwd)
	vol, err := NewVolumeFromPath(cwd)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotZero(t, vol.BytesTotal)
	assert.NotZero(t, vol.BytesFree)
	assert.True(t, vol.BytesTotal >= vol.BytesFree, "Total disk space should be greater or equal thant the free space")
	t.Logf("Total space: %d, Free space: %d", vol.BytesTotal, vol.BytesFree)
}
