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

func TestCSVHandler_Render(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		opts    handlers.Options
		golden  string
	}{
		{
			name:    "basic csv",
			fixture: "test.csv",
			opts:    handlers.Options{Flat: false},
			golden:  "csv_basic",
		},
		{
			name:    "flat csv",
			fixture: "test.csv",
			opts:    handlers.Options{Flat: true},
			golden:  "csv_flat",
		},
	}

	h := &handlers.CSVHandler{}

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

func TestCSVHandler_CanHandle(t *testing.T) {
	h := &handlers.CSVHandler{}
	require.True(t, h.CanHandle("text/csv", ".csv"))
	require.True(t, h.CanHandle("", ".tsv"))
	require.False(t, h.CanHandle("text/plain", ".txt"))
}
