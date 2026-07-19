package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/md2docx/cli/internal/converter"
	"github.com/md2docx/cli/internal/i18n"
	"github.com/md2docx/cli/internal/style"
)

type AgentOutput struct {
	Success    bool                     `json:"success"`
	OutputPath string                   `json:"outputPath,omitempty"`
	Bytes      int64                    `json:"bytes,omitempty"`
	Error      string                   `json:"error,omitempty"`
	Style      *converter.StyleTemplate `json:"style,omitempty"`
	Presets    []string                 `json:"presets,omitempty"`
}

type ConvertOptions struct {
	InputPath     string
	OutputPath    string
	StyleRef      string
	PlainOutput   bool
	Mermaid       bool
	MermaidServer string
	MermaidTheme  string
	Lang          i18n.Lang
}

func Convert(opts ConvertOptions) {
	var result AgentOutput
	t := func(k string) string { return i18n.T(opts.Lang, k) }

	if filepath.Ext(opts.OutputPath) != ".docx" {
		opts.OutputPath += ".docx"
	}

	var convertOpts []converter.ConversionOption
	if opts.Mermaid {
		r := &converter.MermaidInkRenderer{Theme: opts.MermaidTheme}
		if opts.MermaidServer != "" {
			r.BaseURL = opts.MermaidServer
		}
		convertOpts = append(convertOpts, converter.WithMermaid(r))
	}

	cr, err := style.ResolveAndConvertWithOptions(opts.InputPath, opts.OutputPath, opts.StyleRef, convertOpts...)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		printResult(result, opts.PlainOutput, t)
		os.Exit(1)
		return
	}

	result.Success = true
	result.OutputPath = cr.OutputPath
	result.Bytes = cr.Bytes
	printResult(result, opts.PlainOutput, t)
}

func ListPresets(plainOutput bool, lang i18n.Lang) {
	t := func(k string) string { return i18n.T(lang, k) }
	names := style.AllPresetNames()

	if plainOutput {
		for _, name := range names {
			fmt.Printf("%-20s %s\n", name, i18n.PresetDescription(lang, name))
		}
		return
	}

	result := AgentOutput{Success: true, Presets: names}
	printResult(result, false, t)
}

func ShowPreset(name string, plainOutput bool, lang i18n.Lang) {
	t := func(k string) string { return i18n.T(lang, k) }
	st, err := style.LoadPreset(name)
	if err != nil {
		result := AgentOutput{Success: false, Error: err.Error()}
		printResult(result, plainOutput, t)
		os.Exit(1)
		return
	}

	if plainOutput {
		fmt.Printf("Name:           %s\n", name)
		fmt.Printf("Title Font:     %s (%.0fpt)\n", st.TitleFont, st.TitleSize)
		fmt.Printf("Heading Font:   %s (%.0fpt)\n", st.HeadingFont, st.HeadingSize)
		fmt.Printf("Body Font:      %s (%.0fpt)\n", st.BodyFont, st.BodySize)
		fmt.Printf("Code Font:      %s (%.0fpt)\n", st.CodeFont, st.CodeSize)
		fmt.Printf("Text Color:     %s\n", st.TextColor)
		fmt.Printf("Accent Color:   %s\n", st.AccentColor)
		fmt.Printf("Page Margin:    %.2f in\n", st.PageMarginInches)
		return
	}

	result := AgentOutput{Success: true, Style: st}
	printResult(result, false, t)
}

func CreateTemplate(outputPath, presetName string, plainOutput bool, lang i18n.Lang) {
	t := func(k string) string { return i18n.T(lang, k) }
	st, err := style.LoadPresetOrDefault(presetName)
	if err != nil {
		result := AgentOutput{Success: false, Error: err.Error()}
		printResult(result, plainOutput, t)
		os.Exit(1)
		return
	}

	if err := style.SaveTemplateFile(outputPath, st); err != nil {
		result := AgentOutput{Success: false, Error: err.Error()}
		printResult(result, plainOutput, t)
		os.Exit(1)
		return
	}

	if plainOutput {
		fmt.Printf("%s: %s\n", t("cli_template_created"), outputPath)
	} else {
		result := AgentOutput{Success: true, OutputPath: outputPath}
		printResult(result, false, t)
	}
}

func printResult(result AgentOutput, plain bool, t func(string) string) {
	if plain {
		if result.Success {
			if result.OutputPath != "" {
				fmt.Printf("%s: %s (%d %s)\n", t("cli_ok"), result.OutputPath, result.Bytes, t("cli_bytes"))
			} else if result.Style != nil {
				fmt.Printf("%s\n", t("cli_ok"))
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s: %s\n", t("cli_error"), result.Error)
		}
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result)
}
