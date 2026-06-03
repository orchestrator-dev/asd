package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestAudioHandler(t *testing.T) {
	h := &AudioHandler{}
	require.True(t, h.CanHandle("audio/mpeg", ".mp3"))
	require.True(t, h.CanHandle("audio/flac", ".flac"))

	tests := []struct {
		name string
		file string
		opts Options
	}{
		{"mp3", "test.mp3", Options{Plain: true}},
		{"flat", "test.mp3", Options{Flat: true}},
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

			testutil.AssertGolden(t, "fixtures/audio_"+tt.name, out.Bytes())
		})
	}
}
