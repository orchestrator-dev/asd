package handlers

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"asd/errors"
)

type TOMLHandler struct{}

func (h *TOMLHandler) CanHandle(mime, ext string) bool {
	return ext == ".toml" || mime == "application/toml"
}

func (h *TOMLHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {

	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var dummy interface{}
	err = toml.Unmarshal(data, &dummy)
	if err != nil {
		return fmt.Errorf("toml: %w: %v", errors.ErrParseFailed, err)
	}

	formatted, err := toml.Marshal(&dummy)
	if err != nil {
		return err
	}

	return h.highlight(w, string(formatted), opts)
}

func (h *TOMLHandler) highlight(w io.Writer, data string, opts Options) error {
	lexer := lexers.Get("toml")
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

	iterator, err := lexer.Tokenise(nil, data)
	if err != nil {
		return err
	}

	return formatter.Format(w, style, iterator)
}
