package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/md2docx/cli/internal/converter"
	"github.com/md2docx/cli/internal/style"
)

// AgentOutput is structured output for agent consumption.
type AgentOutput struct {
	Success    bool                    `json:"success"`
	OutputPath string                  `json:"outputPath,omitempty"`
	Bytes      int64                   `json:"bytes,omitempty"`
	Error      string                  `json:"error,omitempty"`
	Style      *converter.StyleTemplate `json:"style,omitempty"`
	Presets    []string                `json:"presets,omitempty"`
}

// Convert converts markdown to DOCX and prints structured JSON output.
func Convert(inputPath, outputPath, styleRef string, plainOutput bool) {
	var result AgentOutput

	// Ensure .docx extension
	if filepath.Ext(outputPath) != ".docx" {
		outputPath += ".docx"
	}

	cr, err := style.ResolveAndConvert(inputPath, outputPath, styleRef)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		printResult(result, plainOutput)
		os.Exit(1)
		return
	}

	result.Success = true
	result.OutputPath = cr.OutputPath
	result.Bytes = cr.Bytes
	printResult(result, plainOutput)
}

// ListPresets prints all available built-in style presets.
func ListPresets(plainOutput bool) {
	names := style.AllPresetNames()
	descs := style.PresetDescriptions()

	if plainOutput {
		for _, name := range names {
			fmt.Printf("%-20s %s\n", name, descs[name])
		}
		return
	}

	result := AgentOutput{
		Success: true,
		Presets: names,
	}
	printResult(result, false)
}

// ShowPreset displays the details of a specific preset.
func ShowPreset(name string, plainOutput bool) {
	st, err := style.LoadPreset(name)
	if err != nil {
		result := AgentOutput{Success: false, Error: err.Error()}
		printResult(result, plainOutput)
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
	printResult(result, false)
}

// CreateTemplate creates a new style template JSON file from a preset.
func CreateTemplate(outputPath, presetName string, plainOutput bool) {
	st, err := style.LoadPresetOrDefault(presetName)
	if err != nil {
		result := AgentOutput{Success: false, Error: err.Error()}
		printResult(result, plainOutput)
		os.Exit(1)
		return
	}

	if err := style.SaveTemplateFile(outputPath, st); err != nil {
		result := AgentOutput{Success: false, Error: err.Error()}
		printResult(result, plainOutput)
		os.Exit(1)
		return
	}

	if plainOutput {
		fmt.Printf("Style template created: %s\n", outputPath)
	} else {
		result := AgentOutput{Success: true, OutputPath: outputPath}
		printResult(result, false)
	}
}

func printResult(result AgentOutput, plain bool) {
	if plain {
		if result.Success {
			if result.OutputPath != "" {
				fmt.Printf("OK: %s (%d bytes)\n", result.OutputPath, result.Bytes)
			} else if result.Style != nil {
				fmt.Printf("OK\n")
			}
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", result.Error)
		}
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result)
}
