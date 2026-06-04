package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func explainFile(w io.Writer, data []byte, filename string, theme string) error {
	// Truncate data to avoid sending massive payloads
	content := string(data)
	if len(content) > 10000 {
		content = content[:10000] + "\n... (truncated)"
	}

	prompt := fmt.Sprintf("You are an expert developer assistant. Explain the following file (%s) in 3-4 concise bullet points. Focus on its purpose, key components, and any notable patterns. Keep it short and technical.\n\n%s", filename, content)

	reqBody := OllamaRequest{
		Model:  "llama3", // default or maybe we can auto-detect available models, but llama3 is standard
		Prompt: prompt,
		Stream: false,
	}

	reqBytes, _ := json.Marshal(reqBody)

	// Try Ollama
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		fmt.Fprintf(w, "%s Could not connect to local Ollama (http://localhost:11434).\nMake sure Ollama is running, or disable --explain.\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("✖ Error:"))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama API returned status: %d", resp.StatusCode)
	}

	var oResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return err
	}

	// Render the markdown response
	markdown := "### 🤖 AI Explanation\n" + oResp.Response + "\n---"

	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(theme),
		glamour.WithWordWrap(80),
	)

	out, err := r.Render(markdown)
	if err == nil {
		fmt.Fprint(w, out)
	} else {
		fmt.Fprintln(w, markdown)
	}

	return nil
}
