package detect_test

import (
	"os"
	"path/filepath"
	"testing"

	"asd/detect"
	_ "asd/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		expected string
	}{
		{
			name:     "unknown binary",
			fixture:  "unknown.bin",
			expected: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "testdata", "fixtures", tt.fixture)
			b, err := os.ReadFile(path)
			require.NoError(t, err)

			header := b
			if len(header) > 512 {
				header = header[:512]
			}

			mime := detect.Pipeline(header, tt.fixture)
			require.Equal(t, tt.expected, mime)
		})
	}
}
