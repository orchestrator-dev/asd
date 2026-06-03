sed -i -e '/func (h \*CSVHandler) CanHandle/,/}/c\
func (h *CSVHandler) CanHandle(mime, ext string) bool {\
	return ext == ".csv" || ext == ".tsv" || mime == "text/csv"\
}' handlers/csv.go
