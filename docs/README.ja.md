# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/Aknirex/md2docx)](https://go.dev)
[![License](https://img.shields.io/github/license/Aknirex/md2docx)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/Aknirex/md2docx)](https://github.com/Aknirex/md2docx/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/Aknirex/md2docx/ci.yml?branch=main)](https://github.com/Aknirex/md2docx/actions)
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
go install github.com/Aknirex/md2docx/cmd/md2docx@latest
```

### ビルド済みバイナリ

[GitHub Releases](https://github.com/Aknirex/md2docx/releases) からダウンロード：
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
| cn-official  | 中国   | 小标宋_GBK / 仿宋_GB2312 / 楷体_GB2312（公文形式） |
| cn-modern    | 中国   | Noto Sans SC / Noto Sans Mono SC          |
| jp-formal    | 日本   | Yu Mincho / Yu Gothic                     |
| eu-clean     | 欧州   | Helvetica / Arial / Fira Code             |
| kr-standard  | 韓国   | Malgun Gothic / Nanum Gothic / D2Coding   |
| academic     | 汎用   | Times New Roman / Courier New             |
| default      | 汎用   | Aptos Display / Cascadia Mono             |

## エージェントスキル

md2docx には SKILL.md が含まれており、AI コーディングエージェント（Kilo、Claude Code など）が自動的に発見して呼び出せます。

**npx skills でインストール：**

```bash
npx skills add Aknirex/md2docx
```

インストール後、エージェントは `md2docx convert -i <input> -o <output> --json` を使って Markdown から DOCX への変換方法を理解します。

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
git clone https://github.com/Aknirex/md2docx
cd md2docx
go mod tidy
make build
```

## ライセンス

MIT
