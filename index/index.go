package index

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
)

type Indexer struct {
	fs       afero.Fs
	rootPath string
}

func NewIndexer(rootPath string) *Indexer {
	return NewFsIndexer(rootPath, afero.NewOsFs())
}

func NewFsIndexer(rootPath string, fs afero.Fs) *Indexer {
	return &Indexer{
		fs:       fs,
		rootPath: rootPath,
	}
}

func (i *Indexer) Run() {
	err := afero.Walk(i.fs, i.rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if errors.Is(err, os.ErrPermission) {
					fmt.Println("Permission denied:", path)
					return nil
				}
				return err
			}
			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
