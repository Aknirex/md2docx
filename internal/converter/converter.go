package converter

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// headingRE matches markdown headings (# through ######).
var headingRE = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// unorderedListRE matches unordered list items (-, +, *).
var unorderedListRE = regexp.MustCompile(`^\s*[-+*]\s+(.+)$`)

// orderedListRE matches ordered list items (1., 1), etc.).
var orderedListRE = regexp.MustCompile(`^\s*\d+[.)]\s+(.+)$`)

// blockquoteRE matches blockquotes (> text).
var blockquoteRE = regexp.MustCompile(`^>\s?(.*)$`)

// codeFenceRE matches fenced code block markers.
var codeFenceRE = regexp.MustCompile(`^\s*\x60{3,}`)

// mermaidFenceRE matches ```mermaid (with optional whitespace after backticks).
var mermaidFenceRE = regexp.MustCompile(`^\s*\x60{3,}\s*mermaid\s*$`)

// inlineMarkupRE matches inline bold, italic, and code spans.
var inlineMarkupRE = regexp.MustCompile(`(\*\*.+?\*\*|\x60[^\x60]+\x60|\*[^*\n]+\*|_[^_\n]+_)`)

// mermaidBlock holds a collected mermaid diagram.
type mermaidBlock struct {
	diagram string
	index   int // position in the paragraph list where the placeholder goes
}

// parseResult holds the output of markdown parsing.
type parseResult struct {
	paragraphs    []string
	mermaidBlocks []mermaidBlock
}

// convertInlineMarkdown parses inline markdown (bold, italic, code) into XML runs.
func convertInlineMarkdown(text string, st *StyleTemplate) string {
	var runs strings.Builder
	matches := inlineMarkupRE.FindAllStringIndex(text, -1)
	pos := 0

	for _, m := range matches {
		if m[0] > pos {
			runs.WriteString(runXML(text[pos:m[0]], st.BodyFont, st.BodySize, st.TextColor, false, false, false))
		}
		value := text[m[0]:m[1]]

		switch {
		case len(value) > 1 && value[0] == '`':
			inner := value[1 : len(value)-1]
			runs.WriteString(runXML(inner, st.CodeFont, st.CodeSize, st.TextColor, false, false, true))
		case len(value) > 3 && value[0] == '*' && value[1] == '*':
			inner := value[2 : len(value)-2]
			runs.WriteString(runXML(inner, st.BodyFont, st.BodySize, st.TextColor, true, false, false))
		case len(value) > 1 && value[0] == '*':
			inner := value[1 : len(value)-1]
			runs.WriteString(runXML(inner, st.BodyFont, st.BodySize, st.TextColor, false, true, false))
		case len(value) > 1 && value[0] == '_':
			inner := value[1 : len(value)-1]
			runs.WriteString(runXML(inner, st.BodyFont, st.BodySize, st.TextColor, false, true, false))
		default:
			runs.WriteString(runXML(value, st.BodyFont, st.BodySize, st.TextColor, false, false, false))
		}
		pos = m[1]
	}

	if pos < len(text) || runs.Len() == 0 {
		runs.WriteString(runXML(text[pos:], st.BodyFont, st.BodySize, st.TextColor, false, false, false))
	}

	return runs.String()
}

// parseMarkdown converts markdown text into a list of paragraph XML strings
// and collects mermaid diagram blocks.
func parseMarkdown(markdown string, st *StyleTemplate, enableMermaid bool) *parseResult {
	var paragraphs []string
	var mermaidBlocks []mermaidBlock
	inCodeBlock := false
	inMermaidBlock := false
	var mermaidBuf strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(markdown))
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024) // 1MB max line length
	for scanner.Scan() {
		line := scanner.Text()

		// Inside a mermaid block: collect lines until closing fence
		if inMermaidBlock {
			if codeFenceRE.MatchString(line) {
				// End of mermaid block
				inMermaidBlock = false
				diagram := strings.TrimSpace(mermaidBuf.String())
				if diagram != "" {
					// Insert a placeholder paragraph for this mermaid diagram
					idx := len(paragraphs)
					paragraphs = append(paragraphs, mermaidPlaceholder(idx))
					mermaidBlocks = append(mermaidBlocks, mermaidBlock{
						diagram: diagram,
						index:   idx,
					})
				}
				mermaidBuf.Reset()
				continue
			}
			mermaidBuf.WriteString(line)
			mermaidBuf.WriteString("\n")
			continue
		}

		// Detect opening of a mermaid code block
		if enableMermaid && mermaidFenceRE.MatchString(line) {
			inMermaidBlock = true
			mermaidBuf.Reset()
			continue
		}

		// Handle regular fenced code blocks
		if codeFenceRE.MatchString(line) {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			paragraphs = append(paragraphs,
				paragraphXML(runXML(line, st.CodeFont, st.CodeSize, st.TextColor, false, false, true), "CodeBlock", 0))
			continue
		}

		switch {
		// Headings
		case headingRE.MatchString(line):
			m := headingRE.FindStringSubmatch(line)
			level := len(m[1])
			text := m[2]
			font := st.HeadingFont
			fontSize := st.HeadingSize
			if level == 1 {
				font = st.TitleFont
				fontSize = st.TitleSize
			} else {
				fontSize = st.HeadingSize - float64(level-1)*1.25
				if fontSize < 12 {
					fontSize = 12
				}
			}
			runs := runXML(text, font, fontSize, st.AccentColor, true, false, false)
			paragraphs = append(paragraphs, paragraphXML(runs, fmt.Sprintf("Heading%d", level), 0))

		// Unordered list
		case unorderedListRE.MatchString(line):
			m := unorderedListRE.FindStringSubmatch(line)
			runs := convertInlineMarkdown(m[1], st)
			paragraphs = append(paragraphs, paragraphXML(runs, "", 1))

		// Ordered list
		case orderedListRE.MatchString(line):
			m := orderedListRE.FindStringSubmatch(line)
			runs := convertInlineMarkdown(m[1], st)
			paragraphs = append(paragraphs, paragraphXML(runs, "", 2))

		// Blockquote
		case blockquoteRE.MatchString(line):
			m := blockquoteRE.FindStringSubmatch(line)
			runs := runXML(m[1], st.BodyFont, st.BodySize, st.AccentColor, false, true, false)
			paragraphs = append(paragraphs, paragraphXML(runs, "Quote", 0))

		// Empty line
		case strings.TrimSpace(line) == "":
			paragraphs = append(paragraphs, emptyParagraphXML())

		// Regular paragraph
		default:
			runs := convertInlineMarkdown(line, st)
			paragraphs = append(paragraphs, paragraphXML(runs, "", 0))
		}
	}

	return &parseResult{
		paragraphs:    paragraphs,
		mermaidBlocks: mermaidBlocks,
	}
}

// resolveDefaultStyle returns a default style if nil is passed.
func resolveDefaultStyle(st *StyleTemplate) *StyleTemplate {
	if st == nil {
		return &StyleTemplate{
			TitleFont:        "Aptos Display",
			TitleSize:        28,
			HeadingFont:      "Aptos Display",
			HeadingSize:      18,
			BodyFont:         "Aptos",
			BodySize:         11,
			CodeFont:         "Cascadia Mono",
			CodeSize:         10,
			TextColor:        "#1F2937",
			AccentColor:      "#2563EB",
			PageMarginInches: 0.75,
		}
	}
	return st
}

// ConvertMarkdownToBytes converts markdown content to DOCX bytes.
// Use WithMermaid(r) to enable mermaid diagram rendering.
func ConvertMarkdownToBytes(markdown string, st *StyleTemplate, opts ...ConversionOption) ([]byte, error) {
	cfg := &conversionConfig{Style: resolveDefaultStyle(st)}
	for _, opt := range opts {
		opt(cfg)
	}

	enableMermaid := cfg.Mermaid != nil
	result := parseMarkdown(markdown, cfg.Style, enableMermaid)

	// Render mermaid diagrams if enabled
	var mermaidImages []MermaidImage
	if enableMermaid && len(result.mermaidBlocks) > 0 {
		for i, block := range result.mermaidBlocks {
			pngBytes, wPx, hPx, err := cfg.Mermaid.Render(block.diagram)
			if err != nil {
				// On render failure, fall back to rendering as a code block
				// Replace placeholder with the mermaid source as code
				continue
			}
			imageName := fmt.Sprintf("media/image%d.png", i+1)
			mermaidImages = append(mermaidImages, MermaidImage{
				Index:     block.index,
				ImageName: imageName,
				PNGBytes:  pngBytes,
				WidthEMU:  pixelToEMU(wPx),
				HeightEMU: pixelToEMU(hPx),
			})
		}
	}

	// For mermaid blocks that failed to render, replace placeholders with
	// code paragraphs showing the original mermaid source
	if len(result.mermaidBlocks) > 0 {
		// Build a set of successfully rendered indices
		rendered := make(map[int]bool)
		for _, img := range mermaidImages {
			rendered[img.Index] = true
		}
		// Replace failed placeholders with code paragraphs
		for i, p := range result.paragraphs {
			if idx, ok := parseMermaidPlaceholder(p); ok {
				if !rendered[idx] {
					// Find the original mermaid source
					for _, block := range result.mermaidBlocks {
						if block.index == idx {
							// Render as a code block with language label
							lines := strings.Split(strings.TrimSpace(block.diagram), "\n")
							var codeParas []string
							// Add a label line
							codeParas = append(codeParas,
								paragraphXML(runXML("mermaid", cfg.Style.CodeFont, cfg.Style.CodeSize, cfg.Style.AccentColor, true, false, true), "CodeBlock", 0))
							for _, line := range lines {
								codeParas = append(codeParas,
									paragraphXML(runXML(line, cfg.Style.CodeFont, cfg.Style.CodeSize, cfg.Style.TextColor, false, false, true), "CodeBlock", 0))
							}
							result.paragraphs[i] = strings.Join(codeParas, "")
							break
						}
					}
				}
			}
		}
	}

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// Package-level relationships
	packageRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
		`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>` +
		`</Relationships>`

	// Numbering
	numbering := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:numbering xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:abstractNum w:abstractNumId="0"><w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="bullet"/><w:lvlText w:val="&#x2022;"/></w:lvl></w:abstractNum>` +
		`<w:abstractNum w:abstractNumId="1"><w:lvl w:ilvl="0"><w:start w:val="1"/><w:numFmt w:val="decimal"/><w:lvlText w:val="%1."/></w:lvl></w:abstractNum>` +
		`<w:num w:numId="1"><w:abstractNumId w:val="0"/></w:num>` +
		`<w:num w:numId="2"><w:abstractNumId w:val="1"/></w:num>` +
		`</w:numbering>`

	// Core properties
	coreProps := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/">` +
		`<dc:creator>md2docx</dc:creator>` +
		`</cp:coreProperties>`

	hasImages := len(mermaidImages) > 0

	entries := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", buildContentTypesXML(hasImages)},
		{"_rels/.rels", packageRels},
		{"word/document.xml", documentXML(result.paragraphs, cfg.Style.PageMarginInches, mermaidImages)},
		{"word/styles.xml", stylesXML(cfg.Style)},
		{"word/numbering.xml", numbering},
		{"word/_rels/document.xml.rels", buildDocRelsXML(mermaidImages)},
		{"docProps/core.xml", coreProps},
	}

	for _, entry := range entries {
		fw, err := w.Create(entry.name)
		if err != nil {
			return nil, fmt.Errorf("creating zip entry %s: %w", entry.name, err)
		}
		_, err = fw.Write([]byte(entry.content))
		if err != nil {
			return nil, fmt.Errorf("writing zip entry %s: %w", entry.name, err)
		}
	}

	// Write embedded image files
	for _, img := range mermaidImages {
		fw, err := w.Create("word/" + img.ImageName)
		if err != nil {
			return nil, fmt.Errorf("creating zip entry %s: %w", img.ImageName, err)
		}
		if _, err := fw.Write(img.PNGBytes); err != nil {
			return nil, fmt.Errorf("writing image %s: %w", img.ImageName, err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("closing zip: %w", err)
	}

	return buf.Bytes(), nil
}

// ConvertMarkdownToFile converts a markdown file to a DOCX file using the given style.
func ConvertMarkdownToFile(inputPath, outputPath string, st *StyleTemplate, opts ...ConversionOption) (*ConversionResult, error) {
	markdownBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("reading markdown file %s: %w", inputPath, err)
	}

	docxBytes, err := ConvertMarkdownToBytes(string(markdownBytes), st, opts...)
	if err != nil {
		return nil, fmt.Errorf("converting: %w", err)
	}

	// Ensure output directory exists
	if dir := filepath.Dir(outputPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating output directory: %w", err)
		}
	}

	if err := os.WriteFile(outputPath, docxBytes, 0644); err != nil {
		return nil, fmt.Errorf("writing docx file %s: %w", outputPath, err)
	}

	return &ConversionResult{
		OutputPath: outputPath,
		Bytes:      int64(len(docxBytes)),
	}, nil
}

// ValidateStyle checks that a style template has all required fields.
func ValidateStyle(st *StyleTemplate) error {
	if st.TitleFont == "" {
		return fmt.Errorf("titleFont is required")
	}
	if st.TitleSize <= 0 {
		return fmt.Errorf("titleSize must be positive")
	}
	if st.HeadingFont == "" {
		return fmt.Errorf("headingFont is required")
	}
	if st.HeadingSize <= 0 {
		return fmt.Errorf("headingSize must be positive")
	}
	if st.BodyFont == "" {
		return fmt.Errorf("bodyFont is required")
	}
	if st.BodySize <= 0 {
		return fmt.Errorf("bodySize must be positive")
	}
	if st.CodeFont == "" {
		return fmt.Errorf("codeFont is required")
	}
	if st.CodeSize <= 0 {
		return fmt.Errorf("codeSize must be positive")
	}
	if !strings.HasPrefix(st.TextColor, "#") || len(st.TextColor) != 7 {
		return fmt.Errorf("textColor must be a #RRGGBB value")
	}
	if !strings.HasPrefix(st.AccentColor, "#") || len(st.AccentColor) != 7 {
		return fmt.Errorf("accentColor must be a #RRGGBB value")
	}
	if st.PageMarginInches <= 0 {
		return fmt.Errorf("pageMarginInches must be positive")
	}
	return nil
}
