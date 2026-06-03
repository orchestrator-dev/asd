package handlers
import "io"
type CertHandler struct{}
func (h *CertHandler) CanHandle(mime, ext string) bool { return false }
func (h *CertHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error { return nil }
