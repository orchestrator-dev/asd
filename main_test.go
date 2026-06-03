package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "asd")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	err := cmd.Run()
	require.NoError(t, err, "failed to build binary")
	return bin
}

func TestCLI_FlatMode(t *testing.T) {
	bin := buildBinary(t)
	fixture := filepath.Join("testdata", "fixtures", "unknown.bin")

	cmd := exec.Command(bin, "-f", fixture)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	require.NoError(t, err)

	catCmd := exec.Command("cat", fixture)
	var catStdout bytes.Buffer
	catCmd.Stdout = &catStdout
	err = catCmd.Run()
	require.NoError(t, err)

	require.Equal(t, catStdout.Bytes(), stdout.Bytes())
}
