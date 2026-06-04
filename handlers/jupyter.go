package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type JupyterHandler struct{}

func (h *JupyterHandler) CanHandle(mime, ext string) bool {
	return ext == ".ipynb"
}

type JupyterNotebook struct {
	Cells []JupyterCell `json:"cells"`
}

type JupyterCell struct {
	CellType string          `json:"cell_type"`
	Source   interface{}     `json:"source"`
	Outputs  []JupyterOutput `json:"outputs"`
}

type JupyterOutput struct {
	Text interface{}            `json:"text"`
	Data map[string]interface{} `json:"data"`
}

func parseSource(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []interface{}:
		var b strings.Builder
		for _, s := range v {
			if str, ok := s.(string); ok {
				b.WriteString(str)
			}
		}
		return b.String()
	}
	return ""
}

func (h *JupyterHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	var nb JupyterNotebook
	if err := json.NewDecoder(r).Decode(&nb); err != nil {
		return fmt.Errorf("failed to parse jupyter notebook: %v", err)
	}

	theme := opts.Theme
	if theme == "" || theme == "auto" {
		theme = "dracula"
	}

	rMarkdown, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(theme),
		glamour.WithWordWrap(opts.Width),
	)

	cellStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	for i, cell := range nb.Cells {
		src := parseSource(cell.Source)

		if cell.CellType == "markdown" {
			out, err := rMarkdown.Render(src)
			if err == nil {
				fmt.Fprint(w, out)
			} else {
				fmt.Fprintln(w, src)
			}
		} else if cell.CellType == "code" {
			fmt.Fprintf(w, "\n[%d] Code:\n", i+1)

			var outBuf strings.Builder
			if opts.Plain || opts.NoColor {
				outBuf.WriteString(src)
			} else {
				quick.Highlight(&outBuf, src, "python", "terminal256", theme)
			}

			if !opts.Clean && !opts.Plain {
				fmt.Fprintln(w, cellStyle.Render(outBuf.String()))
			} else {
				fmt.Fprintln(w, outBuf.String())
			}

			// Render Outputs
			for _, out := range cell.Outputs {
				if out.Text != nil {
					text := parseSource(out.Text)
					fmt.Fprintf(w, "  Output: %s\n", strings.TrimSpace(text))
				} else if out.Data != nil {
					if textPlain, ok := out.Data["text/plain"]; ok {
						text := parseSource(textPlain)
						fmt.Fprintf(w, "  Output: %s\n", strings.TrimSpace(text))
					}
				}
			}
		}
	}

	return nil
}
