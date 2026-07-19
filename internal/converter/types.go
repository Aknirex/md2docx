package converter

// StyleTemplate defines the visual style for a DOCX output.
type StyleTemplate struct {
	TitleFont        string  `json:"titleFont"`
	TitleSize        float64 `json:"titleSize"`
	HeadingFont      string  `json:"headingFont"`
	HeadingSize      float64 `json:"headingSize"`
	BodyFont         string  `json:"bodyFont"`
	BodySize         float64 `json:"bodySize"`
	CodeFont         string  `json:"codeFont"`
	CodeSize         float64 `json:"codeSize"`
	TextColor        string  `json:"textColor"`
	AccentColor      string  `json:"accentColor"`
	PageMarginInches float64 `json:"pageMarginInches"`
}

// ConversionResult holds the result of a successful conversion.
type ConversionResult struct {
	OutputPath string
	Bytes      int64
}
