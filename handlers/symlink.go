package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"asd/detect"
)

type SymlinkHandler struct{}

func (h *SymlinkHandler) CanHandle(mime, ext string) bool {
	return mime == "inode/symlink"
}

func (h *SymlinkHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	target, err := os.Readlink(meta.Name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "-> %s\n", target)

	// Since we can't access globalReg, we can do a basic check and use text/binary handlers
	header := make([]byte, 512)
	n, _ := r.Read(header)

	// Reset reader? r is usually a file, we can't seek it if it's io.Reader without interface cast
	if seeker, ok := r.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}

	targetExt := filepath.Ext(target)
	mime := detect.Pipeline(header[:n], target)

	// Simple fallback dispatch
	var handler Handler
	if mime == "inode/directory" {
		handler = &DirectoryHandler{}
	} else if targetExt == ".json" || targetExt == ".jsonl" || mime == "application/json" {
		handler = &JSONHandler{}
	} else if mime == "text/plain" || mime == "text/markdown" || mime == "text/xml" || mime == "text/html" {
		handler = &TextHandler{}
	} else if mime == "application/octet-stream" {
		handler = &BinaryHandler{}
	} else {
		handler = &TextHandler{} // fallback
	}

	meta.Name = target
	return handler.Render(w, r, meta, opts)
}
