# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/md2docx/cli)](https://go.dev)
[![License](https://img.shields.io/github/license/md2docx/cli)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/md2docx/cli)](https://github.com/md2docx/cli/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/md2docx/cli/ci.yml?branch=main)](https://github.com/md2docx/cli/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

Converta Markdown em documentos DOCX profissionais — sem dependencias, nao requer Word nem Pandoc.

Escrito em Go. Distribuido como um unico binario estatico sem dependencias de tempo de execucao.

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## Inicio Rapido

### TUI Interativa (para humanos)

```bash
md2docx
```

Interface de terminal com navegacao por setas:
- Selecionar arquivo Markdown de entrada
- Escolher localizacao e nome do arquivo de saida
- Escolher um preset de estilo integrado (EUA, China, Japao, Europa, Coreia, Academico) ou um template JSON personalizado
- Confirmar e converter

### CLI (para agentes / automacao)

```bash
# Converter com estilo padrao
md2docx convert -i notas.md -o notas.docx --json

# Converter com preset especifico do pais
md2docx convert -i relatorio.md -o relatorio.docx -s cn-official --json

# Listar todos os presets
md2docx presets --json

# Criar um template personalizado a partir de um preset
md2docx template create -o meu-estilo.json -s jp-formal

# Converter com template personalizado
md2docx convert -i doc.md -o doc.docx -s meu-estilo.json --json
```

O flag `--json` produz JSON estruturado para consumo de agentes:
```json
{"success": true, "outputPath": "/caminho/para/saida.docx", "bytes": 12345}
```

## Instalacao

### Via Go

```bash
go install github.com/md2docx/cli/cmd/md2docx@latest
```

### Binarios pre-compilados

Baixar de [GitHub Releases](https://github.com/md2docx/cli/releases) para:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Gerenciadores de pacotes

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## Presets de Estilo Integrados

| Preset       | Regiao  | Fontes                                   |
|-------------|---------|------------------------------------------|
| us-business | EUA     | Cambria / Calibri / Consolas             |
| us-modern   | EUA     | Segoe UI / Cascadia Code                 |
| cn-official | China   | SimHei / SimSun (estilo documento oficial) |
| cn-modern   | China   | Noto Sans SC / Noto Sans Mono SC         |
| jp-formal   | Japao   | Yu Mincho / Yu Gothic                    |
| eu-clean    | Europa  | Helvetica / Arial / Fira Code            |
| kr-standard | Coreia  | Malgun Gothic / Nanum Gothic / D2Coding  |
| academic    | Global  | Times New Roman / Courier New            |
| default     | Global  | Aptos Display / Cascadia Mono            |

## Skill para Agentes

md2docx inclui uma skill integrada para que agentes de IA (Kilo, Claude Code, etc.) possam descobri-la e invoca-la automaticamente.

**Instalar a skill:**

```bash
# Auto-detectar .kilo/skills no projeto atual (ou usar ~/.config/kilo/skills)
md2docx skill install

# Instalar em um caminho explicito
md2docx skill install --path /caminho/.kilo/skills/md2docx
```

Apos a instalacao, agentes que escaneiam `.kilo/skills/` ou `~/.config/kilo/skills/` encontrarao a skill `md2docx` e saberao como invoca-la para conversoes de Markdown para DOCX.

## Templates de Estilo

Templates de estilo personalizados sao arquivos JSON:

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

Criar um a partir de um preset:
```bash
md2docx template create -o meu-estilo.json -s default
```

## Markdown Suportado

- Cabecalhos (h1-h6)
- Paragrafos
- Listas nao ordenadas (`-`, `+`, `*`)
- Listas ordenadas (`1.`, `1)`)
- Citacoes (`>`)
- Blocos de codigo delimitados (`` ``` ``)
- **Negrito**, *italico*, `codigo inline`

## Compilar do Codigo Fonte

```bash
git clone https://github.com/md2docx/cli
cd cli
go mod tidy
make build
```

## Licenca

MIT
