package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/md2docx/cli/internal/i18n"
)

func TestPath(t *testing.T) {
	p, err := Path()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == "" {
		t.Error("expected non-empty path")
	}
	if filepath.Ext(p) != ".json" {
		t.Errorf("expected .json extension, got %s", filepath.Ext(p))
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	// Save and restore home config if it exists
	cfgPath, _ := Path()
	backup, _ := os.ReadFile(cfgPath)
	defer func() {
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		}
	}()

	// Remove config temporarily
	os.Remove(cfgPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.FirstRun {
		t.Error("expected FirstRun=true when no config file")
	}
	if cfg.Lang != "" {
		t.Error("expected empty Lang when no config file")
	}
}

func TestLoad_CorruptConfig(t *testing.T) {
	cfgPath, _ := Path()
	dir := filepath.Dir(cfgPath)
	os.MkdirAll(dir, 0755)
	backup, _ := os.ReadFile(cfgPath)
	defer func() {
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		} else {
			os.Remove(cfgPath)
		}
	}()

	os.WriteFile(cfgPath, []byte("not json!!!"), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error for corrupt config: %v", err)
	}
	if !cfg.FirstRun {
		t.Error("expected FirstRun=true for corrupt config")
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	cfgPath, _ := Path()
	dir := filepath.Dir(cfgPath)
	os.MkdirAll(dir, 0755)
	backup, _ := os.ReadFile(cfgPath)
	defer func() {
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		} else {
			os.Remove(cfgPath)
		}
	}()

	valid := &Config{Lang: i18n.ZH_CN, DefaultStyle: "cn-official", FirstRun: false}
	data, _ := json.Marshal(valid)
	os.WriteFile(cfgPath, data, 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Lang != i18n.ZH_CN {
		t.Errorf("Lang = %q, want %q", cfg.Lang, i18n.ZH_CN)
	}
	if cfg.DefaultStyle != "cn-official" {
		t.Errorf("DefaultStyle = %q, want %q", cfg.DefaultStyle, "cn-official")
	}
	if cfg.FirstRun {
		t.Error("FirstRun should be false")
	}
}

func TestSave_LoadRoundTrip(t *testing.T) {
	cfgPath, _ := Path()
	dir := filepath.Dir(cfgPath)
	os.MkdirAll(dir, 0755)
	backup, _ := os.ReadFile(cfgPath)
	defer func() {
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		} else {
			os.Remove(cfgPath)
		}
	}()

	original := &Config{Lang: i18n.JA, DefaultStyle: "jp-formal", FirstRun: false}
	if err := Save(original); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Lang != original.Lang {
		t.Errorf("Lang mismatch: %q vs %q", loaded.Lang, original.Lang)
	}
	if loaded.DefaultStyle != original.DefaultStyle {
		t.Errorf("DefaultStyle mismatch: %q vs %q", loaded.DefaultStyle, original.DefaultStyle)
	}
}

func TestSetLanguage_NewConfig(t *testing.T) {
	cfgPath, _ := Path()
	dir := filepath.Dir(cfgPath)
	os.MkdirAll(dir, 0755)
	backup, _ := os.ReadFile(cfgPath)
	defer func() {
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		} else {
			os.Remove(cfgPath)
		}
	}()

	// Remove existing config
	os.Remove(cfgPath)

	cfg, err := SetLanguage(i18n.KO)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Lang != i18n.KO {
		t.Errorf("Lang = %q, want %q", cfg.Lang, i18n.KO)
	}
	if cfg.FirstRun {
		t.Error("FirstRun should be false after SetLanguage")
	}
	if cfg.DefaultStyle != "kr-standard" {
		t.Errorf("DefaultStyle = %q, want %q", cfg.DefaultStyle, "kr-standard")
	}
}

func TestSetLanguage_PreservesExistingStyle(t *testing.T) {
	cfgPath, _ := Path()
	dir := filepath.Dir(cfgPath)
	os.MkdirAll(dir, 0755)
	backup, _ := os.ReadFile(cfgPath)
	defer func() {
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		} else {
			os.Remove(cfgPath)
		}
	}()

	// Save existing config with a custom style
	existing := &Config{Lang: i18n.EN, DefaultStyle: "custom-style", FirstRun: false}
	data, _ := json.Marshal(existing)
	os.WriteFile(cfgPath, data, 0644)

	cfg, err := SetLanguage(i18n.ZH_CN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should preserve existing custom style
	if cfg.DefaultStyle != "custom-style" {
		t.Errorf("DefaultStyle = %q, want %q (should be preserved)", cfg.DefaultStyle, "custom-style")
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	cfgPath, _ := Path()
	dir := filepath.Dir(cfgPath)
	backup, _ := os.ReadFile(cfgPath)
	// Remove the entire md2docx config dir
	os.RemoveAll(dir)
	defer func() {
		os.MkdirAll(dir, 0755)
		if backup != nil {
			os.WriteFile(cfgPath, backup, 0644)
		}
	}()

	cfg := &Config{Lang: i18n.FR, DefaultStyle: "eu-clean", FirstRun: false}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save should create directory: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("config file should exist after Save")
	}
}
