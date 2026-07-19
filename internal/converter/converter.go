package converter

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"os"
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

// inlineMarkupRE matches inline bold, italic, and code spans.
var inlineMarkupRE = regexp.MustCompile(`(\*\*.+?\*\*|\x60[^\x60]+\x60|\*[^*\n]+\*|_[^_\n]+_)`)

// inlineCodeRE matches inline code
var inlineCodeRE = regexp.MustCompile(`^\x60([^\x60]+)\x60$`)

// boldRE matches **bold**
var boldRE = regexp.MustCompile(`^\*\*(.+?)\*\*$`)

// italicAsteriskRE matches *italic*
var italicAsteriskRE = regexp.MustCompile(`^\*(.+?)\*$`)

// italicUnderscoreRE matches _italic_
var italicUnderscoreRE = regexp.MustCompile(`^_(.+?)_$`)

// convertInlineMarkdown parses inline markdown (bold, italic, code) into XML runs.
func convertInlineMarkdown(text string, st *StyleTemplate) string {
	var runs strings.Builder
	matches := inlineMarkupRE.FindAllStringIndex(text, -1)
	pos := 0

	for _, m := range matches {
		// Add text before this match as a plain run
		if m[0] > pos {
			runs.WriteString(runXML(text[pos:m[0]], st.BodyFont, st.BodySize, st.TextColor, false, false, false))
		}
		value := text[m[0]:m[1]]

		switch {
		case inlineCodeRE.MatchString(value):
			inner := inlineCodeRE.FindStringSubmatch(value)[1]
			runs.WriteString(runXML(inner, st.CodeFont, st.CodeSize, st.TextColor, false, false, true))
		case boldRE.MatchString(value):
			inner := boldRE.FindStringSubmatch(value)[1]
			runs.WriteString(runXML(inner, st.BodyFont, st.BodySize, st.TextColor, true, false, false))
		case italicAsteriskRE.MatchString(value), italicUnderscoreRE.MatchString(value):
			inner := value[1 : len(value)-1]
			runs.WriteString(runXML(inner, st.BodyFont, st.BodySize, st.TextColor, false, true, false))
		default:
			runs.WriteString(runXML(value, st.BodyFont, st.BodySize, st.TextColor, false, false, false))
		}
		pos = m[1]
	}

	// Remaining text
	if pos < len(text) || runs.Len() == 0 {
		runs.WriteString(runXML(text[pos:], st.BodyFont, st.BodySize, st.TextColor, false, false, false))
	}

	return runs.String()
}

// parseMarkdown converts markdown text into a list of paragraph XML strings.
func parseMarkdown(markdown string, st *StyleTemplate) []string {
	var paragraphs []string
	inCodeBlock := false

	scanner := bufio.NewScanner(strings.NewReader(markdown))
	for scanner.Scan() {
		line := scanner.Text()

		// Handle fenced code blocks
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

	return paragraphs
}

// ConvertMarkdownToBytes converts markdown content to DOCX bytes using the given style.
func ConvertMarkdownToBytes(markdown string, st *StyleTemplate) ([]byte, error) {
	paragraphs := parseMarkdown(markdown, st)

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// [Content_Types].xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
		`<Default Extension="xml" ContentType="application/xml"/>` +
		`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
		`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>` +
		`<Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>` +
		`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>` +
		`</Types>`

	// Relationships
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
		`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>` +
		`</Relationships>`

	// Document rels
	docRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>` +
		`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"/>` +
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

	entries := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", contentTypes},
		{"_rels/.rels", rels},
		{"word/document.xml", documentXML(paragraphs, st.PageMarginInches)},
		{"word/styles.xml", stylesXML(st)},
		{"word/numbering.xml", numbering},
		{"word/_rels/document.xml.rels", docRels},
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

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("closing zip: %w", err)
	}

	return buf.Bytes(), nil
}

// ConvertMarkdownToFile converts a markdown file to a DOCX file using the given style.
func ConvertMarkdownToFile(inputPath, outputPath string, st *StyleTemplate) (*ConversionResult, error) {
	markdownBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("reading markdown file %s: %w", inputPath, err)
	}

	docxBytes, err := ConvertMarkdownToBytes(string(markdownBytes), st)
	if err != nil {
		return nil, fmt.Errorf("converting: %w", err)
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
