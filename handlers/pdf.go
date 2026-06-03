package handlers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// PDFHandler handles PDF files.
type PDFHandler struct{}

func (h *PDFHandler) CanHandle(mime, ext string) bool {
	return strings.ToLower(ext) == ".pdf" || mime == "application/pdf"
}

func (h *PDFHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	rs, ok := r.(io.ReadSeeker)
	if !ok {
		f, err := os.CreateTemp("", "asd-pdf-*")
		if err != nil {
			return err
		}
		defer os.Remove(f.Name())
		defer f.Close()
		if _, err := io.Copy(f, r); err != nil {
			return err
		}
		if _, err := f.Seek(0, 0); err != nil {
			return err
		}
		rs = f
	}

	conf := model.NewDefaultConfiguration()
	ctx, err := api.ReadContext(rs, conf)
	if err != nil {
		return err
	}

	title := ctx.Title
	if title == "" {
		title = "(none)"
	}
	author := ctx.Author
	if author == "" {
		author = "(none)"
	}

	fmt.Fprintf(w, "PDF Version: %s\n", ctx.HeaderVersion)
	fmt.Fprintf(w, "Title: %s\n", title)
	fmt.Fprintf(w, "Author: %s\n", author)
	fmt.Fprintf(w, "Pages: %d\n\n", ctx.PageCount)

	tempDir, err := os.MkdirTemp("", "asd-pdfcpu-extract-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if _, err := rs.Seek(0, 0); err != nil {
		return err
	}

	err = api.ExtractContent(rs, tempDir, "out", []string{"1"}, conf)
	if err != nil {
		fmt.Fprintf(w, "Failed to extract content: %v\n", err)
	} else {
		contentPath := filepath.Join(tempDir, "out_Content_page_1.txt")
		cf, err := os.Open(contentPath)
		if err == nil {
			defer cf.Close()
			scanner := bufio.NewScanner(cf)
			lineCount := 0
			for scanner.Scan() {
				if lineCount >= 100 {
					break
				}
				fmt.Fprintln(w, scanner.Text())
				lineCount++
			}
		}
	}

	fmt.Fprintf(w, "\n... (%d pages total)\n", ctx.PageCount)
	return nil
}
