package handlers

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	if opts.Flat {
		// flat mode: just list files
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, e := range entries {
			fmt.Fprintln(w, e.Name())
		}
		return nil
	}

	dirName := filepath.Base(dir)
	if !opts.NoColor && !opts.Plain {
		dirName = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Render("📁 " + dirName)
	} else {
		dirName = "📁 " + dirName
	}
	fmt.Fprintf(w, "%s/\n", dirName)

	stats := &dirStats{}
	h.printTree(w, dir, "", opts, 0, stats)

	// Dashboard Summary
	if !opts.Clean {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Repeat("━", 40))
		fmt.Fprintf(w, "📊 Dashboard: %s\n", dir)
		fmt.Fprintf(w, "   Files: %d | Directories: %d\n", stats.files, stats.dirs)

		// Try Git
		if out, err := exec.Command("git", "-C", dir, "branch", "--show-current").Output(); err == nil && len(out) > 0 {
			branch := strings.TrimSpace(string(out))
			fmt.Fprintf(w, "   Git Branch: %s\n", branch)
			if status, err := exec.Command("git", "-C", dir, "status", "-s").Output(); err == nil {
				lines := strings.Split(strings.TrimSpace(string(status)), "\n")
				modified := 0
				for _, l := range lines {
					if len(l) > 0 {
						modified++
					}
				}
				fmt.Fprintf(w, "   Git Status: %d modified/untracked files\n", modified)
			}
		}
		fmt.Fprintln(w, strings.Repeat("━", 40))
	}

	return nil
}

type dirStats struct {
	files int
	dirs  int
}

func (h *DirectoryHandler) printTree(w io.Writer, dir string, prefix string, opts Options, depth int, stats *dirStats) {
	if depth > 4 {
		return // max depth 4 to prevent crazy output
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.IsDir() {
			stats.dirs++
		} else {
			stats.files++
		}

		isLast := i == len(entries)-1
		pointer := "├── "
		if isLast {
			pointer = "└── "
		}

		name := entry.Name()
		mode := info.Mode()

		if !opts.NoColor && !opts.Plain {
			if mode.IsDir() {
				name = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(name)
			} else if mode&os.ModeSymlink != 0 {
				name = lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Render(name)
			} else if mode&0111 != 0 {
				name = lipgloss.NewStyle().Foreground(lipgloss.Color("40")).Render(name)
			} else {
				ext := filepath.Ext(name)
				switch ext {
				case ".json", ".yaml", ".yml", ".toml", ".csv":
					name = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Render(name)
				case ".log", ".dat", ".txt":
					name = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Render(name)
				case ".go", ".js", ".py", ".rs":
					name = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(name)
				default:
					name = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(name)
				}
			}
		}

		suffix := ""
		if mode.IsDir() {
			suffix = "/"
		} else if mode&os.ModeSymlink != 0 {
			suffix = "@"
		} else if mode&0111 != 0 {
			suffix = "*"
		}

		fmt.Fprintf(w, "%s%s%s%s\n", prefix, pointer, name, suffix)

		if mode.IsDir() {
			extension := "│   "
			if isLast {
				extension = "    "
			}
			h.printTree(w, filepath.Join(dir, entry.Name()), prefix+extension, opts, depth+1, stats)
		}
	}
}
