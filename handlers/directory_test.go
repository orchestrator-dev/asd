package handlers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryHandler(t *testing.T) {
	h := &DirectoryHandler{}

	assert.True(t, h.CanHandle("inode/directory", ""))

	meta := FileMeta{Name: "../testdata"}
	var buf bytes.Buffer
	err := h.Render(&buf, nil, meta, Options{NoColor: true})
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "fixtures")
	assert.Contains(t, output, "binary_unknown.golden")
}
