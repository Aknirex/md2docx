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

// documentXML assembles the full word/document.xml.
func documentXML(paragraphs []string, marginInches float64) string {
	margin := int(marginInches * 1440)
	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
			`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`+
			`<w:body>%s`+
			`<w:sectPr>`+
			`<w:pgSz w:w="12240" w:h="15840"/>`+
			`<w:pgMar w:top="%d" w:right="%d" w:bottom="%d" w:left="%d" w:header="720" w:footer="720" w:gutter="0"/>`+
			`</w:sectPr>`+
			`</w:body>`+
			`</w:document>`,
		strings.Join(paragraphs, ""),
		margin, margin, margin, margin,
	)
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
