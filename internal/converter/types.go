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

// MermaidImage holds a rendered mermaid diagram ready for DOCX embedding.
type MermaidImage struct {
	Index     int    // position in the paragraph list
	ImageName string // filename within word/media/, e.g. "image1.png"
	PNGBytes  []byte // PNG image data
	WidthEMU  int64  // image width in EMU (English Metric Units)
	HeightEMU int64  // image height in EMU
}

// MermaidRenderer converts a mermaid diagram definition to a PNG image.
type MermaidRenderer interface {
	// Render converts a mermaid diagram string to PNG bytes.
	// Returns the PNG bytes, the image width and height in pixels.
	Render(diagram string) (pngBytes []byte, widthPx, heightPx int, err error)
}

// conversionConfig holds optional conversion settings.
type conversionConfig struct {
	Style   *StyleTemplate
	Mermaid MermaidRenderer // nil = skip mermaid rendering, render as code
}

// ConversionOption is a functional option for conversions.
type ConversionOption func(*conversionConfig)

// WithMermaid enables mermaid diagram rendering with the given renderer.
func WithMermaid(r MermaidRenderer) ConversionOption {
	return func(cfg *conversionConfig) {
		cfg.Mermaid = r
	}
}
