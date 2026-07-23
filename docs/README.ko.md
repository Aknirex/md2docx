# md2docx

[![Go Version](https://img.shields.io/github/go-mod/go-version/Aknirex/md2docx)](https://go.dev)
[![License](https://img.shields.io/github/license/Aknirex/md2docx)](../LICENSE)
[![Release](https://img.shields.io/github/v/release/Aknirex/md2docx)](https://github.com/Aknirex/md2docx/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/Aknirex/md2docx/ci.yml?branch=main)](https://github.com/Aknirex/md2docx/actions)
[![Platforms](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-blue)]()

Markdown을 전문적인 DOCX 문서로 변환 — Word나 Pandoc 불필요, 외부 의존성 제로.

Go로 작성되었으며 단일 정적 바이너리로 배포됩니다. 런타임 의존성이 없습니다.

[English](../README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Español](./README.es.md) | [Português](./README.pt-BR.md) | [Deutsch](./README.de.md) | [Français](./README.fr.md)

## 빠른 시작

### 대화형 TUI (사용자용)

```bash
md2docx
```

화살표 키로 탐색하는 터미널 UI:
- Markdown 입력 파일 선택
- 출력 위치 및 파일명 선택
- 내장 스타일 프리셋(미국, 중국, 일본, 유럽, 한국, 학술) 또는 사용자 정의 JSON 템플릿 선택
- 확인 후 변환

### CLI (AI 에이전트 / 자동화용)

```bash
# 기본 스타일로 변환
md2docx convert -i notes.md -o notes.docx --json

# 국가별 프리셋으로 변환
md2docx convert -i report.md -o report.docx -s cn-official --json

# 모든 프리셋 목록
md2docx presets --json

# 프리셋에서 사용자 정의 템플릿 생성
md2docx template create -o my-style.json -s jp-formal

# 사용자 정의 템플릿으로 변환
md2docx convert -i doc.md -o doc.docx -s my-style.json --json
```

`--json` 플래그는 에이전트 소비에 적합한 구조화된 JSON을 출력합니다:
```json
{"success": true, "outputPath": "/path/to/output.docx", "bytes": 12345}
```

## 설치

### Go를 통한 설치

```bash
go install github.com/Aknirex/md2docx/cmd/md2docx@latest
```

### 사전 빌드된 바이너리

[GitHub Releases](https://github.com/Aknirex/md2docx/releases)에서 다운로드:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### 패키지 관리자

```bash
# Homebrew
brew install md2docx/homebrew-tap/md2docx

# Debian/Ubuntu
dpkg -i md2docx_*.deb

# RPM
rpm -i md2docx_*.rpm
```

## 내장 스타일 프리셋

| 프리셋        | 지역   | 글꼴                                      |
|-------------|--------|-------------------------------------------|
| us-business | 미국   | Cambria / Calibri / Consolas              |
| us-modern   | 미국   | Segoe UI / Cascadia Code                  |
| cn-official | 중국   | 小标宋_GBK / 仿宋_GB2312 / 楷体_GB2312 (공문서 스타일) |
| cn-modern   | 중국   | Noto Sans SC / Noto Sans Mono SC          |
| jp-formal   | 일본   | Yu Mincho / Yu Gothic                     |
| eu-clean    | 유럽   | Helvetica / Arial / Fira Code             |
| kr-standard | 한국   | Malgun Gothic / Nanum Gothic / D2Coding   |
| academic    | 범용   | Times New Roman / Courier New             |
| default     | 범용   | Aptos Display / Cascadia Mono             |

## 에이전트 스킬

md2docx에는 SKILL.md가 포함되어 있어 AI 코딩 에이전트(Kilo, Claude Code 등)가 자동으로 발견하고 호출할 수 있습니다.

**npx skills로 설치:**

```bash
npx skills add Aknirex/md2docx
```

설치 후 에이전트는 `md2docx convert -i <input> -o <output> --json`을 사용하여 Markdown을 DOCX로 변환하는 방법을 알게 됩니다.

## 스타일 템플릿

사용자 정의 스타일 템플릿은 JSON 파일입니다:

```json
{
  "titleFont": "맑은 고딕",
  "titleSize": 24,
  "headingFont": "맑은 고딕",
  "headingSize": 16,
  "bodyFont": "나눔고딕",
  "bodySize": 11,
  "codeFont": "D2Coding",
  "codeSize": 10,
  "textColor": "#1C1C1C",
  "accentColor": "#1E6B4E",
  "pageMarginInches": 0.8
}
```

프리셋에서 생성:
```bash
md2docx template create -o my-style.json -s kr-standard
```

## 지원하는 Markdown

- 제목 (h1–h6)
- 단락
- 순서 없는 목록 (`-`, `+`, `*`)
- 순서 있는 목록 (`1.`, `1)`)
- 인용구 (`>`)
- 코드 블록 (` ``` `)
- **굵게**, *기울임*, `인라인 코드`

## 소스에서 빌드

```bash
git clone https://github.com/Aknirex/md2docx
cd md2docx
go mod tidy
make build
```

## 라이선스

MIT
