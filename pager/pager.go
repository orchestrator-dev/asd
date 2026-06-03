package pager

import (
	"bytes"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type Pager struct {
	buf bytes.Buffer
	out io.Writer
}

func New(out io.Writer) *Pager {
	return &Pager{out: out}
}

func (p *Pager) Write(b []byte) (int, error) {
	return p.buf.Write(b)
}

func (p *Pager) Close() error {
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if !isTTY {
		_, err := p.out.Write(p.buf.Bytes())
		return err
	}

	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		_, err = p.out.Write(p.buf.Bytes())
		return err
	}

	lines := strings.Split(p.buf.String(), "\n")
	if len(lines) <= h {
		_, err = p.out.Write(p.buf.Bytes())
		return err
	}

	m := &pagerModel{
		lines:  lines,
		height: h,
	}

	if _, err := tea.NewProgram(m, tea.WithOutput(p.out), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

type pagerModel struct {
	lines  []string
	width  int
	height int
	offset int
}

func (m *pagerModel) Init() tea.Cmd { return nil }

func (m *pagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			if m.offset < len(m.lines)-m.height+1 {
				m.offset++
			}
		case "pgup", "b":
			m.offset -= m.height - 1
			if m.offset < 0 {
				m.offset = 0
			}
		case "pgdown", "f", " ":
			m.offset += m.height - 1
			max := len(m.lines) - m.height + 1
			if m.offset > max {
				m.offset = max
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *pagerModel) View() string {
	var b strings.Builder
	end := m.offset + m.height - 1
	if end > len(m.lines) {
		end = len(m.lines)
	}
	for i := m.offset; i < end; i++ {
		b.WriteString(m.lines[i])
		b.WriteString("\n")
	}
	b.WriteString("\x1b[7m : \x1b[0m") // simple prompt
	return b.String()
}
