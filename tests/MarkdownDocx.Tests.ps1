$modulePath = Join-Path $PSScriptRoot '..\MarkdownDocx.psm1'
Import-Module $modulePath -Force

Describe 'Markdown DOCX conversion' {
    $testRoot = Join-Path ([IO.Path]::GetTempPath()) "markdown-docx-tests-$([guid]::NewGuid())"
    New-Item -ItemType Directory -Path $testRoot | Out-Null

    try {
        It 'creates a valid DOCX package with converted content' {
            $input = Join-Path $testRoot 'input.md'
            $output = Join-Path $testRoot 'output.docx'
            Set-Content -LiteralPath $input -Encoding UTF8 -Value "# Heading`n`nText with **bold** and `inline`.`n`n- Item"

            Convert-MarkdownToDocx -InputPath $input -OutputPath $output | Out-Null

            Test-Path -LiteralPath $output | Should Be $true
            $archive = [System.IO.Compression.ZipFile]::OpenRead($output)
            try {
                $entryNames = @($archive.Entries.FullName)
                ($entryNames -contains 'word/document.xml') | Should Be $true
                ($entryNames -contains 'word/styles.xml') | Should Be $true
                $document = [IO.StreamReader]::new(($archive.GetEntry('word/document.xml')).Open()).ReadToEnd()
                $document | Should Match 'Heading'
                $document | Should Match 'bold'
            } finally { $archive.Dispose() }
        }

        It 'uses a supplied JSON style template' {
            $template = Join-Path $testRoot 'style.json'
            $input = Join-Path $testRoot 'template-input.md'
            $output = Join-Path $testRoot 'template-output.docx'
            $style = [pscustomobject]@{ titleFont='Impact'; titleSize=31; headingFont='Arial'; headingSize=20; bodyFont='Arial'; bodySize=12; codeFont='Consolas'; codeSize=10; textColor='#112233'; accentColor='#AA5500'; pageMarginInches=1 }
            New-StyleTemplate -Path $template -Style $style | Out-Null
            Set-Content -LiteralPath $input -Encoding UTF8 -Value '# Styled heading'

            Convert-MarkdownToDocx -InputPath $input -OutputPath $output -TemplatePath $template | Out-Null

            $archive = [System.IO.Compression.ZipFile]::OpenRead($output)
            try {
                $styles = [IO.StreamReader]::new(($archive.GetEntry('word/styles.xml')).Open()).ReadToEnd()
                $document = [IO.StreamReader]::new(($archive.GetEntry('word/document.xml')).Open()).ReadToEnd()
                $styles | Should Match 'AA5500'
                $styles | Should Match 'Arial'
                $document | Should Match 'w:ascii="Impact"'
                $document | Should Match 'w:sz w:val="62"'
            } finally { $archive.Dispose() }
        }

        It 'rejects an incomplete template' {
            $template = Join-Path $testRoot 'invalid.json'
            Set-Content -LiteralPath $template -Encoding UTF8 -Value '{"bodyFont":"Arial"}'
            $threw = $false
            try { Read-StyleTemplate -Path $template | Out-Null } catch { $threw = $true }
            $threw | Should Be $true
        }
    } finally {
        Remove-Item -LiteralPath $testRoot -Recurse -Force -ErrorAction SilentlyContinue
    }
}
