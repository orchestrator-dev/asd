package handlers

import (
	"io"
)

// Options holds all explicit configurations passed around.
type Options struct {
	Flat    bool
	Lines   bool
	Plain   bool
	Theme   string
	NoColor bool // auto-set when stdout is not a TTY
	Width   int  // terminal width, auto-detected
	Clean   bool // disable all decorations like headers, borders, line numbers
	Follow  bool // tail -f mode
	Diff    bool // side-by-side diff
	NoPager bool // disable auto-paging
}

// FileMeta carries information about the file being handled.
type FileMeta struct {
	Name string
	Size int64
}

// Handler interface defines a strategy for handling a specific file type.
type Handler interface {
	// CanHandle returns true if this handler owns this MIME+ext combo.
	CanHandle(mime, ext string) bool
	// Render writes output to w. opts carries CLI flags.
	Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error
}
