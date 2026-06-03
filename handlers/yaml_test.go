package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestYAMLHandler(t *testing.T) {
	h := &YAMLHandler{}
	require.True(t, h.CanHandle("application/yaml", ".yaml"))
	require.True(t, h.CanHandle("text/yaml", ".yml"))
	require.False(t, h.CanHandle("text/plain", ".txt"))

	tests := []struct {
		name    string
		file    string
		opts    Options
		wantErr bool
	}{
		{"flat", "valid.yaml", Options{Flat: true}, false},
		{"plain", "valid.yaml", Options{Plain: true}, false},
		{"colored", "valid.yaml", Options{Theme: "monokai"}, false},
		{"invalid", "invalid.yaml", Options{}, true},
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
				require.Contains(t, err.Error(), "yaml:")
				return
			}

			require.NoError(t, err)

			testutil.AssertGolden(t, "fixtures/yaml_"+tt.name, out.Bytes())
		})
	}
}
