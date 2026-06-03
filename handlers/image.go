package handlers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"asd/exec"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type ImageHandler struct {
	Runner exec.Runner
}

func (h *ImageHandler) CanHandle(mime, ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".svg":
		return true
	}
	return strings.HasPrefix(mime, "image/")
}

func (h *ImageHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	out := w

	if opts.Flat {
		_, err := io.Copy(out, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if h.Runner == nil {
		h.Runner = &exec.OSExecRunner{}
	}

	format := "unknown"
	var width, height int

	isSVG := strings.HasSuffix(strings.ToLower(meta.Name), ".svg")
	var viewBox string
	var elementCount int

	if isSVG {
		format = "svg"
		viewBox, elementCount = parseSVG(bytes.NewReader(data))
	} else {
		config, fmtName, err := image.DecodeConfig(bytes.NewReader(data))
		if err == nil {
			format = fmtName
			width = config.Width
			height = config.Height
		}
	}

	if h.Runner.Available("chafa") && meta.Name != "stdin" && meta.Name != "" {
		output, err := h.Runner.Run("chafa", meta.Name)
		if err == nil {
			fmt.Fprintln(out, strings.TrimRight(string(output), "\n"))
		}
	}

	fmt.Fprintf(out, "File: %s\n", meta.Name)
	if isSVG {
		fmt.Fprintf(out, "Format: %s\n", format)
		fmt.Fprintf(out, "ViewBox: %s\n", viewBox)
		fmt.Fprintf(out, "Elements: %d\n", elementCount)
		fmt.Fprintf(out, "Size: %d bytes\n", meta.Size)
	} else {
		fmt.Fprintf(out, "Format: %s\n", format)
		fmt.Fprintf(out, "Resolution: %dx%d px\n", width, height)
		fmt.Fprintf(out, "Size: %d bytes\n", meta.Size)
	}

	return nil
}

func parseSVG(r io.Reader) (viewBox string, count int) {
	decoder := xml.NewDecoder(r)
	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			count++
			if strings.ToLower(se.Name.Local) == "svg" {
				for _, attr := range se.Attr {
					if strings.ToLower(attr.Name.Local) == "viewbox" {
						viewBox = attr.Value
					}
				}
			}
		}
	}
	return
}
