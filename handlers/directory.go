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
		dirName = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(dirName)
	}
	fmt.Fprintf(w, "%s/\n", dirName)

	h.printTree(w, dir, "", opts, 0)
	return nil
}

func (h *DirectoryHandler) printTree(w io.Writer, dir string, prefix string, opts Options, depth int) {
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
				// Regular files get a nice cyan or pink depending on extension to match the screenshot
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
			h.printTree(w, filepath.Join(dir, entry.Name()), prefix+extension, opts, depth+1)
		}
	}
}
