package converter

import (
	"testing"
)

func TestEncodeMermaidInk_ProducesNonEmptyOutput(t *testing.T) {
	result := encodeMermaidInk("graph TD\n  A-->B", "default")
	if result == "" {
		t.Error("expected non-empty encoded output")
	}
}

func TestEncodeMermaidInk_DifferentThemes(t *testing.T) {
	themes := []string{"default", "neutral", "dark", "forest"}
	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			result := encodeMermaidInk("graph TD", theme)
			if result == "" {
				t.Errorf("expected non-empty output for theme %q", theme)
			}
		})
	}
}

func TestEncodeMermaidInk_DifferentDiagrams(t *testing.T) {
	diagrams := []string{
		"graph TD\n  A-->B",
		"sequenceDiagram\n  A->>B: Hello",
		"classDiagram\n  class Animal",
		"stateDiagram-v2\n  [*] --> Active",
		"pie title Pets\n  \"Dogs\" : 50\n  \"Cats\" : 40",
	}
	for _, diagram := range diagrams {
		t.Run(diagram[:10], func(t *testing.T) {
			result := encodeMermaidInk(diagram, "default")
			if result == "" {
				t.Error("expected non-empty output")
			}
		})
	}
}

func TestEncodeMermaidInk_URLSafe(t *testing.T) {
	result := encodeMermaidInk("graph TD\n  A-->B", "default")
	// base64url should not contain +, /, or =
	for _, c := range result {
		if c == '+' || c == '/' || c == '=' {
			t.Errorf("URL-safe base64 should not contain %q, got %q", string(c), result)
		}
	}
}

func TestReadPNGDimensions_InvalidData(t *testing.T) {
	_, _, err := readPNGDimensions([]byte("not a png"))
	if err == nil {
		t.Error("expected error for invalid PNG data")
	}
}

func TestReadPNGDimensions_EmptyData(t *testing.T) {
	_, _, err := readPNGDimensions([]byte{})
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestReadPNGDimensions_TruncatedData(t *testing.T) {
	// Valid PNG header but truncated
	data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	_, _, err := readPNGDimensions(data)
	if err == nil {
		t.Error("expected error for truncated PNG")
	}
}

func TestMermaidCLIRenderer_DefaultPath(t *testing.T) {
	r := &MermaidCLIRenderer{}
	// Just verify it doesn't panic on construction
	if r.MMDCPath != "" {
		t.Errorf("expected empty default path, got %q", r.MMDCPath)
	}
	if r.Theme != "" {
		t.Errorf("expected empty default theme, got %q", r.Theme)
	}
}

func TestMermaidInkRenderer_DefaultFields(t *testing.T) {
	r := &MermaidInkRenderer{}
	// Verify zero values
	if r.BaseURL != "" {
		t.Errorf("expected empty BaseURL, got %q", r.BaseURL)
	}
	if r.Theme != "" {
		t.Errorf("expected empty Theme, got %q", r.Theme)
	}
}
