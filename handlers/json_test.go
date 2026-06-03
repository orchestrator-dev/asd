package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestJSONHandler(t *testing.T) {
	h := &JSONHandler{}
	require.True(t, h.CanHandle("application/json", ".json"))
	require.True(t, h.CanHandle("", ".jsonl"))
	require.False(t, h.CanHandle("text/plain", ".txt"))

	tests := []struct {
		name    string
		file    string
		opts    Options
		wantErr bool
	}{
		{"flat", "valid.json", Options{Flat: true}, false},
		{"plain", "valid.json", Options{Plain: true}, false},
		{"colored", "valid.json", Options{Theme: "monokai"}, false},
		{"jsonl", "valid.jsonl", Options{Plain: true}, false},
		{"invalid", "invalid.json", Options{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "testdata", "fixtures", tt.file)
			data, err := os.ReadFile(path)
			require.NoError(t, err)

			var out bytes.Buffer
			meta := FileMeta{Name: tt.file, Size: int64(len(data))}
			err = h.Render(&out, bytes.NewReader(data), meta, tt.opts)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "json:")
				return
			}

			require.NoError(t, err)

			testutil.AssertGolden(t, "fixtures/json_"+tt.name, out.Bytes())
		})
	}
}
