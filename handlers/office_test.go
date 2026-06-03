package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestOfficeHandler(t *testing.T) {
	tests := []struct {
		name string
		file string
		mime string
		opts Options
	}{
		{"docx", "../testdata/fixtures/test.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", Options{}},
		{"xlsx", "../testdata/fixtures/test.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Options{}},
		{"pptx", "../testdata/fixtures/test.pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation", Options{}},
		{"odt", "../testdata/fixtures/test.odt", "application/vnd.oasis.opendocument.text", Options{}},
	}

	h := &OfficeHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, h.CanHandle(tt.mime, filepath.Ext(tt.file)))

			data, err := os.ReadFile(tt.file)
			assert.NoError(t, err)

			meta := FileMeta{Name: filepath.Base(tt.file), Size: int64(len(data))}
			var buf bytes.Buffer
			err = h.Render(&buf, bytes.NewReader(data), meta, tt.opts)
			assert.NoError(t, err)

			output := buf.String()
			assert.NotEmpty(t, output)

			testutil.AssertGolden(t, "fixtures/office_"+tt.name, []byte(output))
		})
	}
}
