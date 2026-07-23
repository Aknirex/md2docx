package style

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/md2docx/cli/internal/converter"
)

// Preset names for built-in style templates.
const (
	PresetUSBusiness   = "us-business"
	PresetUSModern     = "us-modern"
	PresetCNOfficial   = "cn-official"
	PresetCNModern     = "cn-modern"
	PresetJPFormal     = "jp-formal"
	PresetEUClean      = "eu-clean"
	PresetKRStandard   = "kr-standard"
	PresetAcademic     = "academic"
	PresetDefault      = "default"
)

//go:embed presets/us-business.json
var usBusinessJSON []byte

//go:embed presets/us-modern.json
var usModernJSON []byte

//go:embed presets/cn-official.json
var cnOfficialJSON []byte

//go:embed presets/cn-modern.json
var cnModernJSON []byte

//go:embed presets/jp-formal.json
var jpFormalJSON []byte

//go:embed presets/eu-clean.json
var euCleanJSON []byte

//go:embed presets/kr-standard.json
var krStandardJSON []byte

//go:embed presets/academic.json
var academicJSON []byte

//go:embed presets/default.json
var defaultJSON []byte

// presetMap maps preset names to their embedded JSON bytes.
var presetMap = map[string][]byte{
	PresetUSBusiness: usBusinessJSON,
	PresetUSModern:   usModernJSON,
	PresetCNOfficial: cnOfficialJSON,
	PresetCNModern:   cnModernJSON,
	PresetJPFormal:   jpFormalJSON,
	PresetEUClean:    euCleanJSON,
	PresetKRStandard: krStandardJSON,
	PresetAcademic:   academicJSON,
	PresetDefault:    defaultJSON,
}

// AllPresetNames returns all available preset names.
func AllPresetNames() []string {
	return []string{
		PresetUSBusiness,
		PresetUSModern,
		PresetCNOfficial,
		PresetCNModern,
		PresetJPFormal,
		PresetEUClean,
		PresetKRStandard,
		PresetAcademic,
		PresetDefault,
	}
}

// PresetDescriptions returns human-readable descriptions for each preset.
func PresetDescriptions() map[string]string {
	return map[string]string{
		PresetUSBusiness: "US Business – Calibri/Cambria, professional serif headings, blue accent",
		PresetUSModern:   "US Modern – Clean sans-serif, dark tones, minimal",
		PresetCNOfficial: "CN Official – 小标宋_GBK/仿宋_GB2312/楷体_GB2312 (公文风格), Chinese official document style",
		PresetCNModern:   "CN Modern – Noto Sans SC, modern Chinese typography",
		PresetJPFormal:   "JP Formal – Yu Gothic/Mincho, Japanese business document style",
		PresetEUClean:    "EU Clean – Helvetica/Arial, European minimalist design",
		PresetKRStandard: "KR Standard – Malgun Gothic/Nanum, Korean standard document style",
		PresetAcademic:   "Academic – Times New Roman, double-spaced margins, scholarly",
		PresetDefault:    "Default – Aptos Display/Cascadia Mono, modern versatile style",
	}
}

// LoadPreset loads a built-in style preset by name.
func LoadPreset(name string) (*converter.StyleTemplate, error) {
	data, ok := presetMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown preset: %s (available: %v)", name, AllPresetNames())
	}
	var st converter.StyleTemplate
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("unmarshaling preset %s: %w", name, err)
	}
	if err := converter.ValidateStyle(&st); err != nil {
		return nil, fmt.Errorf("validating preset %s: %w", name, err)
	}
	return &st, nil
}

// LoadPresetOrDefault loads a preset, falling back to Default if empty.
func LoadPresetOrDefault(name string) (*converter.StyleTemplate, error) {
	if name == "" {
		name = PresetDefault
	}
	return LoadPreset(name)
}

// LoadTemplateFile loads a style template from a JSON file path.
func LoadTemplateFile(path string) (*converter.StyleTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading template file %s: %w", path, err)
	}
	var st converter.StyleTemplate
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("parsing template %s: %w", path, err)
	}
	if err := converter.ValidateStyle(&st); err != nil {
		return nil, fmt.Errorf("validating template %s: %w", path, err)
	}
	return &st, nil
}

// SaveTemplateFile writes a style template to a JSON file.
func SaveTemplateFile(path string, st *converter.StyleTemplate) error {
	if err := converter.ValidateStyle(st); err != nil {
		return fmt.Errorf("validating before save: %w", err)
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling template: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing template file %s: %w", path, err)
	}
	return nil
}

// LoadStyleTemplate is the unified entry: loads from a preset name or a file path.
// If name is empty, loads Default preset. If it looks like a file path, loads from file.
// Otherwise loads as a preset name.
func LoadStyleTemplate(name string) (*converter.StyleTemplate, error) {
	if name == "" {
		return LoadPreset(PresetDefault)
	}
	// Check if it looks like a file path (contains .json or / or \)
	if _, err := os.Stat(name); err == nil {
		return LoadTemplateFile(name)
	}
	// Try as a preset name
	st, err := LoadPreset(name)
	if err == nil {
		return st, nil
	}
	// Try as a file path even if it didn't match stat
	return LoadTemplateFile(name)
}
