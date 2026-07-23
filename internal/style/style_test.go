package style

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/md2docx/cli/internal/converter"
)

func TestAllPresetNames(t *testing.T) {
	names := AllPresetNames()
	if len(names) != 9 {
		t.Errorf("expected 9 presets, got %d", len(names))
	}
	expected := []string{
		"us-business", "us-modern", "cn-official", "cn-modern",
		"jp-formal", "eu-clean", "kr-standard", "academic", "default",
	}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("preset[%d] = %q, want %q", i, names[i], want)
		}
	}
}

func TestLoadPreset_AllPresets(t *testing.T) {
	names := AllPresetNames()
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			st, err := LoadPreset(name)
			if err != nil {
				t.Fatalf("failed to load preset %q: %v", name, err)
			}
			if st == nil {
				t.Fatal("expected non-nil style template")
			}
			if st.TitleFont == "" {
				t.Error("TitleFont should not be empty")
			}
			if st.BodyFont == "" {
				t.Error("BodyFont should not be empty")
			}
			if st.BodySize <= 0 {
				t.Error("BodySize should be positive")
			}
		})
	}
}

func TestLoadPreset_Unknown(t *testing.T) {
	_, err := LoadPreset("nonexistent-preset")
	if err == nil {
		t.Error("expected error for unknown preset")
	}
}

func TestLoadPresetOrDefault_Empty(t *testing.T) {
	st, err := LoadPresetOrDefault("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.TitleFont == "" {
		t.Error("default preset should have TitleFont")
	}
}

func TestLoadPresetOrDefault_Named(t *testing.T) {
	st, err := LoadPresetOrDefault("us-business")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st == nil {
		t.Fatal("expected non-nil style")
	}
}

func TestLoadStyleTemplate_Empty(t *testing.T) {
	st, err := LoadStyleTemplate("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st == nil {
		t.Fatal("expected non-nil for empty string (loads default)")
	}
}

func TestLoadStyleTemplate_PresetName(t *testing.T) {
	st, err := LoadStyleTemplate("cn-official")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st == nil {
		t.Fatal("expected non-nil")
	}
}

func TestLoadStyleTemplate_UnknownPreset(t *testing.T) {
	_, err := LoadStyleTemplate("unknown-preset-name")
	if err == nil {
		t.Error("expected error for unknown preset that is not a file")
	}
}

func TestLoadTemplateFile_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test-style.json")

	st := &converter.StyleTemplate{
		TitleFont: "TestTitle", TitleSize: 28,
		HeadingFont: "TestHeading", HeadingSize: 18,
		BodyFont: "TestBody", BodySize: 11,
		CodeFont: "TestCode", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#FF0000",
		PageMarginInches: 1.0,
	}
	data, _ := json.Marshal(st)
	os.WriteFile(path, data, 0644)

	loaded, err := LoadTemplateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.TitleFont != "TestTitle" {
		t.Errorf("TitleFont = %q, want %q", loaded.TitleFont, "TestTitle")
	}
	if loaded.AccentColor != "#FF0000" {
		t.Errorf("AccentColor = %q, want %q", loaded.AccentColor, "#FF0000")
	}
}

func TestLoadTemplateFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := LoadTemplateFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadTemplateFile_InvalidStyle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.json")
	// Valid JSON but invalid style (missing required fields)
	os.WriteFile(path, []byte(`{"titleFont": ""}`), 0644)

	_, err := LoadTemplateFile(path)
	if err == nil {
		t.Error("expected error for invalid style")
	}
}

func TestLoadTemplateFile_NotFound(t *testing.T) {
	_, err := LoadTemplateFile("/nonexistent/path/style.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveTemplateFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "saved.json")

	st := &converter.StyleTemplate{
		TitleFont: "Arial", TitleSize: 28,
		HeadingFont: "Arial", HeadingSize: 18,
		BodyFont: "Arial", BodySize: 11,
		CodeFont: "Courier", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
		PageMarginInches: 1.0,
	}

	err := SaveTemplateFile(path, st)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it can be loaded back
	loaded, err := LoadTemplateFile(path)
	if err != nil {
		t.Fatalf("failed to load saved template: %v", err)
	}
	if loaded.TitleFont != "Arial" {
		t.Errorf("TitleFont = %q, want %q", loaded.TitleFont, "Arial")
	}
}

func TestSaveTemplateFile_InvalidStyle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	st := &converter.StyleTemplate{} // empty, invalid
	err := SaveTemplateFile(path, st)
	if err == nil {
		t.Error("expected error when saving invalid style")
	}
}

func TestSaveTemplateFile_LoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "roundtrip.json")

	original := &converter.StyleTemplate{
		TitleFont: "CustomFont", TitleSize: 32,
		HeadingFont: "HeadingFont", HeadingSize: 20,
		BodyFont: "BodyFont", BodySize: 12,
		CodeFont: "CodeFont", CodeSize: 9,
		TextColor: "#111111", AccentColor: "#EEEEEE",
		PageMarginInches: 1.25,
	}

	if err := SaveTemplateFile(path, original); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := LoadTemplateFile(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if loaded.TitleFont != original.TitleFont {
		t.Errorf("TitleFont mismatch")
	}
	if loaded.TitleSize != original.TitleSize {
		t.Errorf("TitleSize mismatch")
	}
	if loaded.PageMarginInches != original.PageMarginInches {
		t.Errorf("PageMarginInches mismatch")
	}
}

func TestPresetDescriptions(t *testing.T) {
	descs := PresetDescriptions()
	names := AllPresetNames()
	for _, name := range names {
		if _, ok := descs[name]; !ok {
			t.Errorf("missing description for preset %q", name)
		}
	}
}

func TestLoadStyleTemplate_FilePath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.json")

	st := &converter.StyleTemplate{
		TitleFont: "FileFont", TitleSize: 28,
		HeadingFont: "FileFont", HeadingSize: 18,
		BodyFont: "FileFont", BodySize: 11,
		CodeFont: "FileFont", CodeSize: 10,
		TextColor: "#000000", AccentColor: "#0000FF",
		PageMarginInches: 1.0,
	}
	data, _ := json.Marshal(st)
	os.WriteFile(path, data, 0644)

	loaded, err := LoadStyleTemplate(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.TitleFont != "FileFont" {
		t.Errorf("TitleFont = %q, want %q", loaded.TitleFont, "FileFont")
	}
}
