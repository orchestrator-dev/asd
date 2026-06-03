package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

type videoMockRunner struct {
	available func(string) bool
	run       func(name string, args ...string) ([]byte, error)
}

func (m *videoMockRunner) Available(name string) bool { return m.available != nil && m.available(name) }
func (m *videoMockRunner) Run(name string, args ...string) ([]byte, error) {
	if m.run != nil {
		return m.run(name, args...)
	}
	return nil, nil
}

func TestVideoHandler(t *testing.T) {
	runner := &videoMockRunner{
		available: func(n string) bool { return n == "ffprobe" },
		run: func(name string, args ...string) ([]byte, error) {
			jsonOut := `{"format":{"duration":"120.5","bit_rate":"500000"},"streams":[{"codec_type":"video","codec_name":"h264","width":1920,"height":1080,"r_frame_rate":"30000/1001"}]}`
			return []byte(jsonOut), nil
		},
	}
	h := &VideoHandler{Runner: runner}
	require.True(t, h.CanHandle("video/mp4", ".mp4"))
	require.True(t, h.CanHandle("video/x-matroska", ".mkv"))

	tests := []struct {
		name string
		file string
		opts Options
	}{
		{"mp4", "test.mp4", Options{Plain: true}},
		{"flat", "test.mp4", Options{Flat: true}},
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

			testutil.AssertGolden(t, "fixtures/video_"+tt.name, out.Bytes())
		})
	}
}
