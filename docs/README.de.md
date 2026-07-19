# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/md2docx/cli)](https://go.dev)
[![License](https://img.shields.io/github/license/md2docx/cli)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/md2docx/cli)](https://github.com/md2docx/cli/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/md2docx/cli/ci.yml?branch=main)](https://github.com/md2docx/cli/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

Konvertiert Markdown in professionelle DOCX-Dokumente — ohne Abhaengigkeiten, kein Word oder Pandoc erforderlich.

In Go geschrieben. Als einzelnes statisches Binary ohne Laufzeitabhaengigkeiten verteilt.

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## Schnellstart

### Interaktive TUI (fuer Menschen)

```bash
md2docx
```

Terminal-Oberflaeche mit Pfeiltasten-Navigation:
- Markdown-Eingabedatei auswaehlen
- Ausgabeort und Dateinamen waehlen
- Integrierte Stilvorgabe auswaehlen (USA, China, Japan, Europa, Korea, Akademisch) oder eine benutzerdefinierte JSON-Vorlage
- Bestaetigen und konvertieren

### CLI (fuer Agenten / Automatisierung)

```bash
# Mit Standardstil konvertieren
md2docx convert -i notizen.md -o notizen.docx --json

# Mit laenderspezifischer Vorgabe konvertieren
md2docx convert -i bericht.md -o bericht.docx -s cn-official --json

# Alle Vorgaben auflisten
md2docx presets --json

# Benutzerdefinierte Vorlage aus Vorgabe erstellen
md2docx template create -o mein-stil.json -s jp-formal

# Mit benutzerdefinierter Vorlage konvertieren
md2docx convert -i doc.md -o doc.docx -s mein-stil.json --json
```

Das `--json`-Flag erzeugt strukturiertes JSON fuer die Agentenverarbeitung:
```json
{"success": true, "outputPath": "/pfad/zur/ausgabe.docx", "bytes": 12345}
```

## Installation

### Via Go

```bash
go install github.com/md2docx/cli/cmd/md2docx@latest
```

### Vorkompilierte Binaries

Herunterladen von [GitHub Releases](https://github.com/md2docx/cli/releases) fuer:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Paketmanager

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## Integrierte Stilvorgaben

| Vorgabe     | Region | Schriften                                |
|------------|--------|------------------------------------------|
| us-business | USA    | Cambria / Calibri / Consolas             |
| us-modern   | USA    | Segoe UI / Cascadia Code                 |
| cn-official | China  | SimHei / SimSun (Amtliches Dokument)      |
| cn-modern   | China  | Noto Sans SC / Noto Sans Mono SC         |
| jp-formal   | Japan  | Yu Mincho / Yu Gothic                    |
| eu-clean    | Europa | Helvetica / Arial / Fira Code            |
| kr-standard | Korea  | Malgun Gothic / Nanum Gothic / D2Coding  |
| academic    | Global | Times New Roman / Courier New            |
| default     | Global | Aptos Display / Cascadia Mono            |

## Agenten-Skill

md2docx enthaelt einen integrierten Agenten-Skill, damit KI-Coding-Agenten (Kilo, Claude Code, etc.) ihn automatisch entdecken und aufrufen koennen.

**Skill installieren:**

```bash
# .kilo/skills im aktuellen Projekt automatisch erkennen (oder ~/.config/kilo/skills als Fallback)
md2docx skill install

# An explizitem Pfad installieren
md2docx skill install --path /pfad/.kilo/skills/md2docx
```

Nach der Installation finden Agenten, die `.kilo/skills/` oder `~/.config/kilo/skills/` scannen, den `md2docx`-Skill und wissen, wie sie ihn fuer Markdown-zu-DOCX-Konvertierungen aufrufen.

## Stilvorlagen

Benutzerdefinierte Stilvorlagen sind JSON-Dateien:

```json
{
  "titleFont": "Arial",
  "titleSize": 28,
  "headingFont": "Arial",
  "headingSize": 16,
  "bodyFont": "Times New Roman",
  "bodySize": 12,
  "codeFont": "Courier New",
  "codeSize": 10,
  "textColor": "#1F2937",
  "accentColor": "#2563EB",
  "pageMarginInches": 1.0
}
```

Aus einer Vorgabe erstellen:
```bash
md2docx template create -o mein-stil.json -s default
```

## Unterstuetztes Markdown

- Ueberschriften (h1-h6)
- Absaetze
- Ungeordnete Listen (`-`, `+`, `*`)
- Geordnete Listen (`1.`, `1)`)
- Zitate (`>`)
- Code-Bloecke (`` ``` ``)
- **Fett**, *kursiv*, `Inline-Code`

## Aus dem Quellcode bauen

```bash
git clone https://github.com/md2docx/cli
cd cli
go mod tidy
make build
```

## Lizenz

MIT
