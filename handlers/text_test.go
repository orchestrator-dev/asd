package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestTextHandler(t *testing.T) {
	h := &TextHandler{}
	require.True(t, h.CanHandle("text/plain", ".txt"))
	require.True(t, h.CanHandle("text/html", ".html"))
	require.False(t, h.CanHandle("application/json", ".json"))

	tests := []struct {
		name string
		file string
		opts Options
	}{
		{"flat", "test.txt", Options{Flat: true}},
		{"plain", "test.txt", Options{Plain: true}},
		{"colored_go", "test.go", Options{Theme: "monokai"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "testdata", "fixtures", tt.file)
			data, err := os.ReadFile(path)
			require.NoError(t, err)

			var out bytes.Buffer
			meta := FileMeta{Name: tt.file, Size: int64(len(data))}
			err = h.Render(&out, bytes.NewReader(data), meta, tt.opts)
			require.NoError(t, err)

			testutil.AssertGolden(t, "fixtures/text_"+tt.name, out.Bytes())
		})
	}
}
