import os
for name in ['office.go', 'toml.go', 'yaml.go', 'cert.go', 'directory.go', 'csv.go', 'markdown.go', 'archive.go', 'pdf.go', 'symlink.go']:
    if os.path.exists('handlers/' + name):
        h = name.split('.')[0].capitalize() + 'Handler'
        if name == 'csv.go': h = 'CSVHandler'
        if name == 'yaml.go': h = 'YAMLHandler'
        if name == 'toml.go': h = 'TOMLHandler'
        if name == 'pdf.go': h = 'PDFHandler'
        
        # We write a stub
        c = f"""package handlers
import "io"
type {h} struct{{}}
func (h *{h}) CanHandle(mime, ext string) bool {{ return false }}
func (h *{h}) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {{ return nil }}
"""
        with open('handlers/' + name, 'w') as f:
            f.write(c)

