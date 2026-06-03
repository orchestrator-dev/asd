package handlers

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymlinkHandler(t *testing.T) {
	h := &SymlinkHandler{}

	assert.True(t, h.CanHandle("inode/symlink", ""))

	meta := FileMeta{Name: "../testdata/fixtures/test_symlink.txt"}
	f, err := os.Open(meta.Name) // Open target
	assert.NoError(t, err)
	defer f.Close()

	var buf bytes.Buffer
	err = h.Render(&buf, f, meta, Options{})
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "-> test.txt")
}
