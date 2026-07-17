Set-StrictMode -Version Latest

function Get-DefaultStyleTemplate {
    [pscustomobject]@{
        titleFont        = 'Aptos Display'
        titleSize        = 28
        headingFont      = 'Aptos Display'
        headingSize      = 18
        bodyFont         = 'Aptos'
        bodySize         = 11
        codeFont         = 'Cascadia Mono'
        codeSize         = 10
        textColor        = '#1F2937'
        accentColor      = '#2563EB'
        pageMarginInches = 0.75
    }
}

function Test-StyleTemplate {
    param([Parameter(Mandatory)]$Style)

    $defaults = Get-DefaultStyleTemplate
    foreach ($property in $defaults.PSObject.Properties.Name) {
        if ($null -eq $Style.PSObject.Properties[$property]) {
            throw "Style template is missing '$property'."
        }
    }
    foreach ($color in 'textColor', 'accentColor') {
        if ([string]$Style.$color -notmatch '^#[0-9A-Fa-f]{6}$') {
            throw "Style template property '$color' must be a #RRGGBB value."
        }
    }
    foreach ($size in 'titleSize', 'headingSize', 'bodySize', 'codeSize') {
        if ([double]$Style.$size -le 0) { throw "Style template property '$size' must be positive." }
    }
    if ([double]$Style.pageMarginInches -le 0) { throw "Style template property 'pageMarginInches' must be positive." }
}

function Read-StyleTemplate {
    [CmdletBinding()]
    param([Parameter(Mandatory)][string]$Path)

    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { throw "Style template not found: $Path" }
    try { $style = Get-Content -LiteralPath $Path -Raw -Encoding UTF8 | ConvertFrom-Json }
    catch { throw "Style template is not valid JSON: $Path" }
    Test-StyleTemplate $style
    return $style
}

function New-StyleTemplate {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)][string]$Path,
        [object]$Style = (Get-DefaultStyleTemplate)
    )

    Test-StyleTemplate $Style
    $directory = Split-Path -Parent $Path
    if ($directory -and -not (Test-Path -LiteralPath $directory -PathType Container)) {
        throw "Template directory does not exist: $directory"
    }
    $Style | ConvertTo-Json | Set-Content -LiteralPath $Path -Encoding UTF8
    return (Get-Item -LiteralPath $Path)
}

function Escape-XmlText {
    param([AllowEmptyString()][string]$Value)
    return [System.Security.SecurityElement]::Escape($Value)
}

function Get-HexColor {
    param([string]$Color)
    return $Color.TrimStart('#').ToUpperInvariant()
}

function New-RunXml {
    param(
        [AllowEmptyString()][string]$Text,
        [string]$Font,
        [double]$Size,
        [string]$Color,
        [switch]$Bold,
        [switch]$Italic,
        [switch]$Code
    )

    $properties = "<w:rPr><w:rFonts w:ascii=`"$Font`" w:hAnsi=`"$Font`"/>"
    $properties += "<w:sz w:val=`"$([int]($Size * 2))`"/><w:color w:val=`"$(Get-HexColor $Color)`"/>"
    if ($Bold) { $properties += '<w:b/>' }
    if ($Italic) { $properties += '<w:i/>' }
    if ($Code) { $properties += '<w:shd w:val="clear" w:fill="F3F4F6"/>' }
    $properties += '</w:rPr>'
    $space = if ($Text -match '^\s|\s$') { ' xml:space="preserve"' } else { '' }
    return "<w:r>$properties<w:t$space>$(Escape-XmlText $Text)</w:t></w:r>"
}

function Convert-InlineMarkdownToRuns {
    param([AllowEmptyString()][string]$Text, [Parameter(Mandatory)]$Style)

    $runs = [System.Collections.Generic.List[string]]::new()
    $pattern = '(\*\*.+?\*\*|`[^`]+`|\*[^*\n]+\*|_[^_\n]+_)'
    $position = 0
    foreach ($match in [regex]::Matches($Text, $pattern)) {
        if ($match.Index -gt $position) {
            $runs.Add((New-RunXml $Text.Substring($position, $match.Index - $position) $Style.bodyFont $Style.bodySize $Style.textColor))
        }
        $value = $match.Value
        if ($value.StartsWith('**')) {
            $runs.Add((New-RunXml $value.Substring(2, $value.Length - 4) $Style.bodyFont $Style.bodySize $Style.textColor -Bold))
        } elseif ($value.StartsWith('`')) {
            $runs.Add((New-RunXml $value.Substring(1, $value.Length - 2) $Style.codeFont $Style.codeSize $Style.textColor -Code))
        } else {
            $runs.Add((New-RunXml $value.Substring(1, $value.Length - 2) $Style.bodyFont $Style.bodySize $Style.textColor -Italic))
        }
        $position = $match.Index + $match.Length
    }
    if ($position -lt $Text.Length -or $runs.Count -eq 0) {
        $runs.Add((New-RunXml $Text.Substring($position) $Style.bodyFont $Style.bodySize $Style.textColor))
    }
    return $runs -join ''
}

function New-ParagraphXml {
    param([string]$Runs, [string]$StyleName, [int]$ListId = 0)

    $properties = ''
    if ($StyleName) { $properties += "<w:pStyle w:val=`"$StyleName`"/>" }
    if ($ListId) { $properties += "<w:numPr><w:ilvl w:val=`"0`"/><w:numId w:val=`"$ListId`"/></w:numPr>" }
    if ($properties) { $properties = "<w:pPr>$properties</w:pPr>" }
    return "<w:p>$properties$Runs</w:p>"
}

function Convert-MarkdownToDocumentXml {
    param([Parameter(Mandatory)][string]$Markdown, [Parameter(Mandatory)]$Style)

    $paragraphs = [System.Collections.Generic.List[string]]::new()
    $inCodeBlock = $false
    foreach ($line in ($Markdown -split "`r?`n")) {
        if ($line -match '^\s*```') { $inCodeBlock = -not $inCodeBlock; continue }
        if ($inCodeBlock) {
            $paragraphs.Add((New-ParagraphXml (New-RunXml $line $Style.codeFont $Style.codeSize $Style.textColor -Code) 'CodeBlock'))
            continue
        }
        if ($line -match '^(#{1,6})\s+(.+)$') {
            $level = $Matches[1].Length
            $font = if ($level -eq 1) { $Style.titleFont } else { $Style.headingFont }
            $fontSize = if ($level -eq 1) { [double]$Style.titleSize } else { [math]::Max(12, [double]$Style.headingSize - (($level - 1) * 1.25)) }
            $runs = New-RunXml $Matches[2] $font $fontSize $Style.accentColor -Bold
            $paragraphs.Add((New-ParagraphXml $runs "Heading$level"))
            continue
        }
        if ($line -match '^\s*[-+*]\s+(.+)$') {
            $paragraphs.Add((New-ParagraphXml (Convert-InlineMarkdownToRuns $Matches[1] $Style) '' 1))
            continue
        }
        if ($line -match '^\s*\d+[.)]\s+(.+)$') {
            $paragraphs.Add((New-ParagraphXml (Convert-InlineMarkdownToRuns $Matches[1] $Style) '' 2))
            continue
        }
        if ($line -match '^>\s?(.*)$') {
            $runs = New-RunXml $Matches[1] $Style.bodyFont $Style.bodySize $Style.accentColor -Italic
            $paragraphs.Add((New-ParagraphXml $runs 'Quote'))
            continue
        }
        if ([string]::IsNullOrWhiteSpace($line)) {
            $paragraphs.Add('<w:p/>')
        } else {
            $paragraphs.Add((New-ParagraphXml (Convert-InlineMarkdownToRuns $line $Style)))
        }
    }
    $margin = [int]([double]$Style.pageMarginInches * 1440)
    return @"
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>$($paragraphs -join '')<w:sectPr><w:pgSz w:w="12240" w:h="15840"/><w:pgMar w:top="$margin" w:right="$margin" w:bottom="$margin" w:left="$margin" w:header="720" w:footer="720" w:gutter="0"/></w:sectPr></w:body></w:document>
"@
}

function Get-StylesXml {
    param([Parameter(Mandatory)]$Style)
    $headingStyles = foreach ($level in 1..6) {
        $size = [math]::Max(12, [double]$Style.headingSize - (($level - 1) * 1.25))
        "<w:style w:type=`"paragraph`" w:styleId=`"Heading$level`"><w:name w:val=`"heading $level`"/><w:basedOn w:val=`"Normal`"/><w:next w:val=`"Normal`"/><w:qFormat/><w:pPr><w:keepNext/></w:pPr><w:rPr><w:rFonts w:ascii=`"$($Style.headingFont)`" w:hAnsi=`"$($Style.headingFont)`"/><w:b/><w:color w:val=`"$(Get-HexColor $Style.accentColor)`"/><w:sz w:val=`"$([int]($size * 2))`"/></w:rPr></w:style>"
    }
    return @"
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:docDefaults><w:rPrDefault><w:rPr><w:rFonts w:ascii="$($Style.bodyFont)" w:hAnsi="$($Style.bodyFont)"/><w:sz w:val="$([int]($Style.bodySize * 2))"/><w:color w:val="$(Get-HexColor $Style.textColor)"/></w:rPr></w:rPrDefault></w:docDefaults><w:style w:type="paragraph" w:default="1" w:styleId="Normal"><w:name w:val="Normal"/></w:style><w:style w:type="paragraph" w:styleId="CodeBlock"><w:name w:val="Code Block"/><w:basedOn w:val="Normal"/><w:rPr><w:rFonts w:ascii="$($Style.codeFont)" w:hAnsi="$($Style.codeFont)"/><w:sz w:val="$([int]($Style.codeSize * 2))"/></w:rPr></w:style><w:style w:type="paragraph" w:styleId="Quote"><w:name w:val="Quote"/><w:basedOn w:val="Normal"/><w:pPr><w:ind w:left="720"/></w:pPr></w:style>$($headingStyles -join '')</w:styles>
"@
}

function Write-ZipEntry {
    param($Archive, [string]$Name, [string]$Content)
    $entry = $Archive.CreateEntry($Name)
    $stream = $entry.Open()
    try {
        $bytes = [System.Text.UTF8Encoding]::new($false).GetBytes($Content)
        $stream.Write($bytes, 0, $bytes.Length)
    } finally { $stream.Dispose() }
}

function Convert-MarkdownToDocx {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)][string]$InputPath,
        [Parameter(Mandatory)][string]$OutputPath,
        [string]$TemplatePath
    )

    if (-not (Test-Path -LiteralPath $InputPath -PathType Leaf)) { throw "Markdown input not found: $InputPath" }
    if ([IO.Path]::GetExtension($OutputPath) -ne '.docx') { throw 'OutputPath must end in .docx.' }
    $style = if ($TemplatePath) { Read-StyleTemplate $TemplatePath } else { Get-DefaultStyleTemplate }
    $outputDirectory = Split-Path -Parent $OutputPath
    if ($outputDirectory -and -not (Test-Path -LiteralPath $outputDirectory -PathType Container)) { throw "Output directory does not exist: $outputDirectory" }
    $markdown = Get-Content -LiteralPath $InputPath -Raw -Encoding UTF8
    $archive = [System.IO.Compression.ZipFile]::Open($OutputPath, [System.IO.Compression.ZipArchiveMode]::Create)
    try {
        Write-ZipEntry $archive '[Content_Types].xml' '<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/><Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/><Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/><Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/></Types>'
        Write-ZipEntry $archive '_rels/.rels' '<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/></Relationships>'
        Write-ZipEntry $archive 'word/document.xml' (Convert-MarkdownToDocumentXml $markdown $style)
        Write-ZipEntry $archive 'word/styles.xml' (Get-StylesXml $style)
        Write-ZipEntry $archive 'word/numbering.xml' '<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:numbering xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:abstractNum w:abstractNumId="0"><w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="bullet"/><w:lvlText w:val="&#x2022;"/></w:lvl></w:abstractNum><w:abstractNum w:abstractNumId="1"><w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="decimal"/><w:lvlText w:val="%1."/></w:lvl></w:abstractNum><w:num w:numId="1"><w:abstractNumId w:val="0"/></w:num><w:num w:numId="2"><w:abstractNumId w:val="1"/></w:num></w:numbering>'
        Write-ZipEntry $archive 'word/_rels/document.xml.rels' '<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"/></Relationships>'
        Write-ZipEntry $archive 'docProps/core.xml' '<?xml version="1.0" encoding="UTF-8" standalone="yes"?><cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:creator>Markdown DOCX Tool</dc:creator></cp:coreProperties>'
    } finally { $archive.Dispose() }
    return (Get-Item -LiteralPath $OutputPath)
}

function Show-Selector {
    param([string]$Title, [string[]]$Items, [int]$Selected)
    Clear-Host
    Write-Host $Title -ForegroundColor Cyan
    Write-Host 'Up/Down: move  Enter: select  Backspace: parent directory  Esc: cancel' -ForegroundColor DarkGray
    for ($index = 0; $index -lt $Items.Count; $index++) {
        $prefix = if ($index -eq $Selected) { '> ' } else { '  ' }
        Write-Host "$prefix$($Items[$index])"
    }
}

function Select-FileInteractive {
    param([string]$Title, [string]$StartDirectory, [string]$Extension, [switch]$CanCreateTemplate)
    $directory = [IO.Path]::GetFullPath($StartDirectory)
    $selected = 0
    while ($true) {
        $items = @('[..]') + @(Get-ChildItem -LiteralPath $directory -Force | Where-Object { $_.PSIsContainer -or $_.Extension -ieq $Extension } | Sort-Object @{ Expression = 'PSIsContainer'; Descending = $true }, Name | ForEach-Object { if ($_.PSIsContainer) { "[D] $($_.Name)" } else { $_.Name } })
        if ($CanCreateTemplate) { $items += '[Create a new style template]' }
        $selected = [Math]::Min($selected, $items.Count - 1)
        Show-Selector "$Title`n$directory" $items $selected
        $key = [Console]::ReadKey($true).Key
        if ($key -eq 'Escape') { return $null }
        if ($key -eq 'UpArrow') { $selected = ($selected - 1 + $items.Count) % $items.Count; continue }
        if ($key -eq 'DownArrow') { $selected = ($selected + 1) % $items.Count; continue }
        if ($key -eq 'Backspace' -or ($key -eq 'Enter' -and $selected -eq 0)) {
            $parent = Split-Path -Parent $directory
            if ($parent) { $directory = $parent; $selected = 0 }
            continue
        }
        if ($key -eq 'Enter' -and $CanCreateTemplate -and $items[$selected] -eq '[Create a new style template]') {
            $name = Read-Host 'Template filename (without .json is allowed)'
            if ($name) {
                if (-not $name.EndsWith('.json')) { $name += '.json' }
                $path = Join-Path $directory $name
                New-StyleTemplate -Path $path | Out-Null
                return $path
            }
            continue
        }
        if ($key -eq 'Enter' -and $items[$selected].StartsWith('[D] ')) { $directory = Join-Path $directory $items[$selected].Substring(4); $selected = 0; continue }
        if ($key -eq 'Enter' -and $selected -gt 0) { return (Join-Path $directory $items[$selected]) }
    }
}

function Select-OutputPathInteractive {
    param([string]$StartDirectory)
    $directory = [IO.Path]::GetFullPath($StartDirectory)
    $selected = 0
    while ($true) {
        $items = @('[Use this folder]', '[..]') + @(Get-ChildItem -LiteralPath $directory -Directory -Force | Sort-Object Name | ForEach-Object { "[D] $($_.Name)" })
        $selected = [Math]::Min($selected, $items.Count - 1)
        Show-Selector "Choose DOCX output folder`n$directory" $items $selected
        $key = [Console]::ReadKey($true).Key
        if ($key -eq 'Escape') { return $null }
        if ($key -eq 'UpArrow') { $selected = ($selected - 1 + $items.Count) % $items.Count; continue }
        if ($key -eq 'DownArrow') { $selected = ($selected + 1) % $items.Count; continue }
        if ($key -eq 'Backspace' -or ($key -eq 'Enter' -and $selected -eq 1)) { $parent = Split-Path -Parent $directory; if ($parent) { $directory = $parent; $selected = 0 }; continue }
        if ($key -eq 'Enter' -and $selected -eq 0) { $name = Read-Host 'DOCX filename'; if ($name) { if (-not $name.EndsWith('.docx')) { $name += '.docx' }; return (Join-Path $directory $name) }; continue }
        if ($key -eq 'Enter') { $directory = Join-Path $directory $items[$selected].Substring(4); $selected = 0 }
    }
}

function Start-MarkdownDocxTool {
    [CmdletBinding()]
    param([string]$StartDirectory = (Get-Location).Path)
    $inputPath = $null; $templatePath = $null; $outputPath = $null; $selected = 0
    while ($true) {
        Clear-Host
        Write-Host 'Markdown to DOCX' -ForegroundColor Cyan
        Write-Host 'Up/Down: move  Enter: choose/run  Backspace: clear selected value  Esc: exit' -ForegroundColor DarkGray
        $items = @("Markdown input: $(if ($inputPath) { $inputPath } else { '<select>' })", "Style template: $(if ($templatePath) { $templatePath } else { '<default or select>' })", "DOCX output: $(if ($outputPath) { $outputPath } else { '<select>' })", 'Convert')
        for ($index = 0; $index -lt $items.Count; $index++) { Write-Host "$(if ($index -eq $selected) { '> ' } else { '  ' })$($items[$index])" }
        $key = [Console]::ReadKey($true).Key
        if ($key -eq 'Escape') { return }
        if ($key -eq 'UpArrow') { $selected = ($selected - 1 + $items.Count) % $items.Count; continue }
        if ($key -eq 'DownArrow') { $selected = ($selected + 1) % $items.Count; continue }
        if ($key -eq 'Backspace') { if ($selected -eq 0) { $inputPath = $null } elseif ($selected -eq 1) { $templatePath = $null } elseif ($selected -eq 2) { $outputPath = $null }; continue }
        if ($key -ne 'Enter') { continue }
        if ($selected -eq 0) { $inputPath = Select-FileInteractive 'Choose Markdown input' $StartDirectory '.md' }
        elseif ($selected -eq 1) { $templatePath = Select-FileInteractive 'Choose style template' $StartDirectory '.json' -CanCreateTemplate }
        elseif ($selected -eq 2) { $outputPath = Select-OutputPathInteractive $StartDirectory }
        elseif (-not $inputPath -or -not $outputPath) { Write-Host 'Select both Markdown input and DOCX output first.' -ForegroundColor Yellow; Start-Sleep -Seconds 2 }
        else { Convert-MarkdownToDocx -InputPath $inputPath -OutputPath $outputPath -TemplatePath $templatePath | Out-Null; Write-Host "Created $outputPath" -ForegroundColor Green; [Console]::ReadKey($true) | Out-Null }
    }
}

Export-ModuleMember -Function Convert-MarkdownToDocx, New-StyleTemplate, Read-StyleTemplate, Start-MarkdownDocxTool
