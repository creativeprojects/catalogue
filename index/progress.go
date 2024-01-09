package index

import (
	"os"
	"strings"

	"github.com/pterm/pterm"
)

type Progresser interface {
	Increment(path string, info os.FileInfo, fileCount, dirCount int)
	Error(path string, err error)
	Success(message string)
}

type Progress struct {
	spinner *pterm.SpinnerPrinter
}

func NewProgress(spinner *pterm.SpinnerPrinter) *Progress {
	if spinner == nil {
		spinner = &pterm.DefaultSpinner
	}
	spinner.Start()
	return &Progress{
		spinner: spinner,
	}
}

func (p *Progress) Increment(path string, info os.FileInfo, fileCount, dirCount int) {
	maxWidth := pterm.GetTerminalWidth() - 5
	display := path
	if len(display) > maxWidth {
		display = display[len(display)-maxWidth:]
	}
	if len(display) < maxWidth { // pad with spaces
		display = display + strings.Repeat(" ", maxWidth-len(display))
	}
	p.spinner.Text = display
	p.spinner.WithText(display)
}

func (p *Progress) Error(path string, err error) {
	p.spinner.FailPrinter.Printf("%s: %s\n", path, err)
}

func (p *Progress) Success(message string) {
	p.spinner.Success(message)
	_ = p.spinner.Stop()
}
