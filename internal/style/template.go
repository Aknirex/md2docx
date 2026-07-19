package style

import (
	"fmt"

	"github.com/md2docx/cli/internal/converter"
)

// ResolveAndConvert resolves a style reference (preset name or template file path)
// and converts a markdown file to DOCX.
func ResolveAndConvert(inputPath, outputPath, styleRef string) (*converter.ConversionResult, error) {
	st, err := LoadStyleTemplate(styleRef)
	if err != nil {
		return nil, fmt.Errorf("resolving style %q: %w", styleRef, err)
	}
	return converter.ConvertMarkdownToFile(inputPath, outputPath, st)
}
