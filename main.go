package main

import (
	"asd/cmd"
	"asd/handlers"
	"asd/registry"
)

var version = "dev"

func main() {
	reg := registry.New(&handlers.BinaryHandler{})
	// Register other handlers here as they are implemented
	// e.g., reg.Register(&handlers.TextHandler{})

	cmd.SetRegistry(reg)
	cmd.Execute()
}
