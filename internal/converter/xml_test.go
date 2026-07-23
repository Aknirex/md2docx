package converter

import (
	"archive/zip"
	"strings"
	"testing"
)

func TestEscapeXML_AllChars(t *testing.T) {
	input := `<hello "world" & 'test'`
	want := `&lt;hello &quot;world&quot; &amp; &apos;test&apos;`
	got := escapeXML(input)
	if got != want {
		t.Errorf("escapeXML(%q) = %q, want %q", input, got, want)
	}
}

func TestEscapeXML_Empty(t *testing.T) {
	if got := escapeXML(""); got != "" {
		t.Errorf("escapeXML(\"\") = %q, want \"\"", got)
	}
}

func TestEscapeXML_NoSpecialChars(t *testing.T) {
	input := "hello world"
	if got := escapeXML(input); got != input {
		t.Errorf("escapeXML(%q) = %q, want %q", input, got, input)
	}
}

func TestXmlAttr(t *testing.T) {
	result := xmlAttr("w:ascii", "Arial")
	want := `w:ascii="Arial"`
	if result != want {
		t.Errorf("xmlAttr = %q, want %q", result, want)
	}
}

func TestXmlAttr_WithSpecialChars(t *testing.T) {
	result := xmlAttr("name", `<test>`)
	want := `name="&lt;test&gt;"`
	if result != want {
		t.Errorf("xmlAttr = %q, want %q", result, want)
	}
}

func TestRunXML_ContainsFont(t *testing.T) {
	result := runXML("hello", "Courier New", 10, "#FF0000", false, false, false)
	if !strings.Contains(result, "Courier New") {
		t.Error("should contain font name")
	}
	if !strings.Contains(result, "hello") {
		t.Error("should contain text")
	}
}

func TestRunXML_SizeDoubled(t *testing.T) {
	// Font size in OOXML is in half-points (size * 2)
	result := runXML("text", "Arial", 12, "#000000", false, false, false)
	if !strings.Contains(result, `w:val="24"`) {
		t.Errorf("expected size 24 (12*2) in %s", result)
	}
}

func TestRunXML_CodeBackground(t *testing.T) {
	result := runXML("code", "Courier", 10, "#000000", false, false, true)
	if !strings.Contains(result, "F3F4F6") {
		t.Error("code should have gray background")
	}
	if !strings.Contains(result, "w:shd") {
		t.Error("code should have shading element")
	}
}

func TestRunXML_PreservesLeadingSpace(t *testing.T) {
	result := runXML(" text", "Arial", 11, "#000000", false, false, false)
	if !strings.Contains(result, `xml:space="preserve"`) {
		t.Error("should preserve leading space")
	}
}

func TestRunXML_NoPreserveWithoutSpace(t *testing.T) {
	result := runXML("text", "Arial", 11, "#000000", false, false, false)
	if strings.Contains(result, "preserve") {
		t.Error("should not preserve when no leading/trailing space")
	}
}

func TestDocumentXML_SectionProperties(t *testing.T) {
	paragraphs := []string{"<w:p/>"}
	result := documentXML(paragraphs, 1.0, nil)
	if !strings.Contains(result, "w:sectPr") {
		t.Error("should contain section properties")
	}
	if !strings.Contains(result, "w:pgSz") {
		t.Error("should contain page size")
	}
	if !strings.Contains(result, "w:pgMar") {
		t.Error("should contain page margins")
	}
}

func TestDocumentXML_MarginCalculation(t *testing.T) {
	paragraphs := []string{"<w:p/>"}
	// 1.0 inch = 1440 twips
	result := documentXML(paragraphs, 1.0, nil)
	if !strings.Contains(result, `w:top="1440"`) {
		t.Errorf("expected margin 1440 for 1.0 inch, got: %s", result)
	}
}

func TestDocumentXML_WithMermaidReplacement(t *testing.T) {
	placeholder := mermaidPlaceholder(0)
	paragraphs := []string{placeholder}
	images := []MermaidImage{
		{Index: 0, ImageName: "media/image1.png", WidthEMU: 914400, HeightEMU: 914400},
	}
	result := documentXML(paragraphs, 1.0, images)
	// Placeholder should be replaced with drawing
	if strings.Contains(result, "MERMAID:") {
		t.Error("mermaid placeholder should be replaced")
	}
	if !strings.Contains(result, "w:drawing") {
		t.Error("should contain drawing element")
	}
}

func TestDocumentXML_MermaidPlaceholderKeptWhenNoImage(t *testing.T) {
	placeholder := mermaidPlaceholder(0)
	paragraphs := []string{placeholder}
	result := documentXML(paragraphs, 1.0, nil)
	// Without matching image, placeholder stays as-is
	if !strings.Contains(result, "MERMAID:0") {
		t.Error("placeholder should remain when no matching image")
	}
}

func TestStylesXML_ContainsAllHeadingLevels(t *testing.T) {
	st := &StyleTemplate{
		TitleFont: "A", TitleSize: 28, HeadingFont: "A", HeadingSize: 18,
		BodyFont: "A", BodySize: 11, CodeFont: "A", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF", PageMarginInches: 1.0,
	}
	result := stylesXML(st)
	for i := 1; i <= 6; i++ {
		styleID := "Heading" + string(rune('0'+i))
		if !strings.Contains(result, styleID) {
			t.Errorf("styles.xml should contain %s", styleID)
		}
	}
}

func TestStylesXML_HeadingSizeFloor(t *testing.T) {
	st := &StyleTemplate{
		TitleFont: "A", TitleSize: 28, HeadingFont: "A", HeadingSize: 14,
		BodyFont: "A", BodySize: 11, CodeFont: "A", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF", PageMarginInches: 1.0,
	}
	result := stylesXML(st)
	// H6 = 14 - 5*1.25 = 7.75, floored to 12 -> 24 half-points
	if !strings.Contains(result, `w:val="24"`) {
		t.Error("heading size should be floored at 12pt (24 half-points)")
	}
}

func TestConvertInlineMarkdown_EmptyText(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("", st)
	// Should still produce a run
	if !strings.Contains(result, "<w:r>") {
		t.Error("empty text should produce a run element")
	}
}

func TestConvertInlineMarkdown_ConsecutiveMarkup(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("**bold** and *italic*", st)
	if !strings.Contains(result, "<w:b/>") {
		t.Error("expected bold")
	}
	if !strings.Contains(result, "<w:i/>") {
		t.Error("expected italic")
	}
}

func TestConvertInlineMarkdown_CodeWithSpecialChars(t *testing.T) {
	st := &StyleTemplate{
		BodyFont: "Arial", BodySize: 11, CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
	}
	result := convertInlineMarkdown("`<html>`", st)
	if !strings.Contains(result, "&lt;html&gt;") {
		t.Error("special chars in code should be escaped")
	}
}

func TestParseMarkdown_MultipleCodeBlocks(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "```\nblock1\n```\n\ntext\n\n```\nblock2\n```"
	result := parseMarkdown(input, st, false)
	// 1 code line + 1 empty + 1 text + 1 empty + 1 code line = 5
	if len(result.paragraphs) != 5 {
		t.Fatalf("expected 5 paragraphs, got %d", len(result.paragraphs))
	}
}

func TestParseMarkdown_NestedListMarkers(t *testing.T) {
	st := resolveDefaultStyle(nil)
	// Test different list markers
	input := "- dash\n+ plus\n* star"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 3 {
		t.Fatalf("expected 3 paragraphs, got %d", len(result.paragraphs))
	}
	for _, p := range result.paragraphs {
		if !strings.Contains(p, `<w:numId w:val="1"/>`) {
			t.Errorf("all should be unordered list items: %s", p)
		}
	}
}

func TestParseMarkdown_OrderedListWithParen(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "1) first\n2) second"
	result := parseMarkdown(input, st, false)
	if len(result.paragraphs) != 2 {
		t.Fatalf("expected 2 paragraphs, got %d", len(result.paragraphs))
	}
	for _, p := range result.paragraphs {
		if !strings.Contains(p, `<w:numId w:val="2"/>`) {
			t.Errorf("should be ordered list items: %s", p)
		}
	}
}

func TestParseMarkdown_MermaidBlockContent(t *testing.T) {
	st := resolveDefaultStyle(nil)
	input := "```mermaid\ngraph LR\n  A-->B\n  B-->C\n```"
	result := parseMarkdown(input, st, true)
	if len(result.mermaidBlocks) != 1 {
		t.Fatalf("expected 1 mermaid block, got %d", len(result.mermaidBlocks))
	}
	diagram := result.mermaidBlocks[0].diagram
	if !strings.Contains(diagram, "graph LR") {
		t.Error("diagram should contain graph directive")
	}
	if !strings.Contains(diagram, "A-->B") {
		t.Error("diagram should contain edges")
	}
}

func TestParseMarkdown_CodeBlockFenceVariations(t *testing.T) {
	st := resolveDefaultStyle(nil)
	tests := []struct {
		name  string
		input string
	}{
		{"triple backtick", "```\ncode\n```"},
		{"four backticks", "````\ncode\n````"},
		{"with spaces", "   ```\ncode\n   ```"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMarkdown(tt.input, st, false)
			if len(result.paragraphs) != 1 {
				t.Fatalf("expected 1 code paragraph, got %d", len(result.paragraphs))
			}
			if !strings.Contains(result.paragraphs[0], "CodeBlock") {
				t.Error("should be CodeBlock style")
			}
		})
	}
}

func TestConvertMarkdownToBytes_EmptyDocument(t *testing.T) {
	data, err := ConvertMarkdownToBytes("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected valid DOCX even for empty input")
	}
	r, err := zip.NewReader(strings.NewReader(string(data)), int64(len(data)))
	if err != nil {
		t.Fatalf("should be valid zip: %v", err)
	}
	_ = r
}

func TestConvertMarkdownToBytes_OnlyWhitespace(t *testing.T) {
	data, err := ConvertMarkdownToBytes("   \n\n   \n", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected valid DOCX for whitespace-only input")
	}
}

func TestParseMermaidPlaceholder_Valid(t *testing.T) {
	tests := []struct {
		input   string
		wantIdx int
		wantOk  bool
	}{
		{`<w:p><!--MERMAID:0--></w:p>`, 0, true},
		{`<w:p><!--MERMAID:5--></w:p>`, 5, true},
		{`<w:p><!--MERMAID:999--></w:p>`, 999, true},
		{`<w:p/>`, 0, false},
		{`<w:p><w:r><w:t>text</w:t></w:r></w:p>`, 0, false},
		{`<w:p><!--OTHER:0--></w:p>`, 0, false},
		{``, 0, false},
	}
	for _, tt := range tests {
		name := tt.input
		if len(name) > 30 {
			name = name[:30]
		}
		t.Run(name, func(t *testing.T) {
			idx, ok := parseMermaidPlaceholder(tt.input)
			if ok != tt.wantOk {
				t.Errorf("parseMermaidPlaceholder(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if ok && idx != tt.wantIdx {
				t.Errorf("parseMermaidPlaceholder(%q) idx = %d, want %d", tt.input, idx, tt.wantIdx)
			}
		})
	}
}

func TestParseMermaidPlaceholder_RoundTrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		placeholder := mermaidPlaceholder(i)
		idx, ok := parseMermaidPlaceholder(placeholder)
		if !ok {
			t.Errorf("round trip failed for index %d: not recognized", i)
		}
		if idx != i {
			t.Errorf("round trip failed: got index %d, want %d", idx, i)
		}
	}
}
