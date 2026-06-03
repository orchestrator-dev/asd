# asd

`asd` is a smart universal file viewer in Go. It replaces `cat` with smart, type-aware file rendering in the terminal.
It automatically detects file types and renders them with beautiful syntax highlighting, formatted tables, tree views, and image/media metadata.

## Features

- **Text & Source Code**: Syntax highlighting using Chroma.
- **Structured Data**: CSV/TSV table rendering, JSON/YAML/TOML pretty-printing.
- **Markdown**: Styled markdown rendering.
- **Archives**: Directory tree listing for zip, tar, gz, 7z, rar.
- **Media**: Image dimensions/metadata, audio/video metadata.
- **Office Docs & PDF**: Text extraction from PDF, DOCX, XLSX, PPTX.
- **Security**: X.509 cert info parsing.
- **Smart Pager**: Built-in interactive pager when output exceeds terminal height.
- **Fallback**: Hex dump for unknown binary files.

## Installation

### Using Homebrew (macOS/Linux)
\`\`\`bash
brew install orchestrator-dev/asd/asd
\`\`\`

### Using Go
\`\`\`bash
go install github.com/orchestrator-dev/asd@latest
\`\`\`

### Manual
Download the latest release from the [Releases](https://github.com/orchestrator-dev/asd/releases) page.

## Usage

\`\`\`bash
asd [flags] [file...]
\`\`\`

Flags:
- `-f`, `--flat`: bypass smart rendering; behave like cat
- `-n`, `--lines`: show line numbers (text/source only)
- `-p`, `--plain`: disable all color and styling
- `--theme STR`: chroma highlight theme (default: auto)
- `--version`: print version and exit
- `-h`, `--help`: usage

## Configuration

You can configure `asd` using `~/.config/asd/config.toml`:

\`\`\`toml
theme = "monokai"
\`\`\`