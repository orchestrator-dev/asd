package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

type TextHandler struct{}

func (h *TextHandler) CanHandle(mime, ext string) bool {
	return mime == "text/plain" || strings.HasPrefix(mime, "text/")
}

func (h *TextHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	if opts.Follow {
		// Stream line by line for tail -f
		scanner := bufio.NewScanner(r)

		lexer := lexers.Match(meta.Name)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		style := styles.Get(opts.Theme)
		if style == nil {
			style = styles.Fallback
		}
		formatter := formatters.Get("terminal256")
		if opts.NoColor || opts.Plain {
			formatter = formatters.Get("noop")
		}

		for scanner.Scan() {
			line := scanner.Text()
			iterator, _ := lexer.Tokenise(nil, line)
			formatter.Format(w, style, iterator)
			fmt.Fprintln(w)
		}
		return scanner.Err()
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Sniff language
	lexer := lexers.Match(meta.Name)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(opts.Theme)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if opts.NoColor || opts.Plain {
		formatter = formatters.Get("noop")
	}

	iterator, err := lexer.Tokenise(nil, string(data))
	if err != nil {
		return err
	}

	if opts.Clean {
		// Clean mode: just format without any UI decorations
		return formatter.Format(w, style, iterator)
	}

	// Rich mode: Split into lines and add Git gutter & line numbers
	lines := chroma.SplitTokensIntoLines(iterator.Tokens())
	diff := getGitDiff(meta.Name)

	numWidth := len(strconv.Itoa(len(lines)))
	if numWidth < 3 {
		numWidth = 3
	}

	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Grey line numbers
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))    // Green
	modStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))    // Yellow
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))    // Red

	if opts.NoColor || opts.Plain {
		lineStyle = lipgloss.NewStyle()
		addStyle = lipgloss.NewStyle()
		modStyle = lipgloss.NewStyle()
		delStyle = lipgloss.NewStyle()
	}

	for i, lineTokens := range lines {
		lineNum := i + 1

		// Format the syntax highlighted line
		var buf bytes.Buffer
		formatter.Format(&buf, style, chroma.Literator(lineTokens...))
		code := strings.TrimRight(buf.String(), "\r\n")

		// Determine Git Gutter
		gutter := " "
		d := diff[lineNum]
		if d == "+" {
			gutter = addStyle.Render("+")
		} else if d == "~" {
			gutter = modStyle.Render("~")
		} else if d == "-" {
			gutter = delStyle.Render("-")
		}

		// Print Line Number, Gutter, and Code
		fmt.Fprintf(w, "%s %s │ %s\n", lineStyle.Render(fmt.Sprintf("%*d", numWidth, lineNum)), gutter, code)
	}

	return nil
}

func getGitDiff(filename string) map[int]string {
	diff := make(map[int]string)
	if filename == "" || filename == "stdin" || filename == "-" {
		return diff
	}

	cmd := exec.Command("git", "diff", "-U0", "--", filename)
	out, err := cmd.Output()
	if err != nil {
		// Possibly untracked file. Let's check git ls-files.
		cmdLs := exec.Command("git", "ls-files", "--error-unmatch", filename)
		if errLs := cmdLs.Run(); errLs != nil {
			// Untracked. We could mark all as added, but that would be noisy.
			return diff
		}
		return diff
	}

	lines := strings.Split(string(out), "\n")
	re := regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			oldCount := 1
			if matches[2] != "" {
				oldCount, _ = strconv.Atoi(matches[2])
			}

			newStart, _ := strconv.Atoi(matches[3])
			newCount := 1
			if matches[4] != "" {
				newCount, _ = strconv.Atoi(matches[4])
			}

			if newCount == 0 {
				diff[newStart] = "-"
			} else {
				mark := "+"
				if oldCount > 0 {
					mark = "~"
				}
				for i := 0; i < newCount; i++ {
					diff[newStart+i] = mark
				}
			}
		}
	}
	return diff
}
