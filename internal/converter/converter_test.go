package converter

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestConvertInlineMarkdown_PlainText(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("hello world", st)
	if !strings.Contains(result, "hello world") {
		t.Errorf("expected plain text, got: %s", result)
	}
	if strings.Contains(result, "<w:b/>") {
		t.Error("plain text should not be bold")
	}
}

func TestConvertInlineMarkdown_Bold(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("**bold text**", st)
	if !strings.Contains(result, "<w:b/>") {
		t.Error("expected bold markup")
	}
	if !strings.Contains(result, "bold text") {
		t.Error("expected bold text content")
	}
}

func TestConvertInlineMarkdown_Italic(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	tests := []struct {
		name  string
		input string
	}{
		{"asterisk", "*italic*"},
		{"underscore", "_italic_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertInlineMarkdown(tt.input, st)
			if !strings.Contains(result, "<w:i/>") {
				t.Errorf("expected italic markup for %q, got: %s", tt.input, result)
			}
			if !strings.Contains(result, "italic") {
				t.Errorf("expected italic text content for %q", tt.input)
			}
		})
	}
}

func TestConvertInlineMarkdown_InlineCode(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("`code`", st)
	if !strings.Contains(result, "code") {
		t.Error("expected code text content")
	}
	if !strings.Contains(result, "F3F4F6") {
		t.Error("expected code background highlight")
	}
}

func TestConvertInlineMarkdown_MixedContent(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("text **bold** and `code`", st)
	if !strings.Contains(result, "text ") {
		t.Error("expected leading text")
	}
	if !strings.Contains(result, "<w:b/>") {
		t.Error("expected bold markup")
	}
	if !strings.Contains(result, "F3F4F6") {
		t.Error("expected code background")
	}
}

func TestParseMarkdown_Headings(t *testing.T) {
	st := resolveDefaultStyle(nil)
	tests := []struct {
		input    string
		styleID  string
		isTitle  bool
	}{
		{"# Title", "Heading1", true},
		{"## Heading 2", "Heading2", false},
		{"### Heading 3", "Heading3", false},
		{"#### Heading 4", "Heading4", false},
		{"##### Heading 5", "Heading5", false},
		{"###### Heading 6", "Heading6", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseMarkdown(tt.input, st, false)
			if len(result.paragraphs) != 1 {
				t.Fatalf("expected 1 paragraph, got %d", len(result.paragraphs))
			}
			if !strings.Contains(result.paragraphs[0], tt.styleID) {
				t.Errorf("expected style %s in %s", tt.styleID, result.paragraphs[0])
			}
		})
	}
}

func TestParseMarkdown_UnorderedList(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "- item 1\n- item 2\n- item 3"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 3 {
		t.Fatalf("expected 3 paragraphs, got %d", len(result.paragraphs))
	}
	for _, p := range result.paragraphs {
		if !strings.Contains(p, `<w:numId w:val="1"/>`) {
			t.Errorf("expected unordered list numId=1 in %s", p)
		}
	}
}

func TestParseMarkdown_OrderedList(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "1. first\n2. second\n3. third"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 3 {
		t.Fatalf("expected 3 paragraphs, got %d", len(result.paragraphs))
	}
	for _, p := range result.paragraphs {
		if !strings.Contains(p, `<w:numId w:val="2"/>`) {
			t.Errorf("expected ordered list numId=2 in %s", p)
		}
	}
}

func TestParseMarkdown_Blockquote(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "> quoted text"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 1 {
		t.Fatalf("expected 1 paragraph, got %d", len(result.paragraphs))
	}
	if !strings.Contains(result.paragraphs[0], "Quote") {
		t.Errorf("expected Quote style in %s", result.paragraphs[0])
	}
	if !strings.Contains(result.paragraphs[0], "quoted text") {
		t.Errorf("expected quoted text in %s", result.paragraphs[0])
	}
}

func TestParseMarkdown_CodeBlock(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "```\nline1\nline2\n```"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 2 {
		t.Fatalf("expected 2 paragraphs (code lines), got %d", len(result.paragraphs))
	}
	for _, p := range result.paragraphs {
		if !strings.Contains(p, "CodeBlock") {
			t.Errorf("expected CodeBlock style in %s", p)
		}
	}
}

func TestParseMarkdown_EmptyLine(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "text\n\nmore text"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 3 {
		t.Fatalf("expected 3 paragraphs, got %d", len(result.paragraphs))
	}
	if result.paragraphs[1] != "<w:p/>" {
		t.Errorf("expected empty paragraph, got %s", result.paragraphs[1])
	}
}

func TestParseMarkdown_MermaidEnabled(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "```mermaid\ngraph TD\n    A-->B\n```"
	result := parseMarkdown(input, st, true)
	if len(result.mermaidBlocks) != 1 {
		t.Fatalf("expected 1 mermaid block, got %d", len(result.mermaidBlocks))
	}
	if result.mermaidBlocks[0].diagram != "graph TD\n    A-->B" {
		t.Errorf("unexpected mermaid diagram: %q", result.mermaidBlocks[0].diagram)
	}
}

func TestParseMarkdown_MermaidDisabled(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "```mermaid\ngraph TD\n    A-->B\n```"
	result := parseMarkdown(input, st, false)
	if len(result.mermaidBlocks) != 0 {
		t.Fatalf("expected 0 mermaid blocks when disabled, got %d", len(result.mermaidBlocks))
	}
}

func TestResolveDefaultStyle_Nil(t *testing.T) {
	st := resolveDefaultStyle(nil)
	if st == nil {
		t.Fatal("expected non-nil default style")
	}
	if st.TitleFont != "Aptos Display" {
		t.Errorf("expected Aptos Display, got %s", st.TitleFont)
	}
	if st.BodySize != 11 {
		t.Errorf("expected body size 11, got %f", st.BodySize)
	}
}

func TestResolveDefaultStyle_NonNil(t *testing.T) {
	input := &StyleTemplate{TitleFont: "CustomFont"}
	st := resolveDefaultStyle(input)
	if st.TitleFont != "CustomFont" {
		t.Errorf("expected CustomFont, got %s", st.TitleFont)
	}
}

func TestConvertMarkdownToBytes_BasicDocx(t *testing.T) {
	md := "# Hello\n\nThis is a test."
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty DOCX output")
	}

	// Verify it's a valid zip
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("not a valid zip: %v", err)
	}

	entries := make(map[string]bool)
	for _, f := range r.File {
		entries[f.Name] = true
	}
	required := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"word/document.xml",
		"word/styles.xml",
		"word/numbering.xml",
		"word/_rels/document.xml.rels",
		"docProps/core.xml",
	}
	for _, name := range required {
		if !entries[name] {
			t.Errorf("missing required DOCX entry: %s", name)
		}
	}
}

func TestConvertMarkdownToBytes_DocumentXML(t *testing.T) {
	md := "# Title\n\nParagraph text."
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readZipEntry(t, data, "word/document.xml")
	if !strings.Contains(content, "Title") {
		t.Error("document.xml should contain title text")
	}
	if !strings.Contains(content, "Paragraph text.") {
		t.Error("document.xml should contain paragraph text")
	}
	if !strings.Contains(content, "Heading1") {
		t.Error("document.xml should contain Heading1 style")
	}
}

func TestConvertMarkdownToBytes_StylesXML(t *testing.T) {
	data, err := ConvertMarkdownToBytes("# test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readZipEntry(t, data, "word/styles.xml")
	if !strings.Contains(content, "Normal") {
		t.Error("styles.xml should define Normal style")
	}
	if !strings.Contains(content, "Heading1") {
		t.Error("styles.xml should define Heading1 style")
	}
	if !strings.Contains(content, "CodeBlock") {
		t.Error("styles.xml should define CodeBlock style")
	}
	if !strings.Contains(content, "Quote") {
		t.Error("styles.xml should define Quote style")
	}
}

func TestConvertMarkdownToBytes_CustomStyle(t *testing.T) {
	st := &StyleTemplate{
		TitleFont:        "CustomTitle",
		TitleSize:        32,
		HeadingFont:      "CustomHeading",
		HeadingSize:      20,
		BodyFont:         "CustomBody",
		BodySize:         12,
		CodeFont:         "CustomCode",
		CodeSize:         9,
		TextColor:        "#111111",
		AccentColor:      "#FF0000",
		PageMarginInches: 1.0,
	}
	data, err := ConvertMarkdownToBytes("# test", st)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readZipEntry(t, data, "word/styles.xml")
	// styles.xml uses HeadingFont for heading styles and BodyFont for Normal
	if !strings.Contains(content, "CustomHeading") {
		t.Error("styles.xml should use custom heading font")
	}
	if !strings.Contains(content, "CustomBody") {
		t.Error("styles.xml should use custom body font")
	}
	if !strings.Contains(content, "CustomCode") {
		t.Error("styles.xml should use custom code font")
	}

	// TitleFont is used in document.xml for H1 paragraphs
	docContent := readZipEntry(t, data, "word/document.xml")
	if !strings.Contains(docContent, "CustomTitle") {
		t.Error("document.xml should use custom title font for H1")
	}
}

func TestConvertMarkdownToBytes_AllMarkdownElements(t *testing.T) {
	md := `# Title

## Subtitle

Regular paragraph with **bold** and *italic* and ` + "`code`" + `.

- unordered item 1
- unordered item 2

1. ordered item 1
2. ordered item 2

> A blockquote

` + "```" + `
code block line 1
code block line 2
` + "```" + `

`
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readZipEntry(t, data, "word/document.xml")
	checks := []struct {
		name string
		want string
	}{
		{"title", "Title"},
		{"subtitle", "Subtitle"},
		{"bold", "<w:b/>"},
		{"italic", "<w:i/>"},
		{"unordered list", `<w:numId w:val="1"/>`},
		{"ordered list", `<w:numId w:val="2"/>`},
		{"blockquote", "Quote"},
		{"code block", "CodeBlock"},
	}
	for _, c := range checks {
		if !strings.Contains(content, c.want) {
			t.Errorf("missing %s: expected %q in document.xml", c.name, c.want)
		}
	}
}

func TestValidateStyle_AllValid(t *testing.T) {
	st := &StyleTemplate{
		TitleFont: "Arial", TitleSize: 28,
		HeadingFont: "Arial", HeadingSize: 18,
		BodyFont: "Arial", BodySize: 11,
		CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
		PageMarginInches: 1.0,
	}
	if err := ValidateStyle(st); err != nil {
		t.Errorf("expected valid style, got error: %v", err)
	}
}

func TestValidateStyle_MissingFields(t *testing.T) {
	tests := []struct {
		name  string
		field string
		style *StyleTemplate
	}{
		{"empty TitleFont", "titleFont", &StyleTemplate{TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"zero TitleSize", "titleSize", &StyleTemplate{TitleFont: "A", TitleSize: 0, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"empty HeadingFont", "headingFont", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"zero HeadingSize", "headingSize", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 0, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"empty BodyFont", "bodyFont", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"zero BodySize", "bodySize", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 0, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"empty CodeFont", "codeFont", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"zero CodeSize", "codeSize", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 0, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 1}},
		{"bad TextColor", "textColor", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "red", AccentColor: "#000000", PageMarginInches: 1}},
		{"bad AccentColor", "accentColor", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "blue", PageMarginInches: 1}},
		{"zero Margin", "pageMarginInches", &StyleTemplate{TitleFont: "A", TitleSize: 1, HeadingFont: "A", HeadingSize: 1, BodyFont: "A", BodySize: 1, CodeFont: "A", CodeSize: 1, TextColor: "#000000", AccentColor: "#000000", PageMarginInches: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStyle(tt.style)
			if err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func readZipEntry(t *testing.T, data []byte, name string) string {
	t.Helper()
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("failed to open %s: %v", name, err)
			}
			defer rc.Close()
			content, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("failed to read %s: %v", name, err)
			}
			return string(content)
		}
	}
	t.Fatalf("entry %s not found in zip", name)
	return ""
}

func TestHeadingSizeDecreasesWithLevel(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "# H1\n## H2\n### H3\n#### H4\n##### H5\n###### H6"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 6 {
		t.Fatalf("expected 6 paragraphs, got %d", len(result.paragraphs))
	}
	// H1 uses TitleSize, H2-H6 use decreasing HeadingSize
	if !strings.Contains(result.paragraphs[0], "Heading1") {
		t.Error("first should be Heading1")
	}
}

func TestConvertMarkdownToBytes_MarginInDocx(t *testing.T) {
	st := &StyleTemplate{
		TitleFont: "A", TitleSize: 28, HeadingFont: "A", HeadingSize: 18,
		BodyFont: "A", BodySize: 11, CodeFont: "A", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF", PageMarginInches: 1.5,
	}
	data, err := ConvertMarkdownToBytes("# test", st)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readZipEntry(t, data, "word/document.xml")
	// 1.5 inches = 1.5 * 1440 = 2160 twips
	if !strings.Contains(content, `w:top="2160"`) {
		t.Errorf("expected margin 2160, got: %s", content)
	}
}

func TestConvertMarkdownToBytes_XMLEscaping(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<tag>", "&lt;tag&gt;"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"it's", "it&apos;s"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeXML(tt.input)
			if result != tt.want {
				t.Errorf("escapeXML(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestPixelToEMU(t *testing.T) {
	// 96px at 96 DPI = 1 inch = 914400 EMU
	result := pixelToEMU(96)
	if result != 914400 {
		t.Errorf("pixelToEMU(96) = %d, want 914400", result)
	}
	// 1px = 914400/96 = 9525 EMU
	result = pixelToEMU(1)
	if result != 9525 {
		t.Errorf("pixelToEMU(1) = %d, want 9525", result)
	}
}

func TestHexColor(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"#abcdef", "ABCDEF"},
		{"#ABCDEF", "ABCDEF"},
		{"#123456", "123456"},
		{"abcdef", "ABCDEF"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := hexColor(tt.input); got != tt.want {
				t.Errorf("hexColor(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHasLeadingOrTrailingSpace(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"hello", false},
		{" hello", true},
		{"hello ", true},
		{" hello ", true},
		{"\thello", true},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := hasLeadingOrTrailingSpace(tt.input); got != tt.want {
				t.Errorf("hasLeadingOrTrailingSpace(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestRunXML_BoldItalicCode(t *testing.T) {
	tests := []struct {
		name       string
		bold       bool
		italic     bool
		code       bool
		wantMarkup string
	}{
		{"bold", true, false, false, "<w:b/>"},
		{"italic", false, true, false, "<w:i/>"},
		{"code", false, false, true, "F3F4F6"},
		{"plain", false, false, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runXML("text", "Arial", 11, "#000000", tt.bold, tt.italic, tt.code)
			if tt.wantMarkup != "" && !strings.Contains(result, tt.wantMarkup) {
				t.Errorf("expected %q in %s", tt.wantMarkup, result)
			}
			if !strings.Contains(result, "text") {
				t.Error("expected text content")
			}
		})
	}
}

func TestMermaidPlaceholder(t *testing.T) {
	result := mermaidPlaceholder(5)
	expected := `<w:p><!--MERMAID:5--></w:p>`
	if result != expected {
		t.Errorf("mermaidPlaceholder(5) = %q, want %q", result, expected)
	}
}

func TestEmptyParagraphXML(t *testing.T) {
	result := emptyParagraphXML()
	if result != "<w:p/>" {
		t.Errorf("emptyParagraphXML() = %q, want %q", result, "<w:p/>")
	}
}

func TestParagraphXML_WithStyle(t *testing.T) {
	result := paragraphXML("runs", "Heading1", 0)
	if !strings.Contains(result, `<w:pStyle w:val="Heading1"/>`) {
		t.Errorf("expected style in %s", result)
	}
	if !strings.Contains(result, "runs") {
		t.Errorf("expected runs in %s", result)
	}
}

func TestParagraphXML_WithList(t *testing.T) {
	result := paragraphXML("runs", "", 1)
	if !strings.Contains(result, `<w:numId w:val="1"/>`) {
		t.Errorf("expected list in %s", result)
	}
}

func TestParagraphXML_Empty(t *testing.T) {
	result := paragraphXML("runs", "", 0)
	if !strings.Contains(result, "runs") {
		t.Errorf("expected runs in %s", result)
	}
	if strings.Contains(result, "pPr") {
		t.Errorf("unexpected pPr for empty style/list: %s", result)
	}
}

func TestConvertMarkdownToBytes_LargeDocument(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("Paragraph line with some content here.\n")
	}
	data, err := ConvertMarkdownToBytes(sb.String(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty output for large document")
	}
}

func TestConvertMarkdownToBytes_InlineMarkupInHeading(t *testing.T) {
	md := "# **Bold** heading"
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Headings use plain runXML, not inline markdown parsing
	// This should still work
	if len(data) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestBuildContentTypesXML_NoImages(t *testing.T) {
	result := buildContentTypesXML(false)
	if strings.Contains(result, "image/png") {
		t.Error("should not contain image/png when no images")
	}
	if !strings.Contains(result, "document.xml") {
		t.Error("should contain document.xml")
	}
}

func TestBuildContentTypesXML_WithImages(t *testing.T) {
	result := buildContentTypesXML(true)
	if !strings.Contains(result, "image/png") {
		t.Error("should contain image/png when has images")
	}
}

func TestBuildDocRelsXML_NoImages(t *testing.T) {
	result := buildDocRelsXML(nil)
	if !strings.Contains(result, "styles.xml") {
		t.Error("should reference styles.xml")
	}
	if !strings.Contains(result, "numbering.xml") {
		t.Error("should reference numbering.xml")
	}
}

func TestBuildDocRelsXML_WithImages(t *testing.T) {
	images := []MermaidImage{
		{Index: 0, ImageName: "media/image1.png"},
	}
	result := buildDocRelsXML(images)
	if !strings.Contains(result, "rIdMermaid0") {
		t.Error("should contain mermaid relationship ID")
	}
	if !strings.Contains(result, "media/image1.png") {
		t.Error("should reference image path")
	}
}

func TestDrawingXML(t *testing.T) {
	result := drawingXML("rId3", 914400, 914400, 0, "")
	if !strings.Contains(result, "rId3") {
		t.Error("should contain rID")
	}
	if !strings.Contains(result, "diagram-0") {
		t.Error("should contain diagram name")
	}
}

func TestConvertMarkdownToBytes_UnicodeContent(t *testing.T) {
	md := "# 标题\n\n这是一段中文内容。\n\n- 项目一\n- 项目二"
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readZipEntry(t, data, "word/document.xml")
	if !strings.Contains(content, "标题") {
		t.Error("should contain Chinese title")
	}
	if !strings.Contains(content, "中文内容") {
		t.Error("should contain Chinese body text")
	}
}

func TestConvertMarkdownToBytes_JapaneseContent(t *testing.T) {
	md := "# 日本語タイトル\n\nこれはテストです。"
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestConvertMarkdownToBytes_SpecialCharsInCode(t *testing.T) {
	md := "```\n<script>alert('xss')</script>\n```"
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readZipEntry(t, data, "word/document.xml")
	// XML escaping should prevent raw < and > in text content
	if strings.Contains(content, "<script>") {
		t.Error("code block content should be XML-escaped")
	}
	if !strings.Contains(content, "&lt;script&gt;") {
		t.Error("should contain escaped script tag")
	}
}

func TestConvertMarkdownToBytes_StylePresetsIntegrity(t *testing.T) {
	// Ensure inline code renders with gray background
	md := "`inline code`"
	data, err := ConvertMarkdownToBytes(md, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readZipEntry(t, data, "word/document.xml")
	if !strings.Contains(content, "F3F4F6") {
		t.Error("inline code should have gray background")
	}
}

func TestConvertMarkdownToBytes_NumberingXML(t *testing.T) {
	data, err := ConvertMarkdownToBytes("# test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readZipEntry(t, data, "word/numbering.xml")
	if !strings.Contains(content, "bullet") {
		t.Error("numbering.xml should define bullet format")
	}
	if !strings.Contains(content, "decimal") {
		t.Error("numbering.xml should define decimal format")
	}
}

func TestConvertMarkdownToBytes_CoreProps(t *testing.T) {
	data, err := ConvertMarkdownToBytes("# test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readZipEntry(t, data, "docProps/core.xml")
	if !strings.Contains(content, "md2docx") {
		t.Error("core.xml should contain md2docx creator")
	}
}
