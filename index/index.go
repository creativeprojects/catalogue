package index

import (
	"context"
	"io/fs"
	"os"

	"github.com/creativeprojects/catalogue/platform"
	"github.com/creativeprojects/catalogue/volume"
)

type FileIndexed struct {
	Path  string
	Info  os.FileInfo
	Error error
}

type Indexer struct {
	fs                 fs.FS
	deviceID           uint64
	fileIndexedChannel chan<- FileIndexed
}

func NewIndexer(volume *volume.Volume, fileIndexedChannel chan<- FileIndexed) *Indexer {
	return NewFsIndexer(volume, fileIndexedChannel, os.DirFS(volume.Mountpoint))
}

func NewFsIndexer(volume *volume.Volume, fileIndexedChannel chan<- FileIndexed, fs fs.FS) *Indexer {
	return &Indexer{
		fs:                 fs,
		deviceID:           volume.DeviceID,
		fileIndexedChannel: fileIndexedChannel,
	}
}

// Run starts the indexing process. It will walk the filesystem and send the results to the fileIndexedChannel.
// The Run method will return after all files have been indexed.
func (i *Indexer) Run(ctx context.Context) error {
	return fs.WalkDir(i.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err != nil {
			i.fileIndexedChannel <- FileIndexed{Path: path, Error: err}
			return nil
		}
		fileInfo, err := d.Info()
		if err != nil {
			i.fileIndexedChannel <- FileIndexed{Path: path, Error: err}
			return nil
		}
		if !platform.IsWindows() && deviceID(fileInfo) != i.deviceID {
			return fs.SkipDir
		}
		i.fileIndexedChannel <- FileIndexed{Path: path, Info: fileInfo}
		return nil
	})
}
