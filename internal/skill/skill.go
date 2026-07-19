package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Agent detection — known agent skill directory patterns
// ---------------------------------------------------------------------------

// AgentInfo describes a supported AI coding agent and where it looks for skills.
type AgentInfo struct {
	Name        string   // e.g. "kilo", "claude-code"
	Label       string   // human-readable, e.g. "Kilo"
	SkillSubDir string   // subdirectory name inside the skill dir, e.g. "md2docx"
	Dirs        []string // candidate skill-root directories (checked for existence)
}

// KnownAgents returns the list of agents that md2docx can install skills for.
func KnownAgents() []AgentInfo {
	home, _ := os.UserHomeDir()
	return []AgentInfo{
		{
			Name:        "kilo",
			Label:       "Kilo",
			SkillSubDir: "md2docx",
			Dirs:        skillDirCandidates(".kilo/skills", home, ".config/kilo/skills"),
		},
		{
			Name:        "claude",
			Label:       "Claude Code",
			SkillSubDir: "md2docx",
			Dirs:        skillDirCandidates(".claude/skills", home, ".claude/skills"),
		},
		{
			Name:        "kilocode",
			Label:       "KiloCode",
			SkillSubDir: "md2docx",
			Dirs:        skillDirCandidates(".kilocode/skills", home, ".kilocode/skills"),
		},
	}
}

// skillDirCandidates returns candidate paths: project-local (relative) and global (absolute in home).
func skillDirCandidates(projectRel string, home string, globalRel string) []string {
	cwd, _ := os.Getwd()

	var candidates []string
	// Project-local: walk upward from CWD
	for dir := cwd; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		candidates = append(candidates, filepath.Join(dir, projectRel))
	}
	// Global in home
	if home != "" {
		candidates = append(candidates, filepath.Join(home, globalRel))
	}
	return candidates
}

// DetectAgents scans the filesystem and returns agents that have at least one
// existing skill directory, plus their best candidate path.
func DetectAgents() []AgentInstallTarget {
	var targets []AgentInstallTarget
	for _, agent := range KnownAgents() {
		for _, dir := range agent.Dirs {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				targetDir := filepath.Join(dir, agent.SkillSubDir)
				targets = append(targets, AgentInstallTarget{
					Agent:     agent,
					TargetDir: targetDir,
					Exists:    dirExists(targetDir),
				})
				break // first match wins per agent
			}
		}
	}
	return targets
}

// AgentInstallTarget pairs an agent with a concrete installation directory.
type AgentInstallTarget struct {
	Agent     AgentInfo
	TargetDir string // e.g. /home/user/.config/kilo/skills/md2docx
	Exists    bool   // whether SKILL.md already exists there
}

// ---------------------------------------------------------------------------
// Installation tracking
// ---------------------------------------------------------------------------

type installRecord struct {
	Type        string `json:"type"`        // "symlink" or "copy"
	Target      string `json:"target"`      // the installed path
	Source      string `json:"source"`      // canonical SKILL.md source
	InstalledAt string `json:"installedAt"` // ISO 8601
	Agent       string `json:"agent"`       // agent name
}

type installManifest struct {
	Installed []installRecord `json:"installed"`
}

func manifestPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "md2docx")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating config dir: %w", err)
	}
	return filepath.Join(dir, "installations.json"), nil
}

func loadManifest() (*installManifest, error) {
	path, err := manifestPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &installManifest{}, nil
		}
		return nil, err
	}
	var m installManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return &installManifest{}, nil // corrupt → start fresh
	}
	return &m, nil
}

func saveManifest(m *installManifest) error {
	path, err := manifestPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func addInstallRecord(targetDir, sourcePath, agentName, installType string) error {
	m, err := loadManifest()
	if err != nil {
		return err
	}
	// Remove any existing record for this target
	var filtered []installRecord
	for _, r := range m.Installed {
		if r.Target != targetDir {
			filtered = append(filtered, r)
		}
	}
	m.Installed = append(filtered, installRecord{
		Type:        installType,
		Target:      targetDir,
		Source:      sourcePath,
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
		Agent:       agentName,
	})
	return saveManifest(m)
}

func removeInstallRecord(targetDir string) error {
	m, err := loadManifest()
	if err != nil {
		return err
	}
	var filtered []installRecord
	for _, r := range m.Installed {
		if r.Target != targetDir {
			filtered = append(filtered, r)
		}
	}
	m.Installed = filtered
	return saveManifest(m)
}

// ---------------------------------------------------------------------------
// Skill content generation
// ---------------------------------------------------------------------------

// GenerateSkillContent returns the full SKILL.md text.
// The skill describes md2docx as an external tool the agent should install and invoke.
func GenerateSkillContent() string {
	installNotes := binaryInstallInstructions()
	return fmt.Sprintf(`---
name: md2docx
description: Convert Markdown files to professional DOCX documents with region-specific style presets. Supports mermaid diagram rendering, US/CN/JP/EU/KR/Academic styles. No Word or Pandoc required.
---

# md2docx — Markdown to DOCX converter

A dependency-free CLI tool that converts Markdown to professional DOCX (Open XML) documents with built-in regional style presets and Mermaid diagram rendering.

## Tool Installation

Before using, ensure the **md2docx** binary is available on your PATH.

%s

Verify with: ` + "`md2docx version`" + `

## Usage (Agent Mode — always use ` + "`--json`" + `)

### Convert Markdown to DOCX
` + "```bash" + `
md2docx convert -i <input.md> -o <output.docx> --json
` + "```" + `

### Convert with regional style preset
` + "```bash" + `
md2docx convert -i input.md -o output.docx -s cn-official --json
` + "```" + `

### Convert with Mermaid diagrams rendered as images
` + "```bash" + `
md2docx convert -i input.md -o output.docx --mermaid --mermaid-theme default --json
` + "```" + `

### List available style presets
` + "```bash" + `
md2docx presets --json
` + "```" + `

### Create custom style template
` + "```bash" + `
md2docx template create -o my-style.json -s us-business --json
` + "```" + `

## Style Presets

| Preset       | Region  | Characteristic                             |
|--------------|---------|--------------------------------------------|
| us-business  | US      | Cambria/Calibri, professional blue accent  |
| us-modern    | US      | Segoe UI, minimal dark tones               |
| cn-official  | China   | SimHei/SimSun (公文风格), red accent        |
| cn-modern    | China   | Noto Sans SC, modern Chinese               |
| jp-formal    | Japan   | Yu Mincho/Yu Gothic, business formal       |
| eu-clean     | Europe  | Helvetica/Arial, clean minimalist          |
| kr-standard  | Korea   | Malgun Gothic/Nanum Gothic                 |
| academic     | Global  | Times New Roman, scholarly                 |
| default      | Global  | Aptos Display/Cascadia Mono                |

## Mermaid Rendering

When ` + "`--mermaid`" + ` is set, ` + "`" + "```mermaid" + "`" + ` blocks are rendered as embedded PNG images
(via the public mermaid.ink API). Options:
- ` + "`--mermaid-theme`" + `: default, neutral, dark, forest
- ` + "`--mermaid-server`" + `: custom self-hosted mermaid.ink URL

## JSON Output Format

Success:
` + "```json" + `
{"success": true, "outputPath": "/path/to/output.docx", "bytes": 12345}
` + "```" + `

Error:
` + "```json" + `
{"success": false, "error": "error message"}
` + "```" + `

## Requirements

- No Word, Pandoc, or LibreOffice needed
- Mermaid rendering requires network access (mermaid.ink)
- Self-contained static binary, zero runtime dependencies
`, installNotes)
}

// binaryInstallInstructions returns human-readable install instructions for the current platform.
func binaryInstallInstructions() string {
	var sb strings.Builder
	sb.WriteString("### Install via Go\n\n")
	sb.WriteString("```bash\ngo install github.com/md2docx/cli/cmd/md2docx@latest\n```\n\n")

	// Platform-specific direct download
	sb.WriteString("### Direct Download\n\n")
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goarch == "amd64" {
		goarch = "x86_64"
	}

	sb.WriteString(fmt.Sprintf("Download the latest binary for **%s/%s**:\n", goos, goarch))
	sb.WriteString("[GitHub Releases](https://github.com/md2docx/cli/releases/latest)\n\n")

	switch goos {
	case "windows":
		sb.WriteString("```powershell\n")
		sb.WriteString("# After downloading md2docx-windows-amd64.zip, extract and place in PATH:\n")
		sb.WriteString("Expand-Archive md2docx-windows-amd64.zip -DestinationPath $env:USERPROFILE\\AppData\\Local\\md2docx\n")
		sb.WriteString("```\n")
	case "darwin":
		sb.WriteString("```bash\n")
		sb.WriteString("# Homebrew\nbrew install md2docx/homebrew-tap/md2docx\n\n")
		sb.WriteString("# Or manually:\n")
		sb.WriteString("sudo cp md2docx-darwin-amd64 /usr/local/bin/md2docx\n")
		sb.WriteString("sudo chmod +x /usr/local/bin/md2docx\n")
		sb.WriteString("```\n")
	default: // linux
		sb.WriteString("```bash\n")
		sb.WriteString("# Debian/Ubuntu\ndpkg -i md2docx_*.deb\n\n")
		sb.WriteString("# Or manually:\n")
		sb.WriteString("sudo cp md2docx-linux-amd64 /usr/local/bin/md2docx\n")
		sb.WriteString("sudo chmod +x /usr/local/bin/md2docx\n")
		sb.WriteString("```\n")
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// Installation operations
// ---------------------------------------------------------------------------

// installedResult is the per-agent result of an install operation.
type InstalledResult struct {
	Agent     string
	TargetDir string
	Type      string // "symlink" or "copy"
	Error     error
}

// canonicalSkillPath returns the path to the canonical SKILL.md source in this project.
func canonicalSkillPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		// Fallback: search from CWD upward
		cwd, _ := os.Getwd()
		for dir := cwd; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
			candidate := filepath.Join(dir, "skills", "md2docx", "SKILL.md")
			if _, err := os.Stat(candidate); err == nil {
				return filepath.Abs(candidate)
			}
		}
		return "", fmt.Errorf("cannot find canonical skills/md2docx/SKILL.md")
	}
	// Go up from the binary to find the project root's skills/ directory
	exeDir := filepath.Dir(exe)
	for dir := exeDir; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "skills", "md2docx", "SKILL.md")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Abs(candidate)
		}
	}
	return "", fmt.Errorf("cannot find canonical skills/md2docx/SKILL.md from binary location")
}

// InstallToTarget installs the skill to a specific agent target directory.
// Prefers symlinks on Unix; uses copy on Windows or when symlink fails.
func InstallToTarget(targetDir string, agentName string, preferSymlink bool) (*InstalledResult, error) {
	sourcePath, err := canonicalSkillPath()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", targetDir, err)
	}

	targetFile := filepath.Join(targetDir, "SKILL.md")
	installType := "copy"

	// Try symlink first if preferred
	if preferSymlink && runtime.GOOS != "windows" {
		// Remove existing target (file or symlink)
		os.Remove(targetFile)
		if err := os.Symlink(sourcePath, targetFile); err == nil {
			installType = "symlink"
		}
	}

	if installType == "copy" {
		// Remove existing target
		os.Remove(targetFile)
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("reading skill source: %w", err)
		}
		if err := os.WriteFile(targetFile, data, 0644); err != nil {
			return nil, fmt.Errorf("writing skill to %s: %w", targetFile, err)
		}
	}

	if err := addInstallRecord(targetDir, sourcePath, agentName, installType); err != nil {
		// Non-fatal: installation succeeded, just tracking failed
		fmt.Fprintf(os.Stderr, "warning: failed to record installation: %v\n", err)
	}

	return &InstalledResult{
		Agent:     agentName,
		TargetDir: targetDir,
		Type:      installType,
	}, nil
}

// UninstallFromTarget removes the skill from a target directory.
func UninstallFromTarget(targetDir string) error {
	targetFile := filepath.Join(targetDir, "SKILL.md")
	if err := os.Remove(targetFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing %s: %w", targetFile, err)
	}
	return removeInstallRecord(targetDir)
}

// ---------------------------------------------------------------------------
// High-level install / uninstall / status commands
// ---------------------------------------------------------------------------

// InstallResult is the aggregate result of installing skills for multiple agents.
type InstallResult struct {
	Results  []InstalledResult
	Skipped  []string // agent names that were already installed
}

// InstallSkills installs the skill for all requested agents (or all detected if none specified).
func InstallSkills(agentNames []string, preferSymlink bool) (*InstallResult, error) {
	allTargets := DetectAgents()

	// Determine which agents to install for
	if len(agentNames) == 0 {
		// Install for all detected agents
	} else {
		// Filter to requested agents
		nameSet := make(map[string]bool)
		for _, n := range agentNames {
			nameSet[n] = true
		}
		var filtered []AgentInstallTarget
		for _, t := range allTargets {
			if nameSet[t.Agent.Name] {
				filtered = append(filtered, t)
			}
		}
		// Also add agents that were explicitly requested but not detected — create global dir
		for _, n := range agentNames {
			found := false
			for _, t := range allTargets {
				if t.Agent.Name == n {
					found = true
					break
				}
			}
			if !found {
				// Find the agent definition
				for _, a := range KnownAgents() {
					if a.Name == n {
						// Use the last (global) dir candidate
						globalDir := a.Dirs[len(a.Dirs)-1]
						filtered = append(filtered, AgentInstallTarget{
							Agent:     a,
							TargetDir: filepath.Join(globalDir, a.SkillSubDir),
							Exists:    false,
						})
						break
					}
				}
			}
		}
		allTargets = filtered
	}

	var result InstallResult
	for _, target := range allTargets {
		if target.Exists {
			result.Skipped = append(result.Skipped, target.Agent.Name)
			continue
		}

		installed, err := InstallToTarget(target.TargetDir, target.Agent.Name, preferSymlink)
		if err != nil {
			result.Results = append(result.Results, InstalledResult{
				Agent:     target.Agent.Name,
				TargetDir: target.TargetDir,
				Error:     err,
			})
		} else {
			result.Results = append(result.Results, *installed)
		}
	}

	return &result, nil
}

// InstallToPath installs the skill to an explicit directory path (for custom/unsupported agents).
func InstallToPath(targetDir string, preferSymlink bool) (*InstalledResult, error) {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", targetDir, err)
	}
	return InstallToTarget(targetDir, "custom", preferSymlink)
}

// ListInstallations returns the current installation manifest.
func ListInstallations() ([]installRecord, error) {
	m, err := loadManifest()
	if err != nil {
		return nil, err
	}
	return m.Installed, nil
}

// UninstallAll removes all tracked installations.
func UninstallAll() (int, error) {
	m, err := loadManifest()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, r := range m.Installed {
		if err := UninstallFromTarget(r.Target); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// dirExists checks if a directory exists and contains a SKILL.md file.
func dirExists(path string) bool {
	skillFile := filepath.Join(path, "SKILL.md")
	_, err := os.Stat(skillFile)
	return err == nil
}

// Which returns the path to the md2docx binary, or empty if not found.
func Which() string {
	binName := "md2docx"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	path, err := exec.LookPath(binName)
	if err != nil {
		return ""
	}
	return path
}
