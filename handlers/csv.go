package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type CSVHandler struct{}

func (h *CSVHandler) CanHandle(mime, ext string) bool {
	return strings.HasPrefix(mime, "text/csv") || ext == ".csv" || ext == ".tsv"
}

func (h *CSVHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	reader := csv.NewReader(r)
	if strings.HasSuffix(strings.ToLower(meta.Name), ".tsv") {
		reader.Comma = '\t'
	}
	// Set fields per record to -1 to allow variable number of fields, just in case
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}

	numRows := len(records)
	numCols := 0
	for _, rec := range records {
		if len(rec) > numCols {
			numCols = len(rec)
		}
	}

	// Truncate fields to max 40 chars
	for i := range records {
		for j := range records[i] {
			runes := []rune(records[i][j])
			if len(runes) > 40 {
				records[i][j] = string(runes[:37]) + "..."
			}
		}
	}

	var t *table.Table
	if len(records) > 1 {
		t = table.New().
			Headers(records[0]...).
			Rows(records[1:]...)
	} else {
		t = table.New().
			Headers(records[0]...)
	}

	t.Border(lipgloss.NormalBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true)
			}
			return lipgloss.NewStyle()
		})

	fmt.Fprintln(w, t.String())
	fmt.Fprintf(w, "%d rows x %d columns\n", numRows, numCols)

	return nil
}
