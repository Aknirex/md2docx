# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/Aknirex/md2docx)](https://go.dev)
[![License](https://img.shields.io/github/license/Aknirex/md2docx)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/Aknirex/md2docx)](https://github.com/Aknirex/md2docx/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/Aknirex/md2docx/ci.yml?branch=main)](https://github.com/Aknirex/md2docx/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

将 Markdown 转换为专业 DOCX 文档——无需 Word 或 Pandoc，完全无外部依赖。

使用 Go 编写，以单一静态二进制分发，无需运行时依赖。

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## 快速开始

### 交互式 TUI（面向人类用户）

```bash
md2docx
```

终端界面支持方向键导航：
- 选择 Markdown 输入文件
- 选择输出位置和文件名
- 选择内置风格预设（美国、中国、日本、欧洲、韩国、学术）或自定义 JSON 模板
- 确认并转换

### CLI（面向 AI Agent / 自动化）

```bash
# 使用默认风格转换
md2docx convert -i notes.md -o notes.docx --json

# 使用指定国家预设
md2docx convert -i report.md -o report.docx -s cn-official --json

# 列出所有预设
md2docx presets --json

# 从预设创建自定义模板
md2docx template create -o my-style.json -s jp-formal

# 使用自定义模板转换
md2docx convert -i doc.md -o doc.docx -s my-style.json --json
```

`--json` 标志输出结构化 JSON，适合 Agent 解析：
```json
{"success": true, "outputPath": "/path/to/output.docx", "bytes": 12345}
```

## 安装

### 通过 Go 安装

```bash
go install github.com/Aknirex/md2docx/cmd/md2docx@latest
```

### 预编译二进制

从 [GitHub Releases](https://github.com/Aknirex/md2docx/releases) 下载，支持：
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### 包管理器

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## 内置风格预设

| 预设名称     | 地区   | 字体                                      |
|-------------|--------|-------------------------------------------|
| us-business | 美国   | Cambria / Calibri / Consolas              |
| us-modern   | 美国   | Segoe UI / Cascadia Code                  |
| cn-official | 中国   | 小标宋_GBK / 仿宋_GB2312 / 楷体_GB2312（公文风格） |
| cn-modern   | 中国   | Noto Sans SC / Noto Sans Mono SC          |
| jp-formal   | 日本   | Yu Mincho / Yu Gothic                     |
| eu-clean    | 欧洲   | Helvetica / Arial / Fira Code             |
| kr-standard | 韩国   | Malgun Gothic / Nanum Gothic / D2Coding   |
| academic    | 通用   | Times New Roman / Courier New             |
| default     | 通用   | Aptos Display / Cascadia Mono             |

## Agent Skill

md2docx 包含 SKILL.md，AI 编程助手（Kilo、Claude Code 等）可以自动发现并调用。

**通过 npx skills 安装：**

```bash
npx skills add Aknirex/md2docx
```

安装后，Agent 将知道如何调用 `md2docx convert -i <input> -o <output> --json` 进行 Markdown 到 DOCX 的转换。

## 风格模板

自定义风格模板为 JSON 文件：

```json
{
  "titleFont": "小标宋_GBK",
  "titleSize": 22,
  "headingFont": "楷体_GB2312",
  "headingSize": 16,
  "bodyFont": "仿宋_GB2312",
  "bodySize": 14,
  "codeFont": "楷体_GB2312",
  "codeSize": 12,
  "textColor": "#000000",
  "accentColor": "#C00000",
  "pageMarginInches": 0.71
}
```

从预设创建模板：
```bash
md2docx template create -o my-style.json -s cn-official
```

## 支持的 Markdown

- 标题（h1–h6）
- 段落
- 无序列表（`-`、`+`、`*`）
- 有序列表（`1.`、`1)`）
- 引用块（`>`）
- 围栏代码块（` ``` `）
- **粗体**、*斜体*、`行内代码`

## 从源码构建

```bash
git clone https://github.com/Aknirex/md2docx
cd md2docx
go mod tidy
make build
```

## 许可证

MIT
