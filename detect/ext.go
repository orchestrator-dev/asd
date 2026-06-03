package detect

import (
	"path/filepath"
	"strings"
)

// ExtDetector uses file extensions to suggest MIME types.
type ExtDetector struct{}

func (d *ExtDetector) Detect(header []byte, filename string) (string, string) {
	ext := strings.ToLower(filepath.Ext(filename))
	// Basic mapping, can be expanded later
	switch ext {
	case ".go", ".py", ".js", ".ts", ".html", ".css", ".md", ".json", ".yaml", ".yml", ".toml", ".csv", ".tsv":
		return "text/plain", "medium"
	default:
		return "", "none"
	}
}
