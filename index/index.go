package index

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

type Indexer struct {
	fs          afero.Fs
	rootPath    string
	progressBar *pterm.ProgressbarPrinter
}

func NewIndexer(rootPath string) *Indexer {
	return NewFsIndexer(rootPath, afero.NewOsFs())
}

func NewFsIndexer(rootPath string, fs afero.Fs) *Indexer {
	pbar, _ := pterm.DefaultProgressbar.WithShowElapsedTime().WithShowCount().WithShowTitle().Start()
	return &Indexer{
		fs:          fs,
		rootPath:    rootPath,
		progressBar: pbar,
	}
}

func (i *Indexer) Run(progresser Progresser) {
	_ = i.walk(i.rootPath)
}

func (i *Indexer) walk(path string) error {
	file, err := i.fs.Open(path)
	if err != nil {
		pterm.Error.Printf("%s: %s\n", path, err)
		return err
	}
	names, err := file.Readdirnames(-1)
	file.Close()
	if err != nil {
		pterm.Error.Printf("%s: %s\n", path, err)
		return err
	}
	i.progressBar.Total += len(names)

	for _, name := range names {
		i.progressBar.Increment()
		filename := filepath.Join(path, name)
		fileInfo, err := lstatIfPossible(i.fs, filename)
		if err != nil {
			pterm.Error.Printf("%s: %s\n", filename, err)
			continue
		}
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
