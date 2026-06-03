package handlers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"golang.org/x/net/html"
)

type XMLHandler struct{}

func (h *XMLHandler) CanHandle(mime, ext string) bool {
	return strings.HasPrefix(mime, "text/xml") || strings.HasPrefix(mime, "text/html") || ext == ".xml" || ext == ".html" || ext == ".htm"
}

func (h *XMLHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	isHTML := strings.HasSuffix(strings.ToLower(meta.Name), ".html") || strings.HasSuffix(strings.ToLower(meta.Name), ".htm")

	if isHTML {
		// print title, meta description, link count as header
		title, desc, linkCount := h.parseHTMLMeta(bytes.NewReader(data))
		fmt.Fprintf(w, "Title: %s\n", title)
		fmt.Fprintf(w, "Description: %s\n", desc)
		fmt.Fprintf(w, "Links: %d\n\n", linkCount)
	}

	// Pretty-print with indentation
	var buf bytes.Buffer
	if isHTML {
		// HTML formatting isn't cleanly supported by xml.Encoder, just output as-is or use xml encoder but it might break on unclosed tags.
		// Actually, let's just colorize the raw HTML. The prompt says "Pretty-print with indentation + tag coloring" for XML.
		// "HTML (in xml.go or html.go): also print title, meta description, link count as header."
		// For HTML, let's just colorize the original or we can try formatting if it's XHTML. Let's just colorize original for HTML.
		buf.Write(data)
	} else {
		// XML pretty print
		decoder := xml.NewDecoder(bytes.NewReader(data))
		encoder := xml.NewEncoder(&buf)
		encoder.Indent("", "  ")
		for {
			t, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				// fallback to raw data if not well-formed
				buf.Reset()
				buf.Write(data)
				break
			}
			encoder.EncodeToken(t)
		}
		encoder.Flush()
	}

	// Tag coloring
	lexer := lexers.Match(meta.Name)
	if lexer == nil {
		if isHTML {
			lexer = lexers.Get("html")
		} else {
			lexer = lexers.Get("xml")
		}
	}

	style := styles.Get(opts.Theme)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if opts.NoColor || opts.Plain {
		formatter = formatters.Get("noop")
	}

	iterator, err := lexer.Tokenise(nil, buf.String())
	if err != nil {
		return err
	}

	return formatter.Format(w, style, iterator)
}

func (h *XMLHandler) parseHTMLMeta(r io.Reader) (title, desc string, linkCount int) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			var isDesc bool
			var content string
			for _, a := range n.Attr {
				if a.Key == "name" && strings.ToLower(a.Val) == "description" {
					isDesc = true
				}
				if a.Key == "content" {
					content = a.Val
				}
			}
			if isDesc {
				desc = content
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			linkCount++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return
}
