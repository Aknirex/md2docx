# Markdown to DOCX

A dependency-free PowerShell 7 tool that writes standard DOCX (Open XML) files from Markdown.

## Use

Start the terminal interface:

```powershell
pwsh ./markdown-to-docx.ps1
```

Use Up/Down and Enter to select the Markdown file, a JSON template, and the output folder/name. In every file browser, Backspace goes to the parent folder. Backspace on a selected main-menu field clears it. The template browser also includes **Create a new style template**, which saves a default editable JSON file.

Run non-interactively:

```powershell
pwsh ./markdown-to-docx.ps1 -InputPath .\notes.md -OutputPath .\notes.docx -TemplatePath .\style.json
```

## Templates

Templates are JSON with these required properties: `titleFont`, `titleSize`, `headingFont`, `headingSize`, `bodyFont`, `bodySize`, `codeFont`, `codeSize`, `textColor`, `accentColor`, and `pageMarginInches`. Colors use `#RRGGBB`; sizes and margin must be positive. Create a baseline from the UI or run:

```powershell
Import-Module ./MarkdownDocx.psm1
New-StyleTemplate -Path ./style.json
```

Supported Markdown: headings, paragraphs, unordered and ordered lists, block quotes, fenced code blocks, bold, italics, and inline code.

## Tests

```powershell
Invoke-Pester ./tests
```

The tool produces DOCX directly with .NET ZIP APIs; it does not require Word, Pandoc, LibreOffice, or external modules.
