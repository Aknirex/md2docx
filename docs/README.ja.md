# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/md2docx/cli)](https://go.dev)
[![License](https://img.shields.io/github/license/md2docx/cli)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/md2docx/cli)](https://github.com/md2docx/cli/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/md2docx/cli/ci.yml?branch=main)](https://github.com/md2docx/cli/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

Markdown をプロフェッショナルな DOCX ドキュメントに変換——Word や Pandoc 不要、外部依存ゼロ。

Go で構築。単一の静的バイナリとして配布され、ランタイム依存はありません。

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## クイックスタート

### インタラクティブ TUI（人間向け）

```bash
md2docx
```

矢印キーで操作するターミナル UI：
- Markdown 入力ファイルを選択
- 出力先とファイル名を選択
- ビルトインスタイルプリセット（米国、中国、日本、欧州、韓国、学術）またはカスタム JSON テンプレートを選択
- 確認して変換

### CLI（AI エージェント / 自動化向け）

```bash
# デフォルトスタイルで変換
md2docx convert -i notes.md -o notes.docx --json

# 国別プリセットで変換
md2docx convert -i report.md -o report.docx -s cn-official --json

# 全プリセットを表示
md2docx presets --json

# プリセットからカスタムテンプレートを作成
md2docx template create -o my-style.json -s jp-formal

# カスタムテンプレートで変換
md2docx convert -i doc.md -o doc.docx -s my-style.json --json
```

`--json` フラグはエージェントが消費しやすい構造化 JSON を出力します：
```json
{"success": true, "outputPath": "/path/to/output.docx", "bytes": 12345}
```

## インストール

### Go 経由

```bash
go install github.com/md2docx/cli/cmd/md2docx@latest
```

### ビルド済みバイナリ

[GitHub Releases](https://github.com/md2docx/cli/releases) からダウンロード：
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### パッケージマネージャ

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## ビルトインスタイルプリセット

| プリセット     | 地域   | フォント                                   |
|--------------|--------|-------------------------------------------|
| us-business  | 米国   | Cambria / Calibri / Consolas              |
| us-modern    | 米国   | Segoe UI / Cascadia Code                  |
| cn-official  | 中国   | SimHei / SimSun（公文形式）                 |
| cn-modern    | 中国   | Noto Sans SC / Noto Sans Mono SC          |
| jp-formal    | 日本   | Yu Mincho / Yu Gothic                     |
| eu-clean     | 欧州   | Helvetica / Arial / Fira Code             |
| kr-standard  | 韓国   | Malgun Gothic / Nanum Gothic / D2Coding   |
| academic     | 汎用   | Times New Roman / Courier New             |
| default      | 汎用   | Aptos Display / Cascadia Mono             |

## エージェントスキル

md2docx にはビルトインのエージェントスキルが含まれており、AI コーディングエージェント（Kilo、Claude Code など）が自動的に発見して呼び出せます。

**スキルのインストール：**

```bash
# 現在のプロジェクトの .kilo/skills を自動検出（なければ ~/.config/kilo/skills にフォールバック）
md2docx skill install

# 明示的なパスにインストール
md2docx skill install --path /path/to/.kilo/skills/md2docx
```

インストール後、`.kilo/skills/` または `~/.config/kilo/skills/` をスキャンするエージェントは `md2docx` スキルを発見し、Markdown から DOCX への変換方法を認識します。

## スタイルテンプレート

カスタムスタイルテンプレートは JSON ファイルです：

```json
{
  "titleFont": "游明朝",
  "titleSize": 22,
  "headingFont": "游ゴシック",
  "headingSize": 16,
  "bodyFont": "游明朝",
  "bodySize": 10.5,
  "codeFont": "UD Digi Kyokasho N-R",
  "codeSize": 10,
  "textColor": "#2D2D2D",
  "accentColor": "#1A478A",
  "pageMarginInches": 0.71
}
```

プリセットから作成：
```bash
md2docx template create -o my-style.json -s jp-formal
```

## サポートする Markdown

- 見出し（h1–h6）
- 段落
- 順序なしリスト（`-`、`+`、`*`）
- 順序付きリスト（`1.`、`1)`）
- 引用（`>`）
- フェンスコードブロック（` ``` `）
- **太字**、*斜体*、`インラインコード`

## ソースからビルド

```bash
git clone https://github.com/md2docx/cli
cd cli
go mod tidy
make build
```

## ライセンス

MIT
