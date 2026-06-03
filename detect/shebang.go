package detect

import "strings"

// ShebangDetector checks for a shebang line to suggest text type.
type ShebangDetector struct{}

func (d *ShebangDetector) Detect(header []byte, filename string) (string, string) {
	if len(header) >= 2 && header[0] == '#' && header[1] == '!' {
		return "text/plain", "high"
	}
	// Check for XML declaration as a fallback text hint
	if len(header) >= 5 && strings.HasPrefix(string(header), "<?xml") {
		return "text/xml", "high"
	}
	return "", "none"
}
