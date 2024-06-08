package ui

import (
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pterm/pterm"
)

// DefaultSpinner is the default SpinnerPrinter.
var DefaultSpinner = &SpinnerPrinter{
	Sequence:            []string{"▀ ", " ▀", " ▄", "▄ "},
	Style:               &pterm.ThemeDefault.SpinnerStyle,
	Delay:               time.Millisecond * 200,
	ShowTimer:           true,
	TimerRoundingFactor: time.Second,
	TimerStyle:          &pterm.ThemeDefault.TimerStyle,
	MessageStyle:        &pterm.ThemeDefault.SpinnerTextStyle,
}

// SpinnerPrinter is a loading animation, which can be used if the progress is unknown.
type SpinnerPrinter struct {
	text                atomic.Pointer[string]
	Sequence            []string
	Style               *pterm.Style
	Delay               time.Duration
	MessageStyle        *pterm.Style
	ShowTimer           bool
	TimerRoundingFactor time.Duration
	TimerStyle          *pterm.Style

	IsActive bool

	startedAt       time.Time
	currentSequence string

	Writer io.Writer
}

// UpdateText updates the message of the active SpinnerPrinter.
func (s *SpinnerPrinter) UpdateText(text string) {
	s.text.Store(&text)
}

func (s *SpinnerPrinter) Text() string {
	pointer := s.text.Load()
	if pointer == nil {
		return ""
	}
	return *pointer
}

// Start the SpinnerPrinter.
func (s *SpinnerPrinter) Start() *SpinnerPrinter {
	s.IsActive = true
	s.startedAt = time.Now()

	go func() {
		for s.IsActive {
			for _, seq := range s.Sequence {
				if !s.IsActive {
					continue
				}

				var timer string
				if s.ShowTimer {
					timer = " (" + time.Since(s.startedAt).Round(s.TimerRoundingFactor).String() + ")"
				}
				pterm.Fprinto(s.Writer, s.Style.Sprint(seq)+" "+s.MessageStyle.Sprint(s.Text())+s.TimerStyle.Sprint(timer))
				s.currentSequence = seq
				time.Sleep(s.Delay)
			}
		}
	}()
	return s
}

// Stop terminates the SpinnerPrinter immediately.
// The SpinnerPrinter will not resolve into anything.
func (s *SpinnerPrinter) Stop() {
	if !s.IsActive {
		return
	}
	s.IsActive = false
	pterm.Fprinto(s.Writer, strings.Repeat(" ", pterm.GetTerminalWidth()))
	pterm.Fprinto(s.Writer)
}
