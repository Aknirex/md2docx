---
name: md2docx
description: Convert Markdown files to professional DOCX documents with region-specific style presets. Supports mermaid diagram rendering, US/CN/JP/EU/KR/Academic styles. No Word or Pandoc required.
---

# md2docx — Markdown to DOCX Converter

A dependency-free CLI tool that converts Markdown to professional DOCX (Open XML) documents. Supports built-in regional style presets and Mermaid diagram rendering.

## Tool Installation

Before using, ensure the **md2docx** binary is available on your PATH.

### Install via Go

```bash
go install github.com/Aknirex/md2docx/cmd/md2docx@latest
```

### Direct Download

Download the latest binary for your platform from:
[GitHub Releases](https://github.com/Aknirex/md2docx/releases/latest)

- **Linux**: `md2docx-linux-amd64` / `md2docx-linux-arm64`
- **macOS**: `md2docx-darwin-amd64` / `md2docx-darwin-arm64`  
- **Windows**: `md2docx-windows-amd64.exe`

Place the binary in your PATH.

### Package Managers

```bash
# Homebrew (macOS / Linux)
brew install Aknirex/homebrew-tap/md2docx

# Debian / Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

Verify: `md2docx version`

## Usage (Agent Mode — always use `--json`)

### Convert Markdown to DOCX

```bash
md2docx convert -i <input.md> -o <output.docx> --json
```

### Convert with regional style preset

```bash
md2docx convert -i input.md -o output.docx -s cn-official --json
```

### Convert with Mermaid diagrams rendered as images

```bash
md2docx convert -i input.md -o output.docx --mermaid --mermaid-theme default --json
```

### List available style presets

```bash
md2docx presets --json
```

### Create custom style template

```bash
md2docx template create -o my-style.json -s us-business --json
```

## Style Presets

| Preset       | Region  | Characteristic                             |
|--------------|---------|--------------------------------------------|
| us-business  | US      | Cambria/Calibri, professional blue accent  |
| us-modern    | US      | Segoe UI, minimal dark tones               |
| cn-official  | China   | 小标宋_GBK/仿宋_GB2312/楷体_GB2312 (公文风格), red accent |
| cn-modern    | China   | Noto Sans SC, modern Chinese               |
| jp-formal    | Japan   | Yu Mincho/Yu Gothic, business formal       |
| eu-clean     | Europe  | Helvetica/Arial, clean minimalist          |
| kr-standard  | Korea   | Malgun Gothic/Nanum Gothic                 |
| academic     | Global  | Times New Roman, scholarly                 |
| default      | Global  | Aptos Display/Cascadia Mono                |

## Mermaid Rendering

When `--mermaid` is set, ```` ```mermaid ```` blocks are rendered as embedded PNG images via the public [mermaid.ink](https://mermaid.ink) API:

- `--mermaid-theme`: `default`, `neutral`, `dark`, `forest`
- `--mermaid-server`: custom self-hosted mermaid.ink URL

## JSON Output Format

**Success:**
```json
{"success": true, "outputPath": "/path/to/output.docx", "bytes": 12345}
```

**Error:**
```json
{"success": false, "error": "descriptive error message"}
```

## Requirements

- No Word, Pandoc, or LibreOffice needed
- Mermaid rendering requires network access (mermaid.ink)
- Self-contained static binary, zero runtime dependencies
- Install the skill via: `npx skills add Aknirex/md2docx`
