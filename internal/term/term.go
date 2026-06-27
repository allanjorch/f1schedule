package term

import (
	"os"
)

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Italic  = "\033[3m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
	Muted   = "\033[38;5;250m"
	BgGreen = "\033[42;30m"
	BgBlue  = "\033[44;97m"
)

type Writer struct {
	enabled bool
}

func NewWriter() *Writer {
	return &Writer{enabled: ColorEnabled()}
}

func ColorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func (w *Writer) Enabled() bool {
	return w.enabled
}

func (w *Writer) Paint(style, text string) string {
	if !w.enabled || text == "" {
		return text
	}
	return style + text + Reset
}

func (w *Writer) Bold(text string) string  { return w.Paint(Bold, text) }
func (w *Writer) Dim(text string) string   { return w.Paint(Dim, text) }
func (w *Writer) Green(text string) string { return w.Paint(Green, text) }
func (w *Writer) Yellow(text string) string { return w.Paint(Yellow, text) }
func (w *Writer) Cyan(text string) string  { return w.Paint(Cyan, text) }
func (w *Writer) Gray(text string) string   { return w.Paint(Gray, text) }
func (w *Writer) Muted(text string) string  { return w.Paint(Muted, text) }
func (w *Writer) Blue(text string) string  { return w.Paint(Blue, text) }
func (w *Writer) Magenta(text string) string { return w.Paint(Magenta, text) }

func (w *Writer) Badge(style, text string) string {
	if !w.enabled {
		return text
	}
	return style + " " + text + " " + Reset
}