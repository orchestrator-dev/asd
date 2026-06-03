package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"

	"asd/exec"

	"golang.org/x/image/draw"
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

	var img image.Image

	if isSVG {
		format = "svg"
		viewBox, elementCount = parseSVG(bytes.NewReader(data))
	} else {
		var err error
		img, format, err = image.Decode(bytes.NewReader(data))
		if err == nil {
			bounds := img.Bounds()
			width = bounds.Dx()
			height = bounds.Dy()
		}
	}

	renderedHighRes := false

	// Attempt iTerm2 protocol first for pixel-perfect images
	termProg := os.Getenv("TERM_PROGRAM")
	if termProg == "iTerm.app" || termProg == "WezTerm" || termProg == "ghostty" {
		b64 := base64.StdEncoding.EncodeToString(data)
		// Set width to auto to fit nicely if it's too big
		fmt.Fprintf(out, "\033]1337;File=inline=1;width=auto;preserveAspectRatio=1:%s\a\n", b64)
		renderedHighRes = true
	}

	// Attempt chafa if available and not already rendered
	if !renderedHighRes && h.Runner.Available("chafa") && meta.Name != "stdin" && meta.Name != "" {
		output, err := h.Runner.Run("chafa", meta.Name)
		if err == nil {
			fmt.Fprintln(out, strings.TrimRight(string(output), "\n"))
			renderedHighRes = true
		}
	}

	// Fallback to our own native ANSI block renderer
	if !renderedHighRes && img != nil && !opts.NoColor && !opts.Plain {
		renderANSIImage(out, img, opts.Width)
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

	if !renderedHighRes && !opts.NoColor && !opts.Plain {
		fmt.Fprintf(out, "\n\x1b[90mTip: Use WezTerm, iTerm2, or Ghostty for high-resolution images.\x1b[0m\n")
	}

	return nil
}

func renderANSIImage(w io.Writer, img image.Image, termWidth int) {
	if termWidth <= 0 {
		termWidth = 80
	}

	bounds := img.Bounds()
	imgW := bounds.Dx()
	imgH := bounds.Dy()
	if imgW == 0 || imgH == 0 {
		return
	}

	targetW := imgW
	// Cap the width to termWidth, minus a little margin
	if targetW > termWidth-4 {
		targetW = termWidth - 4
	}

	scale := float64(targetW) / float64(imgW)
	targetH := int(float64(imgH) * scale)

	if targetW <= 0 || targetH <= 0 {
		return
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	for y := 0; y < targetH; y += 2 {
		for x := 0; x < targetW; x++ {
			top := dst.RGBAAt(x, y)

			var bottomR, bottomG, bottomB uint8
			if y+1 < targetH {
				bottom := dst.RGBAAt(x, y+1)
				bottomR, bottomG, bottomB = bottom.R, bottom.G, bottom.B
			}

			// Use ▀ (U+2580) which renders the top half in foreground color, and bottom half in background color
			fmt.Fprintf(w, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", top.R, top.G, top.B, bottomR, bottomG, bottomB)
		}
		fmt.Fprintf(w, "\x1b[0m\n")
	}
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
