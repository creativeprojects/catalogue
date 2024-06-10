//go:build windows

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

	testCases := []struct {
		name          string
		fs            fstest.MapFS
		expectedFiles []string
	}{
		{
			name:          "empty",
			fs:            fstest.MapFS{},
			expectedFiles: []string{"."},
		},
		{
			name: "one empty dir",
			fs: fstest.MapFS{
				"dir": &fstest.MapFile{Mode: fs.ModeDir},
			},
			expectedFiles: []string{".", "dir"},
		},
		{
			name: "one file",
			fs: fstest.MapFS{
				"file": &fstest.MapFile{},
			},
			expectedFiles: []string{".", "file"},
		},
		{
			name: "one child",
			fs: fstest.MapFS{
				"dir/file": &fstest.MapFile{},
			},
			expectedFiles: []string{".", "dir", "dir/file"},
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
				PathIndex: "C:\\",
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
