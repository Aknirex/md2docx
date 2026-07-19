# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/md2docx/cli)](https://go.dev)
[![License](https://img.shields.io/github/license/md2docx/cli)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/md2docx/cli)](https://github.com/md2docx/cli/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/md2docx/cli/ci.yml?branch=main)](https://github.com/md2docx/cli/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

Convertit Markdown en documents DOCX professionnels — sans dependance externe, ni Word ni Pandoc necessaires.

Ecrit en Go. Distribue sous forme d'un seul binaire statique sans dependance d'execution.

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## Demarrage Rapide

### TUI Interactive (pour les humains)

```bash
md2docx
```

Interface terminal avec navigation par fleches :
- Selectionner le fichier Markdown d'entree
- Choisir l'emplacement et le nom du fichier de sortie
- Choisir un preset de style integre (USA, Chine, Japon, Europe, Coree, Academique) ou un modele JSON personnalise
- Confirmer et convertir

### CLI (pour les agents / l'automatisation)

```bash
# Convertir avec le style par defaut
md2docx convert -i notes.md -o notes.docx --json

# Convertir avec un preset specifique au pays
md2docx convert -i rapport.md -o rapport.docx -s cn-official --json

# Lister tous les presets
md2docx presets --json

# Creer un modele personnalise a partir d'un preset
md2docx template create -o mon-style.json -s jp-formal

# Convertir avec un modele personnalise
md2docx convert -i doc.md -o doc.docx -s mon-style.json --json
```

Le flag `--json` produit du JSON structure pour la consommation par les agents :
```json
{"success": true, "outputPath": "/chemin/vers/sortie.docx", "bytes": 12345}
```

## Installation

### Via Go

```bash
go install github.com/md2docx/cli/cmd/md2docx@latest
```

### Binaires pre-compiles

Telecharger depuis [GitHub Releases](https://github.com/md2docx/cli/releases) pour :
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Gestionnaires de paquets

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## Presets de Style Integres

| Preset      | Region  | Polices                                  |
|-------------|---------|------------------------------------------|
| us-business | USA     | Cambria / Calibri / Consolas             |
| us-modern   | USA     | Segoe UI / Cascadia Code                 |
| cn-official | Chine   | SimHei / SimSun (style document officiel) |
| cn-modern   | Chine   | Noto Sans SC / Noto Sans Mono SC         |
| jp-formal   | Japon   | Yu Mincho / Yu Gothic                    |
| eu-clean    | Europe  | Helvetica / Arial / Fira Code            |
| kr-standard | Coree   | Malgun Gothic / Nanum Gothic / D2Coding  |
| academic    | Global  | Times New Roman / Courier New            |
| default     | Global  | Aptos Display / Cascadia Mono            |

## Skill pour Agents

md2docx inclut une skill integree pour que les agents d'IA (Kilo, Claude Code, etc.) puissent la decouvrir et l'invoquer automatiquement.

**Installer la skill :**

```bash
# Auto-detecter .kilo/skills dans le projet actuel (ou utiliser ~/.config/kilo/skills)
md2docx skill install

# Installer a un chemin explicite
md2docx skill install --path /chemin/.kilo/skills/md2docx
```

Apres installation, les agents qui scannent `.kilo/skills/` ou `~/.config/kilo/skills/` trouveront la skill `md2docx` et sauront comment l'invoquer pour les conversions Markdown vers DOCX.

## Modeles de Style

Les modeles de style personnalises sont des fichiers JSON :

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

Creer un modele depuis un preset :
```bash
md2docx template create -o mon-style.json -s default
```

## Markdown Supporte

- Titres (h1-h6)
- Paragraphes
- Listes non ordonnees (`-`, `+`, `*`)
- Listes ordonnees (`1.`, `1)`)
- Citations (`>`)
- Blocs de code delimites (`` ``` ``)
- **Gras**, *italique*, `code en ligne`

## Compiler depuis les Sources

```bash
git clone https://github.com/md2docx/cli
cd cli
go mod tidy
make build
```

## Licence

MIT
