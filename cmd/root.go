package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"asd/detect"
	"asd/handlers"
	"asd/registry"
	"asd/render"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	opts      handlers.Options
	globalReg *registry.Registry
)

func SetRegistry(r *registry.Registry) {
	globalReg = r
}

var rootCmd = &cobra.Command{
	Use:   "asd [file...]",
	Short: "A smart universal file viewer",
	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.Theme == "" || opts.Theme == "auto" {
			opts.Theme = "dracula"
		}

		opts.Width = 80 // fallback
		if term.IsTerminal(int(os.Stdout.Fd())) {
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				opts.Width = w
			}
		} else {
			opts.NoColor = true
		}

		if opts.Diff {
			if len(args) != 2 {
				return fmt.Errorf("--diff requires exactly two files")
			}
			return processDiff(args[0], args[1])
		}

		if len(args) == 0 {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return cmd.Help()
			}
			return processFile("-")
		}

		for i, arg := range args {
			if len(args) > 1 && !opts.Flat {
				if i > 0 {
					fmt.Println()
				}
				fmt.Printf("── %s ──\n", arg)
			}
			if err := processFile(arg); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
		return nil
	},
}

func processFile(filename string) error {
	var f *os.File
	var err error

	var isHTTP = strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://")

	if filename == "-" {
		f = os.Stdin
		filename = "stdin"
	} else if !isHTTP {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	var fileSize int64
	if f != nil {
		stat, err := f.Stat()
		if err == nil {
			fileSize = stat.Size()
		}
	}

	header := make([]byte, 512)
	n := 0
	if f != nil {
		n, _ = f.Read(header)
	}

	var r io.Reader = f
	var filenameForMime = filename

	if f != nil {
		if _, err := f.Seek(0, 0); err != nil {
			r = io.MultiReader(bytes.NewReader(header[:n]), f)
		}
	} else if isHTTP {
		// HTTP Fetching
		resp, err := http.Get(filename)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Read a snippet for mime detection
		header = make([]byte, 512)
		n, _ = io.ReadFull(resp.Body, header)
		r = io.MultiReader(bytes.NewReader(header[:n]), resp.Body)

		// Try to use Content-Type from headers
		contentType := resp.Header.Get("Content-Type")
		if contentType != "" {
			parts := strings.Split(contentType, ";")
			mimeType := strings.TrimSpace(parts[0])
			subtypeParts := strings.Split(mimeType, "/")
			if len(subtypeParts) == 2 {
				filenameForMime = "download." + subtypeParts[1]
			}
		}
	}

	var tailCmd *exec.Cmd
	if opts.Follow && filename != "stdin" && filename != "-" {
		tailCmd = exec.Command("tail", "-f", "-n", "+1", filename)
		if stdout, err := tailCmd.StdoutPipe(); err == nil {
			if err := tailCmd.Start(); err == nil {
				r = stdout
				opts.NoPager = true // Disallow paging in follow mode
				defer tailCmd.Process.Kill()
			}
		}
	}

	mime := detect.Pipeline(header[:n], filenameForMime)
	ext := filepath.Ext(filenameForMime)

	// Heuristic: If it's in a log directory or named log, treat it as a log
	baseName := strings.ToLower(filepath.Base(filenameForMime))
	if strings.Contains(baseName, "log") || strings.Contains(filenameForMime, "/var/log/") {
		ext = ".log"
	}

	handler := globalReg.Dispatch(mime, ext)

	meta := handlers.FileMeta{
		Name: filenameForMime,
		Size: fileSize,
	}

	w := render.NewWriter(os.Stdout, opts)
	defer w.Close()

	if opts.Explain && mime != "inode/directory" {
		if data, err := io.ReadAll(r); err == nil {
			explainFile(w, data, filenameForMime, opts.Theme)
			r = bytes.NewReader(data) // Restore reader
		}
	}

	if opts.Flat {
		_, err = io.Copy(w, r)
		return err
	}

	return handler.Render(w, r, meta, opts)
}

func processDiff(file1, file2 string) error {
	w := render.NewWriter(os.Stdout, opts)
	defer w.Close()

	var cmd *exec.Cmd
	// Check if delta is available
	if err := exec.Command("delta", "--version").Run(); err == nil {
		cmd = exec.Command("delta", "-s", file1, file2)
	} else if err := exec.Command("git", "diff", "--no-index").Run(); err == nil || err.Error() == "exit status 1" {
		cmd = exec.Command("git", "diff", "--no-index", "--color=always", file1, file2)
	} else {
		cmd = exec.Command("diff", "-y", "--color=always", file1, file2)
	}

	out, _ := cmd.CombinedOutput()
	fmt.Fprint(w, string(out))
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&opts.Flat, "flat", "f", false, "bypass smart rendering; behave like cat")
	rootCmd.Flags().BoolVarP(&opts.Lines, "lines", "n", false, "show line numbers (text/source only)")
	rootCmd.Flags().BoolVarP(&opts.Plain, "plain", "p", false, "disable all color and styling")
	rootCmd.Flags().BoolVar(&opts.Clean, "clean", false, "disable UI decorations (headers, line numbers)")
	rootCmd.Flags().BoolVarP(&opts.Follow, "follow", "F", false, "tail/follow mode for continuous reading")
	rootCmd.Flags().BoolVar(&opts.Diff, "diff", false, "render a side-by-side diff of two files")
	rootCmd.Flags().BoolVar(&opts.NoPager, "no-pager", false, "disable auto-paging")
	rootCmd.Flags().BoolVar(&opts.Explain, "explain", false, "use local AI (Ollama) to explain the file")
	rootCmd.Flags().StringVar(&opts.Theme, "theme", "", "chroma highlight theme (default: auto)")
}
