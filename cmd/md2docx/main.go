package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/md2docx/cli/internal/cli"
	"github.com/md2docx/cli/internal/style"
	"github.com/md2docx/cli/internal/tui"
)

var (
	// Global flags
	jsonOutput bool

	// version info (set via ldflags)
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "md2docx",
		Short: "Convert Markdown to professional DOCX documents",
		Long: `md2docx converts Markdown files to DOCX (Open XML) documents.
No Word, Pandoc, or LibreOffice required.

Built-in style presets for US, CN, JP, EU, and KR document standards.
Supports rendering Mermaid diagrams as embedded PNG images.
Run without subcommands to launch the interactive TUI.`,
		Run: func(cmd *cobra.Command, args []string) {
			runTUI()
		},
	}

	// convert subcommand
	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert a Markdown file to DOCX",
		Long: `Convert a Markdown file to a DOCX document.

Style can be a built-in preset name (e.g., cn-official, us-business)
or a path to a JSON style template file.

Mermaid code blocks (~~~mermaid) are rendered as embedded PNG images
when --mermaid is enabled. Uses the public mermaid.ink service by default.

Examples:
  md2docx convert -i notes.md -o notes.docx
  md2docx convert -i report.md -o report.docx -s cn-official --mermaid
  md2docx convert -i doc.md -o doc.docx -s my-style.json --mermaid --mermaid-theme dark`,
		Run: runConvert,
	}
	convertCmd.Flags().StringP("input", "i", "", "Input Markdown file (required)")
	convertCmd.Flags().StringP("output", "o", "", "Output DOCX file (required)")
	convertCmd.Flags().StringP("style", "s", "", "Style preset name or template JSON path")
	convertCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output structured JSON (for agent/automation)")
	convertCmd.Flags().Bool("mermaid", false, "Render mermaid code blocks as embedded diagrams (via mermaid.ink)")
	convertCmd.Flags().String("mermaid-server", "", "Custom mermaid.ink server URL (default: https://mermaid.ink)")
	convertCmd.Flags().String("mermaid-theme", "default", "Mermaid theme: default, neutral, dark, forest")
	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("output")

	// presets subcommand
	presetsCmd := &cobra.Command{
		Use:   "presets",
		Short: "List available built-in style presets",
		Run:   runPresets,
	}
	presetsCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	// preset subcommand (single preset details)
	presetCmd := &cobra.Command{
		Use:   "preset [name]",
		Short: "Show details of a specific preset",
		Args:  cobra.ExactArgs(1),
		Run:   runPreset,
	}
	presetCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	// template subcommand
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Style template management",
	}

	templateCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a style template JSON file from a preset",
		Example: `  md2docx template create -o my-style.json
  md2docx template create -o cn-style.json -s cn-official`,
		Run: runTemplateCreate,
	}
	templateCreateCmd.Flags().StringP("output", "o", "", "Output JSON file path (required)")
	templateCreateCmd.Flags().StringP("style", "s", "", "Preset to base template on (default: default)")
	templateCreateCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	templateCreateCmd.MarkFlagRequired("output")

	templateCmd.AddCommand(templateCreateCmd)

	// version subcommand
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("md2docx %s (commit %s, built %s)\n", version, commit, buildDate)
		},
	}

	rootCmd.AddCommand(convertCmd, presetsCmd, presetCmd, templateCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTUI() {
	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

func runConvert(cmd *cobra.Command, args []string) {
	input, _ := cmd.Flags().GetString("input")
	output, _ := cmd.Flags().GetString("output")
	styleRef, _ := cmd.Flags().GetString("style")
	mermaid, _ := cmd.Flags().GetBool("mermaid")
	mermaidServer, _ := cmd.Flags().GetString("mermaid-server")
	mermaidTheme, _ := cmd.Flags().GetString("mermaid-theme")

	cli.Convert(cli.ConvertOptions{
		InputPath:     input,
		OutputPath:    output,
		StyleRef:      styleRef,
		PlainOutput:   !jsonOutput,
		Mermaid:       mermaid,
		MermaidServer: mermaidServer,
		MermaidTheme:  mermaidTheme,
	})
}

func runPresets(cmd *cobra.Command, args []string) {
	cli.ListPresets(!jsonOutput)
}

func runPreset(cmd *cobra.Command, args []string) {
	cli.ShowPreset(args[0], !jsonOutput)
}

func runTemplateCreate(cmd *cobra.Command, args []string) {
	output, _ := cmd.Flags().GetString("output")
	styleRef, _ := cmd.Flags().GetString("style")

	if styleRef == "" {
		styleRef = style.PresetDefault
	}

	cli.CreateTemplate(output, styleRef, !jsonOutput)
}
