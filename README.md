# asd

`asd` is a smart universal file viewer in Go. It replaces `cat` with smart, type-aware file rendering in the terminal.

## Installation

(To be populated)

## Usage

```bash
asd [flags] [file...]
```

Flags:
- `-f`, `--flat`: bypass smart rendering; behave like cat
- `-n`, `--lines`: show line numbers (text/source only)
- `-p`, `--plain`: disable all color and styling
- `--theme STR`: chroma highlight theme (default: auto)
- `--version`: print version and exit
- `-h`, `--help`: usage