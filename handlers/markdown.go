package handlers

import (
	"io"
	"strings"

	"github.com/charmbracelet/glamour"
)

type MarkdownHandler struct{}

func (h *MarkdownHandler) CanHandle(mime, ext string) bool {
	return strings.HasPrefix(mime, "text/markdown") || ext == ".md"
}

func (h *MarkdownHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	theme := opts.Theme
	if opts.NoColor || opts.Plain {
		theme = "ascii"
	} else if theme == "" {
		theme = "dark" // default theme
	}

	width := opts.Width
	if width <= 0 {
		width = 80 // or termenv size?
	}

	// Maybe glamour handles width, 80 by default. Let's just use width
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(theme),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return err
	}

	out, err := renderer.RenderBytes(data)
	if err != nil {
		return err
	}

	_, err = w.Write(out)
	return err
}
