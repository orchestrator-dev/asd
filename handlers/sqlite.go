package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteHandler struct{}

func (h *SQLiteHandler) CanHandle(mime, ext string) bool {
	return ext == ".sqlite" || ext == ".db"
}

func (h *SQLiteHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	path := meta.Name
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("sqlite handler requires a valid file path on disk: %v", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, tableName := range tables {
		var count int
		_ = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)

		schemaRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s);", tableName))
		if err != nil {
			return err
		}

		var columns []string
		for schemaRows.Next() {
			var cid int
			var name string
			var dtype string
			var notnull int
			var dfltValue any
			var pk int
			if err := schemaRows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk); err != nil {
				schemaRows.Close()
				return err
			}
			columns = append(columns, fmt.Sprintf("%s (%s)", name, dtype))
		}
		schemaRows.Close()

		if !opts.NoColor && !opts.Plain {
			tableName = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Render(tableName)
		}
		fmt.Fprintf(w, "Table: %s (%d rows)\n", tableName, count)
		fmt.Fprintf(w, "Columns: %s\n\n", strings.Join(columns, ", "))
	}

	return nil
}
