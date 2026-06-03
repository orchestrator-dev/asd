package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"asd/errors"
	stderrs "errors"
)

type JSONHandler struct{}

func (h *JSONHandler) CanHandle(mime, ext string) bool {
	return ext == ".json" || ext == ".jsonl" || mime == "application/json"
}

func (h *JSONHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {

	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	isJSONL := strings.HasSuffix(strings.ToLower(meta.Name), ".jsonl")

	if isJSONL {
		return h.renderJSONL(w, data, opts)
	}

	return h.renderJSON(w, data, opts)
}

func (h *JSONHandler) renderJSON(w io.Writer, data []byte, opts Options) error {
	var formatted bytes.Buffer
	err := json.Indent(&formatted, data, "", "  ")
	if err != nil {
		return h.handleParseError(data, err)
	}

	return h.highlight(w, formatted.String()+"\n", opts)
}

func (h *JSONHandler) renderJSONL(w io.Writer, data []byte, opts Options) error {
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var formatted bytes.Buffer
		err := json.Indent(&formatted, []byte(line), "", "  ")
		if err != nil {
			return fmt.Errorf("json: %w at line %d: %v", errors.ErrParseFailed, i+1, err)
		}

		err = h.highlight(w, formatted.String()+"\n", opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *JSONHandler) handleParseError(data []byte, err error) error {
	var syntaxErr *json.SyntaxError
	if stderrs.As(err, &syntaxErr) {
		line, col := h.offsetToLineCol(data, syntaxErr.Offset)
		return fmt.Errorf("json: %w at line %d col %d: %v", errors.ErrParseFailed, line, col, err)
	}
	return fmt.Errorf("json: %w: %v", errors.ErrParseFailed, err)
}

func (h *JSONHandler) offsetToLineCol(data []byte, offset int64) (int, int) {
	line := 1
	col := 1
	for i := 0; i < int(offset) && i < len(data); i++ {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

func (h *JSONHandler) highlight(w io.Writer, data string, opts Options) error {
	lexer := lexers.Get("json")
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
