package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

type DirectoryHandler struct{}

func (h *DirectoryHandler) CanHandle(mime, ext string) bool {
	return mime == "inode/directory"
}

func (h *DirectoryHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	dir := meta.Name
	if dir == "" || dir == "stdin" || dir == "-" {
		return fmt.Errorf("invalid directory path")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	dirInfo, err := os.Stat(dir)
	if err == nil {
		h.printEntry(w, dirInfo, ".", "", false, opts)
	}
	parentInfo, err := os.Stat(filepath.Dir(dir))
	if err == nil {
		h.printEntry(w, parentInfo, "..", "", false, opts)
	}

	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		isLast := i == len(entries)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}
		h.printEntry(w, info, entry.Name(), prefix, isLast, opts)
	}

	return nil
}

func (h *DirectoryHandler) printEntry(w io.Writer, info os.FileInfo, name, prefix string, isLast bool, opts Options) {
	mode := info.Mode()
	size := info.Size()
	modTime := info.ModTime().Format("Jan 02 15:04")

	var uname, gname string
	uname = "-"
	gname = "-"

	displayName := name
	if !opts.NoColor {
		if mode.IsDir() {
			displayName = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render(name)
		} else if mode&os.ModeSymlink != 0 {
			displayName = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render(name)
		} else if mode&0111 != 0 {
			displayName = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(name)
		}
	}

	fmt.Fprintf(w, "%s%s %-8s %-8s %8d %s %s\n", prefix, mode.String(), uname, gname, size, modTime, displayName)
}
