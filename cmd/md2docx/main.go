package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/md2docx/cli/internal/cli"
	"github.com/md2docx/cli/internal/config"
	"github.com/md2docx/cli/internal/i18n"
	"github.com/md2docx/cli/internal/style"
	"github.com/md2docx/cli/internal/tui"
)

var (
	jsonOutput bool
	langFlag   string

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
Run without subcommands to launch the interactive TUI.

Use --lang to set the interface language (en, zh-CN, ja, ko, es, pt-BR, de, fr).
The language is saved and remembered across sessions.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// If --lang was explicitly set, save it
			if cmd.Flags().Changed("lang") {
				if _, err := config.SetLanguage(i18n.Lang(langFlag)); err != nil {
					return err
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runTUI()
		},
	}
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "Interface language (en, zh-CN, ja, ko, es, pt-BR, de, fr)")

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

	presetsCmd := &cobra.Command{
		Use:   "presets",
		Short: "List available built-in style presets",
		Run:   runPresets,
	}
	presetsCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	presetCmd := &cobra.Command{
		Use:   "preset [name]",
		Short: "Show details of a specific preset",
		Args:  cobra.ExactArgs(1),
		Run:   runPreset,
	}
	presetCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

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

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("md2docx %s (commit %s, built %s)\n", version, commit, buildDate)
		},
	}

	rootCmd.AddCommand(convertCmd, presetsCmd, presetCmd, templateCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// resolveLang loads config and determines the effective language.
func resolveLang() (i18n.Lang, *config.Config) {
	cfg, err := config.Load()
	if err != nil || cfg.Lang == "" {
		return i18n.EN, cfg
	}
	return cfg.Lang, cfg
}

func runTUI() {
	lang, cfg := resolveLang()
	m := tui.NewModel(cfg, lang)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", i18n.T(lang, "err_tui"), err)
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
	lang, _ := resolveLang()

	cli.Convert(cli.ConvertOptions{
		InputPath:     input,
		OutputPath:    output,
		StyleRef:      styleRef,
		PlainOutput:   !jsonOutput,
		Mermaid:       mermaid,
		MermaidServer: mermaidServer,
		MermaidTheme:  mermaidTheme,
		Lang:          lang,
	})
}

func runPresets(cmd *cobra.Command, args []string) {
	lang, _ := resolveLang()
	cli.ListPresets(!jsonOutput, lang)
}

func runPreset(cmd *cobra.Command, args []string) {
	lang, _ := resolveLang()
	cli.ShowPreset(args[0], !jsonOutput, lang)
}

func runTemplateCreate(cmd *cobra.Command, args []string) {
	output, _ := cmd.Flags().GetString("output")
	styleRef, _ := cmd.Flags().GetString("style")
	lang, _ := resolveLang()

	if styleRef == "" {
		styleRef = style.PresetDefault
	}
	cli.CreateTemplate(output, styleRef, !jsonOutput, lang)
}
