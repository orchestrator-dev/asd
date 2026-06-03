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

func TestBinaryHandler_Render(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		opts    handlers.Options
		golden  string
	}{
		{
			name:    "unknown binary",
			fixture: "unknown.bin",
			opts:    handlers.Options{Flat: false},
			golden:  "binary_unknown",
		},
		{
			name:    "unknown binary flat",
			fixture: "unknown.bin",
			opts:    handlers.Options{Flat: true},
			golden:  "binary_unknown_flat",
		},
	}

	h := &handlers.BinaryHandler{}

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
