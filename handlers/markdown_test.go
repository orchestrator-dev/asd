package handlers_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/handlers"
	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestMarkdownHandler_Render(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		opts    handlers.Options
		golden  string
	}{
		{
			name:    "basic markdown",
			fixture: "test.md",
			opts:    handlers.Options{Flat: false, Width: 80},
			golden:  "markdown_basic",
		},
		{
			name:    "flat markdown",
			fixture: "test.md",
			opts:    handlers.Options{Flat: true},
			golden:  "markdown_flat",
		},
		{
			name:    "plain markdown",
			fixture: "test.md",
			opts:    handlers.Options{Flat: false, Plain: true, Width: 80},
			golden:  "markdown_plain",
		},
	}

	h := &handlers.MarkdownHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "testdata", "fixtures", tt.fixture)
			f, err := os.Open(path)
			require.NoError(t, err)
			defer f.Close()

			stat, err := f.Stat()
			require.NoError(t, err)

			meta := handlers.FileMeta{
				Name: tt.fixture,
				Size: stat.Size(),
			}

			var buf bytes.Buffer
			err = h.Render(&buf, f, meta, tt.opts)
			require.NoError(t, err)

			testutil.AssertGolden(t, tt.golden, buf.Bytes())
		})
	}
}

func TestMarkdownHandler_CanHandle(t *testing.T) {
	h := &handlers.MarkdownHandler{}
	require.True(t, h.CanHandle("text/markdown", ".md"))
	require.False(t, h.CanHandle("text/plain", ".txt"))
}
