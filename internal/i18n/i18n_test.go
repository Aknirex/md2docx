package i18n

import (
	"testing"
)

func TestAllLangs(t *testing.T) {
	langs := AllLangs()
	if len(langs) != 8 {
		t.Errorf("expected 8 languages, got %d", len(langs))
	}
	expected := []Lang{EN, ZH_CN, JA, KO, ES, PT_BR, DE, FR}
	for i, want := range expected {
		if langs[i] != want {
			t.Errorf("langs[%d] = %q, want %q", i, langs[i], want)
		}
	}
}

func TestLangName_AllLangs(t *testing.T) {
	tests := []struct {
		lang Lang
		want string
	}{
		{EN, "English"},
		{ZH_CN, "简体中文"},
		{JA, "日本語"},
		{KO, "한국어"},
		{ES, "Español"},
		{PT_BR, "Português (BR)"},
		{DE, "Deutsch"},
		{FR, "Français"},
		{"unknown", "unknown"},
	}
	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			if got := LangName(tt.lang); got != tt.want {
				t.Errorf("LangName(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

func TestDefaultStyleForLang(t *testing.T) {
	tests := []struct {
		lang Lang
		want string
	}{
		{ZH_CN, "cn-official"},
		{JA, "jp-formal"},
		{KO, "kr-standard"},
		{EN, "us-business"},
		{ES, "us-business"},
		{FR, "us-business"},
		{DE, "us-business"},
		{PT_BR, "us-business"},
	}
	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			if got := DefaultStyleForLang(tt.lang); got != tt.want {
				t.Errorf("DefaultStyleForLang(%q) = %q, want %q", tt.lang, got, tt.want)
			}
		})
	}
}

func TestT_ExistingKey(t *testing.T) {
	tests := []struct {
		lang Lang
		key  string
	}{
		{EN, "app_title"},
		{ZH_CN, "app_title"},
		{JA, "app_title"},
		{KO, "app_title"},
		{ES, "app_title"},
		{PT_BR, "app_title"},
		{DE, "app_title"},
		{FR, "app_title"},
	}
	for _, tt := range tests {
		t.Run(string(tt.lang)+"_"+tt.key, func(t *testing.T) {
			got := T(tt.lang, tt.key)
			if got == tt.key {
				t.Errorf("T(%q, %q) returned key itself, expected translation", tt.lang, tt.key)
			}
		})
	}
}

func TestT_MissingKey_FallbackToEnglish(t *testing.T) {
	// For a key that doesn't exist in any language, it returns the key
	got := T(EN, "nonexistent_key_xyz")
	if got != "nonexistent_key_xyz" {
		t.Errorf("expected key fallback, got %q", got)
	}
}

func TestT_MissingLang_FallbackToEnglish(t *testing.T) {
	// Unknown language should fall back to English
	got := T("xx", "app_title")
	enGot := T(EN, "app_title")
	if got != enGot {
		t.Errorf("expected English fallback for unknown lang, got %q vs %q", got, enGot)
	}
}

func TestT_MissingLangAndKey(t *testing.T) {
	got := T("xx", "nonexistent_key")
	if got != "nonexistent_key" {
		t.Errorf("expected key as fallback, got %q", got)
	}
}

func TestPresetDescription_ExistingPreset(t *testing.T) {
	tests := []struct {
		lang   Lang
		preset string
	}{
		{EN, "default"},
		{ZH_CN, "cn-official"},
		{JA, "jp-formal"},
	}
	for _, tt := range tests {
		t.Run(string(tt.lang)+"_"+tt.preset, func(t *testing.T) {
			got := PresetDescription(tt.lang, tt.preset)
			if got == tt.preset {
				t.Errorf("expected description, got preset name itself")
			}
		})
	}
}

func TestPresetDescription_UnknownPreset(t *testing.T) {
	got := PresetDescription(EN, "nonexistent")
	if got != "nonexistent" {
		t.Errorf("expected preset name as fallback, got %q", got)
	}
}

func TestT_AllKeysExistInEnglish(t *testing.T) {
	// Verify all critical keys exist in English
	keys := []string{
		"app_title", "app_short",
		"lang_select_title", "lang_select_help", "lang_saved", "lang_saved_hint",
		"tui_input_title", "tui_output_title", "tui_style_title",
		"tui_confirm_title", "tui_confirm_input", "tui_confirm_output",
		"tui_confirm_style", "tui_confirm_mermaid", "tui_confirm_convert",
		"tui_confirm_back", "tui_converting", "tui_done", "tui_done_msg",
		"tui_error_msg", "tui_filename_label", "tui_nav_help", "tui_nav_quit",
		"cli_ok", "cli_error", "cli_bytes", "cli_template_created",
		"err_markdown_not_found", "err_output_ext", "err_style_not_found",
		"err_converting", "err_tui",
	}
	for _, key := range keys {
		got := T(EN, key)
		if got == key {
			t.Errorf("key %q not found in English translations", key)
		}
	}
}

func TestT_AllKeysExistInAllLanguages(t *testing.T) {
	// Get all English keys
	enDict := dict[EN]
	langs := AllLangs()
	for _, lang := range langs {
		if lang == EN {
			continue
		}
		langDict, ok := dict[lang]
		if !ok {
			t.Errorf("language %q has no translation map", lang)
			continue
		}
		for key := range enDict {
			if _, exists := langDict[key]; !exists {
				t.Errorf("language %q missing key %q", lang, key)
			}
		}
	}
}

func TestT_FormatStrings(t *testing.T) {
	// Verify format strings are consistent across languages
	formatKeys := []string{"lang_saved", "tui_done_msg", "tui_error_msg"}
	langs := AllLangs()
	for _, key := range formatKeys {
		enVal := T(EN, key)
		for _, lang := range langs {
			if lang == EN {
				continue
			}
			langVal := T(lang, key)
			if langVal == key {
				// Missing key - will be caught by other test
				continue
			}
			// Both should have same number of format verbs
			enFormats := countFormatVerbs(enVal)
			langFormats := countFormatVerbs(langVal)
			if enFormats != langFormats {
				t.Errorf("key %q in %q: format verb count mismatch (%d vs %d in en)", key, lang, langFormats, enFormats)
			}
		}
	}
}

func countFormatVerbs(s string) int {
	count := 0
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '%' && s[i+1] != '%' {
			count++
		}
	}
	return count
}
