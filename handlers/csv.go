package handlers

type CSVHandler struct{}

func (h *CSVHandler) CanHandle(mime, ext string) bool {
	return ext == ".csv" || ext == ".tsv" || mime == "text/csv"
}
