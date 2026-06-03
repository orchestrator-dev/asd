package detect

import (
	"os"
	"path/filepath"
	"strings"
)

// ExtDetector uses file extensions to suggest MIME types.
type ExtDetector struct{}

func (d *ExtDetector) Detect(header []byte, filename string) (string, string) {
	if filename != "" && filename != "stdin" && filename != "-" {
		info, err := os.Lstat(filename)
		if err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				return "inode/symlink", "high"
			}
			if info.IsDir() {
				return "inode/directory", "high"
			}
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))
	// Basic mapping, can be expanded later
	switch ext {
	case ".csv", ".tsv":
		return "text/csv", "medium"
	case ".md":
		return "text/markdown", "medium"
	case ".xml":
		return "text/xml", "medium"
	case ".html", ".htm":
		return "text/html", "medium"
	case ".go", ".py", ".js", ".ts", ".css", ".json", ".yaml", ".yml", ".toml":
		return "text/plain", "medium"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "high"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "high"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation", "high"
	case ".odt":
		return "application/vnd.oasis.opendocument.text", "high"
	case ".ods":
		return "application/vnd.oasis.opendocument.spreadsheet", "high"
	case ".odp":
		return "application/vnd.oasis.opendocument.presentation", "high"
	case ".pem", ".crt", ".cer", ".der":
		return "application/x-x509-ca-cert", "high"
	case ".key":
		return "application/pkcs8", "high"
	case ".p12", ".pfx":
		return "application/x-pkcs12", "high"
	case ".pub":
		return "application/ssh-public-key", "high"
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".svg":
		return "image/" + strings.TrimPrefix(ext, "."), "medium"
	case ".mp3", ".flac", ".ogg", ".wav", ".aac", ".m4a", ".opus":
		return "audio/" + strings.TrimPrefix(ext, "."), "medium"
	case ".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv":
		return "video/" + strings.TrimPrefix(ext, "."), "medium"
	case ".pdf":
		return "application/pdf", "high"
	case ".zip":
		return "application/zip", "high"
	case ".tar":
		return "application/x-tar", "high"
	case ".gz", ".tgz":
		return "application/gzip", "high"
	case ".bz2":
		return "application/x-bzip2", "high"
	case ".xz":
		return "application/x-xz", "high"
	case ".7z":
		return "application/x-7z-compressed", "high"
	case ".rar":
		return "application/vnd.rar", "high"
	default:
		return "", "none"
	}
}
