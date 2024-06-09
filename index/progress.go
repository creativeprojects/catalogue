package index

import (
	"fmt"
	"os"

	"github.com/creativeprojects/catalogue/ui"
	"github.com/pterm/pterm"
)

type Progresser interface {
	Start()
	Increment(path string, info os.FileInfo)
	Error(path string, err error)
	Stop(message string)
}

type Progress struct {
	spinner    *ui.SpinnerPrinter
	fileCount  int
	dirCount   int
	errorCount int
}

func NewProgress() *Progress {
	return &Progress{
		spinner: ui.DefaultSpinner,
	}
}

func (p *Progress) Start() {
	p.dirCount = 0
	p.fileCount = 0
	p.errorCount = 0
	p.spinner.Start()
}

func (p *Progress) Increment(path string, info os.FileInfo) {
	if info.IsDir() {
		p.dirCount++
	} else {
		p.fileCount++
	}
	p.update()
}

func (p *Progress) Error(path string, err error) {
	pterm.Error.Println(err)
	p.errorCount++
	p.update()
}

func (p *Progress) Stop(message string) {
	p.spinner.Stop()
	pterm.Success.Println(message)
}

func (p *Progress) update() {
	text := fmt.Sprintf("Files: %d, Directories: %d, Errors: %d", p.fileCount, p.dirCount, p.errorCount)
	p.spinner.UpdateText(text)
}
