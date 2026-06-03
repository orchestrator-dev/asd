package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPDFHandler_Render(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{"pdf", "sample.pdf"},
	}

	h := &PDFHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, h.CanHandle("application/pdf", ".pdf"))

			path := filepath.Join("..", "testdata", "fixtures", tt.fileName)
			f, err := os.Open(path)
			require.NoError(t, err)
			defer f.Close()

			stat, err := f.Stat()
			require.NoError(t, err)

			meta := FileMeta{Name: tt.fileName, Size: stat.Size()}
			var buf bytes.Buffer
			err = h.Render(&buf, f, meta, Options{})
			require.NoError(t, err)

			testutil.AssertGolden(t, filepath.Join("fixtures", "pdf_"+tt.name), buf.Bytes())
		})
	}
}
