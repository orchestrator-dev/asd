package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

type OfficeHandler struct{}

func (h *OfficeHandler) CanHandle(mime, ext string) bool {
	switch mime {
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.oasis.opendocument.text",
		"application/vnd.oasis.opendocument.spreadsheet",
		"application/vnd.oasis.opendocument.presentation":
		return true
	}
	switch ext {
	case ".docx", ".xlsx", ".pptx", ".odt", ".ods", ".odp":
		return true
	}
	return false
}

func (h *OfficeHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	ext := strings.ToLower(meta.Name)
	if strings.HasSuffix(ext, ".docx") {
		return h.renderDocx(w, data)
	} else if strings.HasSuffix(ext, ".xlsx") {
		return h.renderXlsx(w, data)
	} else if strings.HasSuffix(ext, ".pptx") {
		return h.renderPptx(w, data)
	} else if strings.HasSuffix(ext, ".odt") || strings.HasSuffix(ext, ".ods") || strings.HasSuffix(ext, ".odp") {
		return h.renderOdf(w, data)
	}

	return fmt.Errorf("unsupported office format")
}

func (h *OfficeHandler) renderDocx(w io.Writer, data []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return fmt.Errorf("word/document.xml not found")
	}

	rc, err := docFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	var paraBuf strings.Builder
	var linesCount int

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch se := t.(type) {
		case xml.StartElement:
		case xml.EndElement:
			if se.Name.Local == "p" {
				text := strings.TrimSpace(paraBuf.String())
				if text != "" {
					fmt.Fprintln(w, text)
					linesCount++
					if linesCount >= 200 {
						return nil
					}
				}
				paraBuf.Reset()
			}
		case xml.CharData:
			paraBuf.Write([]byte(se))
		}
	}
	return nil
}

func (h *OfficeHandler) renderXlsx(w io.Writer, data []byte) error {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil
	}
	fmt.Fprintf(w, "Sheets: %s\n\n", strings.Join(sheets, ", "))

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return err
	}

	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return nil
}

func (h *OfficeHandler) renderPptx(w io.Writer, data []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	var slideFiles []*zip.File
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") {
			slideFiles = append(slideFiles, f)
		}
	}

	sort.Slice(slideFiles, func(i, j int) bool {
		return slideFiles[i].Name < slideFiles[j].Name
	})

	for _, f := range slideFiles {
		rc, err := f.Open()
		if err != nil {
			continue
		}

		fmt.Fprintf(w, "--- %s ---\n", f.Name)
		decoder := xml.NewDecoder(rc)
		var paraBuf strings.Builder
		for {
			t, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}

			switch se := t.(type) {
			case xml.EndElement:
				if se.Name.Local == "p" {
					text := strings.TrimSpace(paraBuf.String())
					if text != "" {
						fmt.Fprintln(w, text)
					}
					paraBuf.Reset()
				}
			case xml.CharData:
				paraBuf.Write([]byte(se))
			}
		}
		rc.Close()
		fmt.Fprintln(w)
	}
	return nil
}

func (h *OfficeHandler) renderOdf(w io.Writer, data []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	var contentFile *zip.File
	for _, f := range zr.File {
		if f.Name == "content.xml" {
			contentFile = f
			break
		}
	}
	if contentFile == nil {
		return fmt.Errorf("content.xml not found")
	}

	rc, err := contentFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	var paraBuf strings.Builder
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch se := t.(type) {
		case xml.EndElement:
			if se.Name.Local == "p" {
				text := strings.TrimSpace(paraBuf.String())
				if text != "" {
					fmt.Fprintln(w, text)
				}
				paraBuf.Reset()
			}
		case xml.CharData:
			paraBuf.Write([]byte(se))
		}
	}
	return nil
}
