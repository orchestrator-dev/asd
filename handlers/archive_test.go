package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"asd/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
)

type mockArchiveRunner struct {
	output []byte
	err    error
}

func (m *mockArchiveRunner) Run(name string, args ...string) ([]byte, error) {
	return m.output, m.err
}

func (m *mockArchiveRunner) Available(name string) bool {
	return true
}

func createTestZip(t *testing.T) string {
	path := filepath.Join("..", "testdata", "fixtures", "sample.zip")
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	w := zip.NewWriter(f)
	fw, err := w.CreateHeader(&zip.FileHeader{
		Name:     "test.txt",
		Method:   zip.Store,
		Modified: time.Date(2026, 6, 3, 23, 15, 47, 0, time.UTC),
	})
	require.NoError(t, err)
	_, err = fw.Write([]byte("hello zip"))
	require.NoError(t, err)
	w.Close()
	return path
}

func createTestTar(t *testing.T) string {
	path := filepath.Join("..", "testdata", "fixtures", "sample.tar")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	w := tar.NewWriter(f)
	err = w.WriteHeader(&tar.Header{
		Name:    "test.txt",
		Size:    9,
		Mode:    0644,
		ModTime: time.Date(2026, 6, 3, 23, 15, 47, 0, time.UTC),
	})
	require.NoError(t, err)
	_, err = w.Write([]byte("hello tar"))
	require.NoError(t, err)
	w.Close()
	return path
}

func createTestTarGz(t *testing.T) string {
	path := filepath.Join("..", "testdata", "fixtures", "sample.tar.gz")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	gw := gzip.NewWriter(f)
	w := tar.NewWriter(gw)
	err = w.WriteHeader(&tar.Header{
		Name:    "test.txt",
		Size:    9,
		Mode:    0644,
		ModTime: time.Date(2026, 6, 3, 23, 15, 47, 0, time.UTC),
	})
	require.NoError(t, err)
	_, err = w.Write([]byte("hello tgz"))
	require.NoError(t, err)
	w.Close()
	gw.Close()
	return path
}

func createTestGz(t *testing.T) string {
	path := filepath.Join("..", "testdata", "fixtures", "sample.gz")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()

	gw := gzip.NewWriter(f)
	_, err = gw.Write([]byte("hello plain gz file"))
	require.NoError(t, err)
	gw.Close()
	return path
}

func TestArchiveHandler_Render(t *testing.T) {
	createTestZip(t)
	createTestTar(t)
	createTestTarGz(t)
	createTestGz(t)

	tests := []struct {
		name     string
		fileName string
		runner   *mockArchiveRunner
		ext      string
	}{
		{"zip", "sample.zip", nil, ".zip"},
		{"tar", "sample.tar", nil, ".tar"},
		{"tgz", "sample.tar.gz", nil, ".tgz"},
		{"gz", "sample.gz", nil, ".gz"},
		{"7z", "sample.7z", &mockArchiveRunner{output: []byte("7z listing mock")}, ".7z"},
		{"rar", "sample.rar", &mockArchiveRunner{output: []byte("rar listing mock")}, ".rar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ArchiveHandler{Runner: tt.runner}
			assert.True(t, h.CanHandle("", tt.ext))

			var r io.Reader = bytes.NewReader([]byte{})
			var size int64 = 0
			if tt.runner == nil {
				path := filepath.Join("..", "testdata", "fixtures", tt.fileName)
				f, err := os.Open(path)
				require.NoError(t, err)
				defer f.Close()
				r = f
				stat, _ := f.Stat()
				size = stat.Size()
			}

			meta := FileMeta{Name: tt.fileName, Size: size}
			var buf bytes.Buffer
			err := h.Render(&buf, r, meta, Options{})
			require.NoError(t, err)

			testutil.AssertGolden(t, filepath.Join("fixtures", "archive_"+tt.name), buf.Bytes())
		})
	}
}
