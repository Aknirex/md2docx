package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/md2docx/cli/internal/i18n"
)

// Config holds persistent user settings.
type Config struct {
	Lang         i18n.Lang `json:"lang"`
	DefaultStyle string    `json:"defaultStyle"`
	FirstRun     bool      `json:"firstRun"`
}

// Path returns the config file path (~/.config/md2docx/config.json).
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "md2docx")
	return filepath.Join(dir, "config.json"), nil
}

// Load reads the config file. If it doesn't exist or is corrupt, returns a zero Config.
func Load() (*Config, error) {
	cfgPath, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{FirstRun: true}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		// Corrupt config — start fresh
		return &Config{FirstRun: true}, nil
	}

	// Validate lang
	if cfg.Lang == "" {
		cfg.FirstRun = true
	}

	return &cfg, nil
}

// Save writes the config to disk.
func Save(cfg *Config) error {
	cfgPath, err := Path()
	if err != nil {
		return err
	}

	dir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, data, 0644)
}

// SetLanguage updates the config language and auto-selects the default style.
func SetLanguage(lang i18n.Lang) (*Config, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	cfg.Lang = lang
	cfg.FirstRun = false

	// Only auto-set default style if not already set
	if cfg.DefaultStyle == "" {
		cfg.DefaultStyle = i18n.DefaultStyleForLang(lang)
	}

	if err := Save(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
