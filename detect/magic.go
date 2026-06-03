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
	// some custom overrides
	if len(header) >= 4 && string(header[0:4]) == "\x1A\x45\xDF\xA3" {
		return "video/x-matroska", "high"
	}
	if len(header) >= 4 && string(header[0:4]) == "OggS" {
		return "audio/ogg", "high"
	}
	if len(header) >= 4 && string(header[0:4]) == "fLaC" {
		return "audio/flac", "high"
	}
	if len(header) >= 6 && string(header[0:6]) == "7z\xBC\xAF\x27\x1C" {
		return "application/x-7z-compressed", "high"
	}
	if len(header) >= 6 && string(header[0:6]) == "Rar!\x1A\x07" {
		return "application/vnd.rar", "high"
	}
	if len(header) >= 3 && string(header[0:3]) == "BZh" {
		return "application/x-bzip2", "high"
	}
	if len(header) >= 6 && string(header[0:6]) == "\xFD7zXZ\x00" {
		return "application/x-xz", "high"
	}
	// tar doesn't have a reliable magic number at the start (it's at offset 257), so we rely on ext.

	return mime, "high"
}
