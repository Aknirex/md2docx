# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/md2docx/cli)](https://go.dev)
[![License](https://img.shields.io/github/license/md2docx/cli)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/md2docx/cli)](https://github.com/md2docx/cli/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/md2docx/cli/ci.yml?branch=main)](https://github.com/md2docx/cli/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

Convierte Markdown a documentos DOCX profesionales — sin dependencias, no requiere Word ni Pandoc.

Escrito en Go. Distribuido como un unico binario estatico sin dependencias de tiempo de ejecucion.

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## Inicio Rapido

### TUI Interactiva (para humanos)

```bash
md2docx
```

Interfaz de terminal con navegacion por flechas:
- Seleccionar archivo Markdown de entrada
- Elegir ubicacion y nombre del archivo de salida
- Elegir un preset de estilo incorporado (EEUU, China, Japon, Europa, Corea, Academico) o una plantilla JSON personalizada
- Confirmar y convertir

### CLI (para agentes / automatizacion)

```bash
# Convertir con estilo por defecto
md2docx convert -i notas.md -o notas.docx --json

# Convertir con preset especifico del pais
md2docx convert -i informe.md -o informe.docx -s cn-official --json

# Listar todos los presets
md2docx presets --json

# Crear una plantilla personalizada desde un preset
md2docx template create -o mi-estilo.json -s jp-formal

# Convertir con plantilla personalizada
md2docx convert -i doc.md -o doc.docx -s mi-estilo.json --json
```

El flag `--json` produce JSON estructurado para consumo de agentes:
```json
{"success": true, "outputPath": "/ruta/al/salida.docx", "bytes": 12345}
```

## Instalacion

### Via Go

```bash
go install github.com/md2docx/cli/cmd/md2docx@latest
```

### Binarios precompilados

Descargar desde [GitHub Releases](https://github.com/md2docx/cli/releases) para:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Gestores de paquetes

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## Presets de Estilo Incorporados

| Preset       | Region  | Fuentes                                  |
|-------------|---------|------------------------------------------|
| us-business | EEUU    | Cambria / Calibri / Consolas             |
| us-modern   | EEUU    | Segoe UI / Cascadia Code                 |
| cn-official | China   | SimHei / SimSun (estilo documento oficial) |
| cn-modern   | China   | Noto Sans SC / Noto Sans Mono SC         |
| jp-formal   | Japon   | Yu Mincho / Yu Gothic                    |
| eu-clean    | Europa  | Helvetica / Arial / Fira Code            |
| kr-standard | Corea   | Malgun Gothic / Nanum Gothic / D2Coding  |
| academic    | Global  | Times New Roman / Courier New            |
| default     | Global  | Aptos Display / Cascadia Mono            |

## Skill para Agentes

md2docx incluye una skill incorporada para que los agentes de IA (Kilo, Claude Code, etc.) puedan descubrirla e invocarla automaticamente.

**Instalar la skill:**

```bash
# Auto-detectar .kilo/skills en el proyecto actual (o usar ~/.config/kilo/skills)
md2docx skill install

# Instalar en una ruta explicita
md2docx skill install --path /ruta/.kilo/skills/md2docx
```

Tras la instalacion, los agentes que escanean `.kilo/skills/` o `~/.config/kilo/skills/` encontraran la skill `md2docx` y sabran como invocarla para conversiones de Markdown a DOCX.

## Plantillas de Estilo

Las plantillas de estilo personalizadas son archivos JSON:

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

Crear una desde un preset:
```bash
md2docx template create -o mi-estilo.json -s default
```

## Markdown Soportado

- Encabezados (h1-h6)
- Parrafos
- Listas no ordenadas (`-`, `+`, `*`)
- Listas ordenadas (`1.`, `1)`)
- Citas (`>`)
- Bloques de codigo delimitados (`` ``` ``)
- **Negrita**, *cursiva*, `codigo en linea`

## Compilar desde el Codigo Fuente

```bash
git clone https://github.com/md2docx/cli
cd cli
go mod tidy
make build
```

## Licencia

MIT
