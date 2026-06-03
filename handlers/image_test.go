package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

type imageMockRunner struct {
	available func(string) bool
	run       func(name string, args ...string) ([]byte, error)
}

func (m *imageMockRunner) Available(name string) bool { return m.available != nil && m.available(name) }
func (m *imageMockRunner) Run(name string, args ...string) ([]byte, error) {
	if m.run != nil {
		return m.run(name, args...)
	}
	return nil, nil
}

func TestImageHandler(t *testing.T) {
	runner := &imageMockRunner{
		available: func(n string) bool { return n == "chafa" },
		run: func(name string, args ...string) ([]byte, error) {
			return []byte("chafa output for " + args[0]), nil
		},
	}
	h := &ImageHandler{Runner: runner}
	require.True(t, h.CanHandle("image/png", ".png"))
	require.True(t, h.CanHandle("image/svg+xml", ".svg"))

	tests := []struct {
		name string
		file string
		opts Options
	}{
		{"gif", "test.gif", Options{Plain: true}},
		{"svg", "test.svg", Options{Plain: true}},
		{"flat", "test.gif", Options{Flat: true}},
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

			testutil.AssertGolden(t, "fixtures/image_"+tt.name, out.Bytes())
		})
	}
}
