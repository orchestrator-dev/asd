package handlers

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"asd/exec"

	"github.com/ulikunitz/xz"
)

type ArchiveHandler struct {
	Runner exec.Runner
}

func (h *ArchiveHandler) CanHandle(mime, ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".zip", ".tar", ".gz", ".tgz", ".bz2", ".xz", ".7z", ".rar":
		return true
	}
	return false
}

func (h *ArchiveHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	ext := strings.ToLower(filepath.Ext(meta.Name))

	nameLower := strings.ToLower(meta.Name)
	if strings.HasSuffix(nameLower, ".tar.gz") || ext == ".tgz" {
		return h.handleTarGz(w, r)
	}
	if strings.HasSuffix(nameLower, ".tar.bz2") {
		return h.handleTarBz2(w, r)
	}
	if strings.HasSuffix(nameLower, ".tar.xz") {
		return h.handleTarXz(w, r)
	}

	switch ext {
	case ".zip":
		return h.handleZip(w, r, meta)
	case ".tar":
		return h.handleTar(w, r)
	case ".gz":
		return h.handleGz(w, r, meta)
	case ".bz2":
		return fmt.Errorf("plain .bz2 listing not supported")
	case ".xz":
		return fmt.Errorf("plain .xz listing not supported")
	case ".7z":
		return h.handle7z(w, meta)
	case ".rar":
		return h.handleRar(w, meta)
	}

	return fmt.Errorf("unsupported archive: %s", ext)
}

func (h *ArchiveHandler) handleZip(w io.Writer, r io.Reader, meta FileMeta) error {
	ra, ok := r.(io.ReaderAt)
	if !ok {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		ra = bytes.NewReader(data)
	}

	zr, err := zip.NewReader(ra, meta.Size)
	if err != nil {
		return err
	}

	var total uint64
	var count int

	for _, f := range zr.File {
		count++
		total += f.UncompressedSize64
		fmt.Fprintf(w, "%s %12d %s %s\n", f.Mode(), f.UncompressedSize64, f.Modified.Format("2006-01-02 15:04:05"), f.Name)
	}
	fmt.Fprintf(w, "\n%d files, total uncompressed size: %d\n", count, total)
	return nil
}

func (h *ArchiveHandler) handleTarReader(w io.Writer, tr *tar.Reader) error {
	var total uint64
	var count int

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		count++
		total += uint64(hdr.Size)
		modTime := hdr.ModTime.Format("2006-01-02 15:04:05")
		mode := os.FileMode(hdr.Mode)
		if hdr.Typeflag == tar.TypeDir {
			mode |= os.ModeDir
		}
		fmt.Fprintf(w, "%s %12d %s %s\n", mode, hdr.Size, modTime, hdr.Name)
	}
	fmt.Fprintf(w, "\n%d files, total uncompressed size: %d\n", count, total)
	return nil
}

func (h *ArchiveHandler) handleTar(w io.Writer, r io.Reader) error {
	return h.handleTarReader(w, tar.NewReader(r))
}

func (h *ArchiveHandler) handleTarGz(w io.Writer, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()
	return h.handleTarReader(w, tar.NewReader(gr))
}

func (h *ArchiveHandler) handleTarBz2(w io.Writer, r io.Reader) error {
	br := bzip2.NewReader(r)
	return h.handleTarReader(w, tar.NewReader(br))
}

func (h *ArchiveHandler) handleTarXz(w io.Writer, r io.Reader) error {
	xr, err := xz.NewReader(r)
	if err != nil {
		return err
	}
	return h.handleTarReader(w, tar.NewReader(xr))
}

func (h *ArchiveHandler) handleGz(w io.Writer, r io.Reader, meta FileMeta) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	var uncompressed int64
	buf := make([]byte, 32*1024)
	for {
		n, err := gr.Read(buf)
		uncompressed += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "Compressed size: %d bytes\nUncompressed size: %d bytes\n", meta.Size, uncompressed)
	return nil
}

func (h *ArchiveHandler) handle7z(w io.Writer, meta FileMeta) error {
	if h.Runner == nil {
		return fmt.Errorf("no runner available for 7z")
	}
	out, err := h.Runner.Run("7z", "l", meta.Name)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

func (h *ArchiveHandler) handleRar(w io.Writer, meta FileMeta) error {
	if h.Runner == nil {
		return fmt.Errorf("no runner available for rar")
	}
	out, err := h.Runner.Run("unrar", "l", meta.Name)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
