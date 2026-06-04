package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type LogHandler struct{}

func (h *LogHandler) CanHandle(mime, ext string) bool {
	ext = strings.ToLower(ext)
	if ext == ".log" {
		return true
	}
	return false
}

func (h *LogHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	scanner := bufio.NewScanner(r)

	// Styles
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Grey
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // Green
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))   // Yellow
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))    // Red
	debugStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))  // Blue

	if opts.NoColor || opts.Plain {
		timeStyle = lipgloss.NewStyle()
		infoStyle = lipgloss.NewStyle()
		warnStyle = lipgloss.NewStyle()
		errStyle = lipgloss.NewStyle()
		debugStyle = lipgloss.NewStyle()
	}

	// Regex to match timestamps and log levels. This is a very generic regex.
	// Matches common date formats at the start, and log levels like [INFO], ERROR:, etc.
	timeRe := regexp.MustCompile(`^([\d-]{10}T?[\d:.,]{8,}Z?|[\w]{3}\s+\d+\s+[\d:]{8})`)
	levelRe := regexp.MustCompile(`(?i)\b(INFO|WARN(?:ING)?|ERR(?:OR)?|DEBUG|TRACE|FATAL)\b`)

	// Need to check diff if we want gutter
	diff := getGitDiff(meta.Name)
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	modStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	if opts.NoColor || opts.Plain {
		lineStyle = lipgloss.NewStyle()
		addStyle = lipgloss.NewStyle()
		modStyle = lipgloss.NewStyle()
		delStyle = lipgloss.NewStyle()
	}

	lineNum := 1
	var lines []string

	// Read all lines first to calculate numWidth if not following
	if !opts.Follow {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		r = bytes.NewReader(data)
		scanner = bufio.NewScanner(r)
		lines = strings.Split(string(data), "\n")
	}

	numWidth := 3
	if len(lines) > 0 {
		numWidth = len(fmt.Sprintf("%d", len(lines)))
		if numWidth < 3 {
			numWidth = 3
		}
	}

	for scanner.Scan() {
		text := scanner.Text()

		// Colorize Timestamp
		text = timeRe.ReplaceAllStringFunc(text, func(m string) string {
			return timeStyle.Render(m)
		})

		// Colorize Log Levels
		text = levelRe.ReplaceAllStringFunc(text, func(m string) string {
			mUpper := strings.ToUpper(m)
			if strings.HasPrefix(mUpper, "INFO") {
				return infoStyle.Render(m)
			} else if strings.HasPrefix(mUpper, "WARN") {
				return warnStyle.Render(m)
			} else if strings.HasPrefix(mUpper, "ERR") || strings.HasPrefix(mUpper, "FATAL") {
				return errStyle.Render(m)
			} else if strings.HasPrefix(mUpper, "DEBUG") || strings.HasPrefix(mUpper, "TRACE") {
				return debugStyle.Render(m)
			}
			return m
		})

		if opts.Clean {
			fmt.Fprintln(w, text)
		} else {
			gutter := " "
			d := diff[lineNum]
			if d == "+" {
				gutter = addStyle.Render("+")
			} else if d == "~" {
				gutter = modStyle.Render("~")
			} else if d == "-" {
				gutter = delStyle.Render("-")
			}
			fmt.Fprintf(w, "%s %s │ %s\n", lineStyle.Render(fmt.Sprintf("%*d", numWidth, lineNum)), gutter, text)
		}
		lineNum++
	}

	return scanner.Err()
}
