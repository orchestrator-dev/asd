package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestTOMLHandler(t *testing.T) {
	h := &TOMLHandler{}
	require.True(t, h.CanHandle("application/toml", ".toml"))
	require.False(t, h.CanHandle("text/plain", ".txt"))

	tests := []struct {
		name    string
		file    string
		opts    Options
		wantErr bool
	}{
		{"flat", "valid.toml", Options{Flat: true}, false},
		{"plain", "valid.toml", Options{Plain: true}, false},
		{"colored", "valid.toml", Options{Theme: "monokai"}, false},
		{"invalid", "invalid.toml", Options{}, true},
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
				require.Contains(t, err.Error(), "toml:")
				return
			}

			require.NoError(t, err)

			testutil.AssertGolden(t, "fixtures/toml_"+tt.name, out.Bytes())
		})
	}
}
