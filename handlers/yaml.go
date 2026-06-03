package handlers

import (
	yaml "gopkg.in/yaml.v3"

	"bytes"
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"asd/errors"
)

type YAMLHandler struct{}

func (h *YAMLHandler) CanHandle(mime, ext string) bool {
	return ext == ".yaml" || ext == ".yml" || mime == "application/yaml" || mime == "text/yaml"
}

func (h *YAMLHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {

	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var formatted bytes.Buffer
	dec := yaml.NewDecoder(bytes.NewReader(data))
	enc := yaml.NewEncoder(&formatted)
	enc.SetIndent(2)

	for {
		var node interface{}
		err := dec.Decode(&node)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("yaml: %w: %v", errors.ErrParseFailed, err)
		}
		err = enc.Encode(&node)
		if err != nil {
			return err
		}
	}
	enc.Close()

	return h.highlight(w, formatted.String(), opts)
}

func (h *YAMLHandler) highlight(w io.Writer, data string, opts Options) error {
	lexer := lexers.Get("yaml")
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
