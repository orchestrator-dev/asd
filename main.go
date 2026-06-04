package main

import (
	"asd/cmd"
	"asd/exec"
	"asd/handlers"
	"asd/registry"
)

var version = "dev"

func main() {
	reg := registry.New(&handlers.BinaryHandler{})
	reg.Register(&handlers.JSONHandler{})
	reg.Register(&handlers.YAMLHandler{})
	reg.Register(&handlers.TOMLHandler{})
	reg.Register(&handlers.CSVHandler{})
	reg.Register(&handlers.MarkdownHandler{})
	reg.Register(&handlers.CertHandler{})
	// TextHandler must come AFTER specific text handlers (e.g. CSV, Markdown, Cert)
	reg.Register(&handlers.TextHandler{})
	reg.Register(&handlers.ImageHandler{Runner: &exec.OSExecRunner{}})
	reg.Register(&handlers.AudioHandler{Runner: &exec.OSExecRunner{}})
	reg.Register(&handlers.VideoHandler{Runner: &exec.OSExecRunner{}})
	reg.Register(&handlers.ArchiveHandler{Runner: &exec.OSExecRunner{}})
	reg.Register(&handlers.OfficeHandler{})
	reg.Register(&handlers.CertHandler{})
	reg.Register(&handlers.DirectoryHandler{})

	cmd.SetRegistry(reg)
	cmd.Execute()
}
