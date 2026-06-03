package handlers

import (
	"encoding/hex"
	"fmt"
	"io"
)

// BinaryHandler provides a hex dump fallback for unknown/binary files.
type BinaryHandler struct{}

func (h *BinaryHandler) CanHandle(mime, ext string) bool {
	return true // Fallback
}

func (h *BinaryHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		// Just copy raw data if flat mode is enabled
		_, err := io.Copy(w, r)
		return err
	}

	fmt.Fprintf(w, "File: %s\nSize: %d bytes\nMIME: application/octet-stream (fallback)\n\n", meta.Name, meta.Size)

	buf := make([]byte, 512)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return err
	}

	dumper := hex.Dumper(w)
	_, _ = dumper.Write(buf[:n])
	dumper.Close()

	if meta.Size > 512 {
		fmt.Fprintf(w, "\n… (%d bytes total)\n", meta.Size)
	}

	return nil
}
