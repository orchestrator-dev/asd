package handlers

import "io"

type MarkdownHandler struct{}

func (h *MarkdownHandler) CanHandle(mime, ext string) bool {
	return ext == ".md" || mime == "text/markdown"
}
func (h *MarkdownHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	return nil
}
