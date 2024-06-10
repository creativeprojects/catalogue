//go:build !windows

package index

import (
	"context"
	"io/fs"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/creativeprojects/catalogue/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexing(t *testing.T) {
	t.Parallel()

	const deviceID = 11
	testCases := []struct {
		name          string
		fs            fstest.MapFS
		expectedFiles []string
	}{
		{
			name: "empty",
			fs: fstest.MapFS{
				".": &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
			},
			expectedFiles: []string{"."},
		},
		{
			name: "one empty dir",
			fs: fstest.MapFS{
				".":   &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
				"dir": &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
			},
			expectedFiles: []string{".", "dir"},
		},
		{
			name: "one file",
			fs: fstest.MapFS{
				".":    &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
				"file": &fstest.MapFile{Sys: fileInfoSys(deviceID)},
			},
			expectedFiles: []string{".", "file"},
		},
		{
			name: "one child",
			fs: fstest.MapFS{
				".":        &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
				"dir":      &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
				"dir/file": &fstest.MapFile{Sys: fileInfoSys(deviceID)},
			},
			expectedFiles: []string{".", "dir", "dir/file"},
		},
		{
			name: "traversing other mounted device",
			fs: fstest.MapFS{
				".":        &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID)},
				"dir":      &fstest.MapFile{Mode: fs.ModeDir, Sys: fileInfoSys(deviceID + 11)},
				"dir/file": &fstest.MapFile{Sys: fileInfoSys(deviceID + 11)},
			},
			expectedFiles: []string{"."},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			indexed := make([]string, 0, 100)
			infoChannel := make(chan FileIndexed, 100)
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func(infoChannel <-chan FileIndexed) {
				defer wg.Done()
				for info := range infoChannel {
					indexed = append(indexed, info.Path)
				}
			}(infoChannel)

			volume := &volume.Volume{
				Mountpoint: "/",
				DeviceID:   deviceID,
			}
			indexer := NewFsIndexer(volume, infoChannel, testCase.fs)
			err := indexer.Run(context.Background())
			require.NoError(t, err)

			close(infoChannel)
			wg.Wait()

			assert.ElementsMatch(t, testCase.expectedFiles, indexed)
		})
	}
}
