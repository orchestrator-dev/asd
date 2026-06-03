sed -i -e '/func (h \*CSVHandler) Render/,/}/c\
func (h *CSVHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {\
	_, err := io.Copy(w, r)\
	return err\
}' handlers/csv.go

cat << 'INNER' >> handlers/csv.go
func (h *CSVHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	_, err := io.Copy(w, r)
	return err
}
INNER
go run golang.org/x/tools/cmd/goimports@latest -w handlers/csv.go
