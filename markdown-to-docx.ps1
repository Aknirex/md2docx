[CmdletBinding()]
param(
    [string]$InputPath,
    [string]$OutputPath,
    [string]$TemplatePath,
    [switch]$Interactive
)

Import-Module (Join-Path $PSScriptRoot 'MarkdownDocx.psm1') -Force

if ($Interactive -or (-not $InputPath -and -not $OutputPath)) {
    Start-MarkdownDocxTool -StartDirectory (Get-Location).Path
} else {
    Convert-MarkdownToDocx -InputPath $InputPath -OutputPath $OutputPath -TemplatePath $TemplatePath
}
