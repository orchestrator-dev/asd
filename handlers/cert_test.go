package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"asd/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestCertHandler(t *testing.T) {
	tests := []struct {
		name string
		file string
		mime string
		opts Options
	}{
		{"pem", "../testdata/fixtures/test.pem", "application/x-x509-ca-cert", Options{}},
		{"key", "../testdata/fixtures/test.key", "application/pkcs8", Options{}},
		{"p12", "../testdata/fixtures/test.p12", "application/x-pkcs12", Options{}},
		{"pub", "../testdata/fixtures/test.pub", "application/ssh-public-key", Options{}},
	}

	h := &CertHandler{}

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

			// The golden file for test.pem / key / p12 etc might have different content generated.
			// Let's use testutil.AssertGolden? But the ssh-keygen output changes fingerprint or random serial?
			// Since I generated them locally and they will be committed, they are static!
			testutil.AssertGolden(t, "fixtures/cert_"+tt.name, []byte(output))
		})
	}
}
