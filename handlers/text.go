package handlers

import (
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type TextHandler struct{}

func (h *TextHandler) CanHandle(mime, ext string) bool {
	return mime == "text/plain" || strings.HasPrefix(mime, "text/")
}

func (h *TextHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {

	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	lexer := lexers.Match(meta.Name)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(opts.Theme)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if opts.NoColor || opts.Plain {
		formatter = formatters.Get("noop")
	}

	iterator, err := lexer.Tokenise(nil, string(data))
	if err != nil {
		return err
	}

	return formatter.Format(w, style, iterator)
}
