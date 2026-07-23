package converter

import (
	"fmt"
	"strings"
)

// xmlAttr builds an XML attribute string.
func xmlAttr(name, value string) string {
	return fmt.Sprintf(`%s="%s"`, name, escapeXML(value))
}

// escapeXML escapes special XML characters in text content.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// hexColor strips '#' and uppercases a CSS color string.
func hexColor(c string) string {
	return strings.ToUpper(strings.TrimPrefix(c, "#"))
}

// hasLeadingOrTrailingSpace checks if text has whitespace at boundaries.
func hasLeadingOrTrailingSpace(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == ' ' || s[0] == '\t' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t'
}

// runXML generates a <w:r> element with run properties and text.
func runXML(text, font string, size float64, color string, bold, italic, code bool) string {
	var pr strings.Builder
	pr.WriteString("<w:rPr>")
	pr.WriteString(fmt.Sprintf(`<w:rFonts %s %s/>`, xmlAttr("w:ascii", font), xmlAttr("w:hAnsi", font)))
	pr.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, int(size*2)))
	pr.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, hexColor(color)))
	if bold {
		pr.WriteString("<w:b/>")
	}
	if italic {
		pr.WriteString("<w:i/>")
	}
	if code {
		pr.WriteString(`<w:shd w:val="clear" w:fill="F3F4F6"/>`)
	}
	pr.WriteString("</w:rPr>")

	space := ""
	if hasLeadingOrTrailingSpace(text) {
		space = ` xml:space="preserve"`
	}

	return fmt.Sprintf("<w:r>%s<w:t%s>%s</w:t></w:r>", pr.String(), space, escapeXML(text))
}

// paragraphXML generates a <w:p> element with optional style and list properties.
func paragraphXML(runs, styleName string, listID int) string {
	var pr strings.Builder
	if styleName != "" || listID > 0 {
		pr.WriteString("<w:pPr>")
		if styleName != "" {
			pr.WriteString(fmt.Sprintf(`<w:pStyle w:val="%s"/>`, styleName))
		}
		if listID > 0 {
			pr.WriteString(fmt.Sprintf(`<w:numPr><w:ilvl w:val="0"/><w:numId w:val="%d"/></w:numPr>`, listID))
		}
		pr.WriteString("</w:pPr>")
	}
	return fmt.Sprintf("<w:p>%s%s</w:p>", pr.String(), runs)
}

// emptyParagraphXML generates an empty <w:p/> element.
func emptyParagraphXML() string {
	return "<w:p/>"
}

// mermaidPlaceholderFormat generates a placeholder comment for mermaid blocks.
// The placeholder is embedded in a <w:p/> so it occupies a paragraph slot.
func mermaidPlaceholder(index int) string {
	return fmt.Sprintf(`<w:p><!--MERMAID:%d--></w:p>`, index)
}

// parseMermaidPlaceholder extracts the mermaid diagram index from a placeholder paragraph.
// Returns the index and true if the paragraph is a mermaid placeholder, or 0 and false otherwise.
func parseMermaidPlaceholder(p string) (int, bool) {
	const prefix = "<w:p><!--MERMAID:"
	if !strings.HasPrefix(p, prefix) {
		return 0, false
	}
	rest := p[len(prefix):]
	end := strings.Index(rest, "-->")
	if end <= 0 {
		return 0, false
	}
	var idx int
	if _, err := fmt.Sscanf(rest[:end], "%d", &idx); err != nil {
		return 0, false
	}
	return idx, true
}

// drawingXML generates a <w:drawing> paragraph for an embedded image.
// rID is the relationship ID in word/_rels/document.xml.rels (e.g., "rId3").
// widthEMU and heightEMU are in English Metric Units (1 inch = 914400 EMU).
// imageID is the unique identifier for this image within the document.
func drawingXML(rID string, widthEMU, heightEMU int64, imageID int, diagramName string) string {
	// docPr name and title
	name := fmt.Sprintf("diagram-%d", imageID)
	descr := "Mermaid diagram"
	if diagramName != "" {
		descr = fmt.Sprintf("Mermaid diagram: %s", diagramName)
	}

	return fmt.Sprintf(
		`<w:p>`+
			`<w:r>`+
			`<w:drawing>`+
			`<wp:inline distT="0" distB="0" distL="0" distR="0">`+
			`<wp:extent cx="%d" cy="%d"/>`+
			`<wp:docPr id="%d" name="%s" descr="%s"/>`+
			`<a:graphic>`+
			`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
			`<pic:pic>`+
			`<pic:nvPicPr>`+
			`<pic:cNvPr id="0" name="%s"/>`+
			`<pic:cNvPicPr/>`+
			`</pic:nvPicPr>`+
			`<pic:blipFill>`+
			`<a:blip r:embed="%s"/>`+
			`<a:stretch><a:fillRect/></a:stretch>`+
			`</pic:blipFill>`+
			`<pic:spPr>`+
			`<a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>`+
			`</pic:spPr>`+
			`</pic:pic>`+
			`</a:graphicData>`+
			`</a:graphic>`+
			`</wp:inline>`+
			`</w:drawing>`+
			`</w:r>`+
			`</w:p>`,
		widthEMU, heightEMU,
		imageID+10, name, descr, // docPr id offset by 10 to avoid collisions
		fmt.Sprintf("%s.png", name),
		rID,
		widthEMU, heightEMU,
	)
}

// documentXML assembles the full word/document.xml.
// When mermaidImages is provided, MERMAID:N placeholders are replaced with drawing XML.
func documentXML(paragraphs []string, marginInches float64, mermaidImages []MermaidImage) string {
	margin := int(marginInches * 1440)

	// Build a map of placeholder index -> drawing XML + relationship ID
	type replacement struct {
		xml  string
		rID  string
		path string // e.g., "media/image1.png"
	}
	replacements := make(map[int]replacement)
	for _, img := range mermaidImages {
		rID := fmt.Sprintf("rIdMermaid%d", img.Index) // unique rID prefix
		xml := drawingXML(rID, img.WidthEMU, img.HeightEMU, img.Index, "")
		replacements[img.Index] = replacement{xml: xml, rID: rID, path: img.ImageName}
	}

	// Build paragraphs, replacing placeholders
	var parts []string
	for _, p := range paragraphs {
		if idx, ok := parseMermaidPlaceholder(p); ok {
			if repl, ok := replacements[idx]; ok {
				parts = append(parts, repl.xml)
				continue
			}
		}
		parts = append(parts, p)
	}

	docBody := strings.Join(parts, "")

	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
			`<w:document`+
			` xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`+
			` xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"`+
			` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`+
			` xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture"`+
			` xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`+
			`<w:body>%s`+
			`<w:sectPr>`+
			`<w:pgSz w:w="12240" w:h="15840"/>`+
			`<w:pgMar w:top="%d" w:right="%d" w:bottom="%d" w:left="%d" w:header="720" w:footer="720" w:gutter="0"/>`+
			`</w:sectPr>`+
			`</w:body>`+
			`</w:document>`,
		docBody,
		margin, margin, margin, margin,
	)
}

// buildDocRelsXML builds the document relationships XML including image relationships.
func buildDocRelsXML(mermaidImages []MermaidImage) string {
	var imageRels strings.Builder
	for _, img := range mermaidImages {
		rID := fmt.Sprintf("rIdMermaid%d", img.Index)
		imageRels.WriteString(fmt.Sprintf(
			`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`,
			rID, img.ImageName,
		))
	}

	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`+
			`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"/>`+
			`%s`+
			`</Relationships>`,
		imageRels.String(),
	)
}

// buildContentTypesXML builds [Content_Types].xml, adding PNG if mermaid is used.
func buildContentTypesXML(hasImages bool) string {
	base := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
		`<Default Extension="xml" ContentType="application/xml"/>` +
		`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
		`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>` +
		`<Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>` +
		`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>`

	if hasImages {
		base += `<Default Extension="png" ContentType="image/png"/>`
	}

	base += `</Types>`
	return base
}

// stylesXML generates the word/styles.xml for the given style template.
func stylesXML(st *StyleTemplate) string {
	headingSize := st.HeadingSize
	headingStyles := make([]string, 6)
	for level := 0; level < 6; level++ {
		size := headingSize - float64(level)*1.25
		if size < 12 {
			size = 12
		}
		headingStyles[level] = fmt.Sprintf(
			`<w:style w:type="paragraph" w:styleId="Heading%d">`+
				`<w:name w:val="heading %d"/>`+
				`<w:basedOn w:val="Normal"/>`+
				`<w:next w:val="Normal"/>`+
				`<w:qFormat/>`+
				`<w:pPr><w:keepNext/></w:pPr>`+
				`<w:rPr>`+
				`<w:rFonts %s %s/>`+
				`<w:b/>`+
				`<w:color w:val="%s"/>`+
				`<w:sz w:val="%d"/>`+
				`</w:rPr>`+
				`</w:style>`,
			level+1, level+1,
			xmlAttr("w:ascii", st.HeadingFont), xmlAttr("w:hAnsi", st.HeadingFont),
			hexColor(st.AccentColor), int(size*2),
		)
	}

	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
			`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`+
			`<w:docDefaults>`+
			`<w:rPrDefault>`+
			`<w:rPr>`+
			`<w:rFonts %s %s/>`+
			`<w:sz w:val="%d"/>`+
			`<w:color w:val="%s"/>`+
			`</w:rPr>`+
			`</w:rPrDefault>`+
			`</w:docDefaults>`+
			`<w:style w:type="paragraph" w:default="1" w:styleId="Normal">`+
			`<w:name w:val="Normal"/>`+
			`</w:style>`+
			`<w:style w:type="paragraph" w:styleId="CodeBlock">`+
			`<w:name w:val="Code Block"/>`+
			`<w:basedOn w:val="Normal"/>`+
			`<w:rPr>`+
			`<w:rFonts %s %s/>`+
			`<w:sz w:val="%d"/>`+
			`</w:rPr>`+
			`</w:style>`+
			`<w:style w:type="paragraph" w:styleId="Quote">`+
			`<w:name w:val="Quote"/>`+
			`<w:basedOn w:val="Normal"/>`+
			`<w:pPr><w:ind w:left="720"/></w:pPr>`+
			`</w:style>`+
			`%s`+
			`</w:styles>`,
		xmlAttr("w:ascii", st.BodyFont), xmlAttr("w:hAnsi", st.BodyFont),
		int(st.BodySize*2), hexColor(st.TextColor),
		xmlAttr("w:ascii", st.CodeFont), xmlAttr("w:hAnsi", st.CodeFont),
		int(st.CodeSize*2),
		strings.Join(headingStyles, ""),
	)
}

// pixelToEMU converts pixels to EMU at 96 DPI (standard DOCX resolution).
func pixelToEMU(px int) int64 {
	// 1 pixel at 96 DPI = 1/96 inch
	// 1 inch = 914400 EMU
	// EMU = px * 914400 / 96
	return int64(float64(px) * 914400.0 / 96.0)
}
