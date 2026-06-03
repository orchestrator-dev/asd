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

func TestXMLHandler_Render(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		opts    handlers.Options
		golden  string
	}{
		{
			name:    "basic xml",
			fixture: "test.xml",
			opts:    handlers.Options{Flat: false, NoColor: true}, // testing with NoColor so we don't depend on terminal color sequences in golden
			golden:  "xml_basic",
		},
		{
			name:    "flat xml",
			fixture: "test.xml",
			opts:    handlers.Options{Flat: true},
			golden:  "xml_flat",
		},
		{
			name:    "basic html",
			fixture: "test.html",
			opts:    handlers.Options{Flat: false, NoColor: true},
			golden:  "html_basic",
		},
	}

	h := &handlers.XMLHandler{}

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

func TestXMLHandler_CanHandle(t *testing.T) {
	h := &handlers.XMLHandler{}
	require.True(t, h.CanHandle("text/xml", ".xml"))
	require.True(t, h.CanHandle("text/html", ".html"))
	require.False(t, h.CanHandle("text/plain", ".txt"))
}
