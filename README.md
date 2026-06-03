<div align="center">
  <img src="assets/hero_banner.png" alt="asd - A smart universal file viewer in Go" />

  <h1>asd</h1>
  
  <p><b>A smart universal file viewer for the modern terminal.</b></p>
  
  [![Build Status](https://github.com/vsmanu/asd/actions/workflows/ci.yml/badge.svg)](https://github.com/vsmanu/asd/actions/workflows/ci.yml)
  [![Release](https://img.shields.io/github/v/release/vsmanu/asd)](https://github.com/vsmanu/asd/releases)
  [![Go Report Card](https://goreportcard.com/badge/github.com/vsmanu/asd)](https://goreportcard.com/report/github.com/vsmanu/asd)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

<br/>

`asd` replaces traditional tools like `cat` with smart, type-aware file rendering. By automatically detecting file types via magic bytes (and falling back to extensions), `asd` formats and renders them with beautiful syntax highlighting, layout tables, tree views, and media metadata. 

When you output more text than your terminal can fit, `asd` instantly buffers into an interactive pager. If you pipe the output to another tool, it smartly defaults back to raw bytes (behaving exactly like `cat`).

<div align="center">
  <img src="assets/screenshot.png" alt="asd terminal screenshot" />
</div>

## 🚀 Features

- **Code & Text**: Beautiful syntax highlighting using [Chroma](https://github.com/alecthomas/chroma), with line numbers and theme support.
- **Structured Data**: JSON, YAML, TOML formatting and valid syntax checking.
- **Data Tables**: CSV and TSV formatted into beautiful, auto-sizing terminal tables via [Lipgloss](https://github.com/charmbracelet/lipgloss).
- **Markdown**: Fully styled Markdown rendering using [Glamour](https://github.com/charmbracelet/glamour).
- **Archives**: Explores zip, tar, tar.gz, 7z, and rar as interactive directory trees showing uncompressed stats.
- **Office Docs & PDF**: Automatically extracts and parses text from DOCX, XLSX, PPTX, ODT, and PDFs without you needing to extract them manually!
- **Media**: Peeks into Image, Audio, and Video files to reveal EXIF, track metadata, resolutions, codecs, and durations.
- **Security**: Displays X.509 PEM/CRT properties and SSH Key parameters effortlessly.
- **Directories**: View styled `ls -la` directory trees and resolve symlinks.
- **Hex Dump**: Safely view unknown binaries with built-in hex-dumping.

## 📦 Installation

### Using Homebrew (macOS/Linux)
\`\`\`bash
brew tap vsmanu/asd
brew install asd
\`\`\`

### Using Go
Ensure you have Go `1.23` or higher.
\`\`\`bash
go install github.com/vsmanu/asd@latest
\`\`\`

### Download Pre-built Binaries
Head over to the [Releases](https://github.com/vsmanu/asd/releases) page to download pre-built binaries for your OS and architecture (Windows, macOS, Linux).

*Note: The GitHub Actions workflows are actively compiling `v0.8.3`, so the binaries will appear on the release page momentarily!*

## 📖 Usage

Using `asd` is as simple as using `cat`.

\`\`\`bash
asd [flags] [file...]
\`\`\`

**Examples:**
\`\`\`bash
# Read a simple text or code file (syntax highlighted)
asd main.go

# Read multiple files at once
asd config.yaml payload.json

# Read a compressed archive without unzipping
asd data.zip

# Pipe from another program (reads stdin)
echo "{\"test\": 123}" | asd
\`\`\`

### Flags

| Flag | Long Flag | Description |
|------|-----------|-------------|
| `-f` | `--flat`  | Bypass all smart rendering and act identically to `cat`. |
| `-n` | `--lines` | Show line numbers alongside text and source code. |
| `-p` | `--plain` | Disable all color and styling (good for older terminals). |
| | `--theme` | Specify a chroma highlight theme (default is `auto`). |
| `-h` | `--help`  | Print usage and help menus. |
| | `--version` | Print the current version and exit. |


## ⚙️ Configuration

You can customize `asd` permanently using a global config file located at `~/.config/asd/config.toml`:

\`\`\`toml
# Set your preferred Chroma theme for syntax highlighting
theme = "dracula"
\`\`\`

## 🛠 Compilation and Development

\`\`\`bash
git clone https://github.com/vsmanu/asd.git
cd asd
go mod download
make build
\`\`\`

You can run the tests using:
\`\`\`bash
make test
\`\`\`

## 📄 License

This project is licensed under the MIT License.