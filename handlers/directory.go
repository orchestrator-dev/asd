package handlers

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

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

	// Re-add . and .. if we want ls -la behavior?
	// The prompt said "Styled ls -la tree"

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
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if u, err := user.LookupId(fmt.Sprintf("%d", stat.Uid)); err == nil {
			uname = u.Username
		} else {
			uname = fmt.Sprintf("%d", stat.Uid)
		}
		if g, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); err == nil {
			gname = g.Name
		} else {
			gname = fmt.Sprintf("%d", stat.Gid)
		}
	}

	displayName := name
	if !opts.NoColor {
		if mode.IsDir() {
			displayName = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render(name) // Blue
		} else if mode&os.ModeSymlink != 0 {
			displayName = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Render(name) // Cyan
		} else if mode&0111 != 0 {
			displayName = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(name) // Green
		}
	}

	fmt.Fprintf(w, "%s%s %-8s %-8s %8d %s %s\n", prefix, mode.String(), uname, gname, size, modTime, displayName)
}
