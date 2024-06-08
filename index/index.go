package index

import (
	"os"
	"path/filepath"

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
	fileIndexedChannel chan<- FileIndexed
}

func NewIndexer(rootPath string, fileIndexedChannel chan<- FileIndexed) *Indexer {
	return NewFsIndexer(rootPath, fileIndexedChannel, afero.NewOsFs())
}

func NewFsIndexer(rootPath string, fileIndexedChannel chan<- FileIndexed, fs afero.Fs) *Indexer {
	return &Indexer{
		fs:                 fs,
		rootPath:           rootPath,
		fileIndexedChannel: fileIndexedChannel,
	}
}

// Run starts the indexing process. It will walk the filesystem and send the results to the fileIndexedChannel.
// The Run method will return after all files have been indexed.
func (i *Indexer) Run() {
	_ = i.walk(i.rootPath)
}

func (i *Indexer) walk(path string) error {
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
		filename := filepath.Join(path, name)
		fileInfo, err := lstatIfPossible(i.fs, filename)
		if err != nil {
			i.fileIndexedChannel <- FileIndexed{Path: filename, Error: err}
			continue
		}
		i.fileIndexedChannel <- FileIndexed{Path: filename, Info: fileInfo}
		if fileInfo.IsDir() {
			_ = i.walk(filename)
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
