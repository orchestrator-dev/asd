package render

import (
	"io"
	"asd/handlers"
)

// NewWriter wraps the base writer with decorators based on Options.
func NewWriter(w io.Writer, opts handlers.Options) io.Writer {
	// For v0.1.0, we keep it simple. Decorators like pager will be added later.
	// if opts.Plain { w = StripANSI(w) }
	// if opts.Lines { w = LineNumber(w) }
	// if opts.Pager && isTTY(w) { w = pager.New(w) }
	return w
}
