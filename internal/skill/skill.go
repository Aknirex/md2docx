package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const skillMarkdownTemplate = `---
name: md2docx
description: Convert Markdown files to professional DOCX documents with region-specific style presets. Supports US Business, CN Official, JP Formal, EU Clean, KR Standard, Academic, and custom JSON style templates.
location: file:///<SKILL_DIR>/SKILL.md
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
` + "```bash" + `
go install github.com/md2docx/cli/cmd/md2docx@latest
` + "```" + `

### Via direct download
Download the latest release for your platform from:
https://github.com/md2docx/cli/releases/latest

Place the binary in your PATH (e.g., ` + "`/usr/local/bin`" + ` on macOS/Linux, ` + "`C:\\Windows\\System32`" + ` on Windows).

### Verify installation
` + "```bash" + `
md2docx version
` + "```" + `

## Usage (Agent Mode)

For agent/automation use, always include ` + "`--json`" + ` for structured output:

### Convert with default style
` + "```bash" + `
md2docx convert -i input.md -o output.docx --json
` + "```" + `

### Convert with a specific preset
` + "```bash" + `
md2docx convert -i input.md -o output.docx -s cn-official --json
` + "```" + `

### Convert with Mermaid diagrams rendered as images
` + "```bash" + `
md2docx convert -i input.md -o output.docx --mermaid --json
md2docx convert -i input.md -o output.docx --mermaid --mermaid-theme dark --json
` + "```" + `

### List available style presets
` + "```bash" + `
md2docx presets --json
` + "```" + `

### Show details of a preset
` + "```bash" + `
md2docx preset cn-official --json
` + "```" + `

### Convert with a custom template
` + "```bash" + `
md2docx convert -i input.md -o output.docx -s /path/to/template.json --json
` + "```" + `

### Create a template from a preset
` + "```bash" + `
md2docx template create -o my-template.json -s cn-official --json
` + "```" + `

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

## Output Format

When using ` + "`--json`" + `, the output is structured JSON:

**Success:**
` + "```json" + `
{
  "success": true,
  "outputPath": "/path/to/output.docx",
  "bytes": 12345
}
` + "```" + `

**Error:**
` + "```json" + `
{
  "success": false,
  "error": "error message"
}
` + "```" + `

## Requirements

- No external dependencies (no Word, Pandoc, or LibreOffice required)
- The binary is self-contained; just download and run
`

// GenerateSkillFile generates the SKILL.md content for the current platform.
func GenerateSkillContent(installDir string) string {
	content := strings.ReplaceAll(skillMarkdownTemplate, "<SKILL_DIR>", installDir)
	return content
}

// InstallToDir installs the skill to the specified directory by writing SKILL.md.
func InstallToDir(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating skill directory %s: %w", dir, err)
	}

	skillPath := filepath.Join(dir, "SKILL.md")
	content := GenerateSkillContent(dir)

	if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("writing SKILL.md: %w", err)
	}

	return skillPath, nil
}

// FindKiloSkillsDir attempts to locate the .kilo/skills directory.
// It searches from the current directory upward, and also checks ~/.config/kilo/skills.
func FindKiloSkillsDir() ([]string, error) {
	var candidates []string

	// 1. Search upward from CWD for .kilo/skills
	cwd, _ := os.Getwd()
	for dir := cwd; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		skillsDir := filepath.Join(dir, ".kilo", "skills")
		if info, err := os.Stat(skillsDir); err == nil && info.IsDir() {
			candidates = append(candidates, skillsDir)
		}
	}

	// 2. Check ~/.config/kilo/skills (global skills)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalDir := filepath.Join(homeDir, ".config", "kilo", "skills")
		if info, err := os.Stat(globalDir); err == nil && info.IsDir() {
			candidates = append(candidates, globalDir)
		}
	}

	return candidates, nil
}

// Install auto-discovers the .kilo/skills directory and installs there.
func Install() (string, error) {
	dirs, err := FindKiloSkillsDir()
	if err != nil {
		return "", err
	}

	if len(dirs) == 0 {
		// Create in home config as fallback
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot find home directory and no .kilo/skills found")
		}
		target := filepath.Join(homeDir, ".config", "kilo", "skills", "md2docx")
		return InstallToDir(target)
	}

	// Use the first discovered .kilo/skills directory
	target := filepath.Join(dirs[0], "md2docx")
	return InstallToDir(target)
}

// InstallToPath installs the skill to an explicit path.
func InstallToPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolving path: %w", err)
	}
	return InstallToDir(absPath)
}

// binaryName returns the platform-appropriate binary name.
func binaryName() string {
	if runtime.GOOS == "windows" {
		return "md2docx.exe"
	}
	return "md2docx"
}
