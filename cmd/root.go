package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"asd/detect"
	"asd/handlers"
	"asd/registry"
	"asd/render"

	"github.com/spf13/cobra"
)

var (
	opts     handlers.Options
	globalReg *registry.Registry
)

func SetRegistry(r *registry.Registry) {
	globalReg = r
}

var rootCmd = &cobra.Command{
	Use:   "asd [file...]",
	Short: "A smart universal file viewer",
	RunE: func(cmd *cobra.Command, args []string) error {
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

	if filename == "-" {
		f = os.Stdin
		filename = "stdin"
	} else {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	header := make([]byte, 512)
	n, _ := f.Read(header)
	f.Seek(0, 0) // Reset read pointer if possible, stdin might fail but that's ok for now

	mime := detect.Pipeline(header[:n], filename)
	ext := filepath.Ext(filename)

	handler := globalReg.Dispatch(mime, ext)

	meta := handlers.FileMeta{
		Name: filename,
		Size: stat.Size(),
	}

	w := render.NewWriter(os.Stdout, opts)

	if opts.Flat {
		_, err = io.Copy(w, f)
		return err
	}

	return handler.Render(w, f, meta, opts)
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
	rootCmd.Flags().StringVar(&opts.Theme, "theme", "", "chroma highlight theme (default: auto)")
}
