package index

import (
	"context"
	"os"
	"path/filepath"

	"github.com/creativeprojects/catalogue/platform"
	"github.com/spf13/afero"
)

type FileIndexed struct {
	Path  string
	Info  os.FileInfo
	Error error
}

type Indexer struct {
	fs                 afero.Fs
	rootPath           string
	deviceID           uint64
	fileIndexedChannel chan<- FileIndexed
}

func NewIndexer(rootPath string, deviceID uint64, fileIndexedChannel chan<- FileIndexed) *Indexer {
	return NewFsIndexer(rootPath, deviceID, fileIndexedChannel, afero.NewOsFs())
}

func NewFsIndexer(rootPath string, deviceID uint64, fileIndexedChannel chan<- FileIndexed, fs afero.Fs) *Indexer {
	return &Indexer{
		fs:                 fs,
		rootPath:           rootPath,
		deviceID:           deviceID,
		fileIndexedChannel: fileIndexedChannel,
	}
}

// Run starts the indexing process. It will walk the filesystem and send the results to the fileIndexedChannel.
// The Run method will return after all files have been indexed.
func (i *Indexer) Run(ctx context.Context) error {
	return i.walk(ctx, i.rootPath)
}

func (i *Indexer) walk(ctx context.Context, path string) error {
	file, err := i.fs.Open(path)
	if err != nil {
		i.fileIndexedChannel <- FileIndexed{Path: path, Error: err}
		return err
	}
	names, err := file.Readdirnames(-1)
	file.Close()
	if err != nil {
		i.fileIndexedChannel <- FileIndexed{Path: path, Error: err}
		return err
	}

	for _, name := range names {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		filename := filepath.Join(path, name)
		fileInfo, err := lstatIfPossible(i.fs, filename)
		if err != nil {
			i.fileIndexedChannel <- FileIndexed{Path: filename, Error: err}
			continue
		}
		if !platform.IsWindows() && platform.DeviceID(fileInfo) != i.deviceID {
			continue
		}
		i.fileIndexedChannel <- FileIndexed{Path: filename, Info: fileInfo}
		if fileInfo.IsDir() {
			_ = i.walk(ctx, filename)
		}
	}
	return nil
}

// if the filesystem supports it, use Lstat, else use fs.Stat
func lstatIfPossible(fs afero.Fs, path string) (os.FileInfo, error) {
	if lfs, ok := fs.(afero.Lstater); ok {
		fi, _, err := lfs.LstatIfPossible(path)
		return fi, err
	}
	return fs.Stat(path)
}
