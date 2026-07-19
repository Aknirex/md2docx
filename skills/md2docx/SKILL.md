---
name: md2docx
description: Convert Markdown files to professional DOCX documents with region-specific style presets. Supports US Business, CN Official (公文风格), JP Formal, EU Clean, KR Standard, Academic, and custom JSON style templates. No Word or Pandoc required.
location: file:///D:/00yw/environments/md2docx/skills/md2docx/SKILL.md
---

# md2docx

Convert Markdown to DOCX (Open XML) documents — dependency-free, no Word or Pandoc required.

## When to Use

- When the user asks to convert a `.md` file to `.docx`
- When the user needs to generate styled Word documents from Markdown
- When the user mentions "convert to word", "markdown to docx", "md to docx"
- When the user wants styled documents for US, CN, JP, EU, KR, or Academic contexts

## Installation

The `md2docx` binary should be available on your PATH. To install:

### Via Go install
```bash
go install github.com/md2docx/cli/cmd/md2docx@latest
```

### Via direct download
Download the latest release for your platform from:
https://github.com/md2docx/cli/releases/latest

Place the binary in your PATH (e.g., `/usr/local/bin` on macOS/Linux, `C:\Windows\System32` on Windows).

### Verify installation
```bash
md2docx version
```

## Usage (Agent Mode)

For agent/automation use, always include `--json` for structured output:

### Convert with default style
```bash
md2docx convert -i input.md -o output.docx --json
```

### Convert with a specific preset
```bash
md2docx convert -i input.md -o output.docx -s cn-official --json
```

### List available style presets
```bash
md2docx presets --json
```

### Show details of a preset
```bash
md2docx preset cn-official --json
```

### Convert with a custom template
```bash
md2docx convert -i input.md -o output.docx -s /path/to/template.json --json
```

### Create a template from a preset
```bash
md2docx template create -o my-template.json -s cn-official --json
```

## Built-in Style Presets

| Preset        | Target Region | Fonts                                    |
|---------------|---------------|------------------------------------------|
| us-business   | US            | Cambria / Calibri / Consolas             |
| us-modern     | US            | Segoe UI / Cascadia Code                 |
| cn-official   | China         | SimHei / SimSun (公文风格)                |
| cn-modern     | China         | Noto Sans SC / Noto Sans Mono SC         |
| jp-formal     | Japan         | Yu Mincho / Yu Gothic                    |
| eu-clean      | Europe        | Helvetica / Arial / Fira Code            |
| kr-standard   | Korea         | Malgun Gothic / Nanum Gothic / D2Coding  |
| academic      | Global        | Times New Roman / Courier New            |
| default       | Global        | Aptos Display / Cascadia Mono            |

## Requirements

- No external dependencies (no Word, Pandoc, or LibreOffice required)
- The binary is self-contained; just download and run

## Output Format (JSON)

When using `--json`, the output is structured JSON:

**Success:**
```json
{
  "success": true,
  "outputPath": "/path/to/output.docx",
  "bytes": 12345
}
```

**Error:**
```json
{
  "success": false,
  "error": "error message"
}
```
