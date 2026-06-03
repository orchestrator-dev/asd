package detect

import "net/http"

// MagicDetector uses magic bytes to detect MIME types.
type MagicDetector struct{}

func (d *MagicDetector) Detect(header []byte, filename string) (string, string) {
	if len(header) == 0 {
		return "", "none"
	}
	mime := http.DetectContentType(header)
	if mime == "application/octet-stream" {
		return "application/octet-stream", "low" // Let extension hint try
	}
	return mime, "high"
}
