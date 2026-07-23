package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/md2docx/cli/internal/config"
	"github.com/md2docx/cli/internal/converter"
	"github.com/md2docx/cli/internal/i18n"
	"github.com/md2docx/cli/internal/style"
)

type step int

const (
	stepSelectLang step = iota
	stepSelectInput
	stepSelectOutput
	stepSelectStyle
	stepSelectStyleFile
	stepConfirm
	stepConverting
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).MarginBottom(1)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).MarginTop(1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	valueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))
)

type Model struct {
	step   step
	width  int
	height int

	// Localization
	lang i18n.Lang
	cfg  *config.Config

	// Language selection
	langList   []i18n.Lang
	langCursor int

	// Input file
	inputPicker filepicker.Model
	inputPath   string

	// Output file
	outputPicker filepicker.Model
	outputInput  textinput.Model
	outputPath   string

	// Style selector
	stylePicker    filepicker.Model
	stylePresets   []string
	styleCursor    int
	styleSource    string
	styleSourceStr string

	// Mermaid
	useMermaid bool

	// Confirm
	confirmCursor int
	confirmItems  []string

	// Result
	err    error
	result string
	done   bool
}

type LangSelectedMsg struct{}

func NewModel(cfg *config.Config, lang i18n.Lang) Model {
	inputPicker := filepicker.New()
	inputPicker.CurrentDirectory, _ = os.Getwd()
	inputPicker.AllowedTypes = []string{".md", ".markdown", ".txt"}
	inputPicker.DirAllowed = true

	outputPicker := filepicker.New()
	outputPicker.CurrentDirectory, _ = os.Getwd()
	outputPicker.DirAllowed = true
	outputPicker.FileAllowed = false

	oi := textinput.New()
	oi.Placeholder = "output.docx"
	oi.CharLimit = 256

	stylePicker := filepicker.New()
	stylePicker.CurrentDirectory, _ = os.Getwd()
	stylePicker.AllowedTypes = []string{".json"}

	// Initial step: language selection if first run, otherwise input
	initialStep := stepSelectInput
	if cfg.FirstRun {
		initialStep = stepSelectLang
	}

	return Model{
		step:        initialStep,
		lang:        lang,
		cfg:         cfg,
		langList:    i18n.AllLangs(),
		langCursor:  0,

		inputPicker:  inputPicker,
		outputPicker: outputPicker,
		outputInput:  oi,
		stylePicker:  stylePicker,
		stylePresets: style.AllPresetNames(),

		// Default style cursor to the language-appropriate preset
		styleCursor:  presetCursorForLang(lang),
		confirmCursor: 0,
	}
}

func presetCursorForLang(lang i18n.Lang) int {
	ds := i18n.DefaultStyleForLang(lang)
	presets := style.AllPresetNames()
	for i, p := range presets {
		if p == ds {
			return i
		}
	}
	return len(presets) - 1 // default
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.inputPicker.Init(), textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.inputPicker.Height = msg.Height - 10
		m.outputPicker.Height = msg.Height - 10
		m.stylePicker.Height = msg.Height - 10
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "esc" {
			// Esc goes back one step, or quits from the first step
			switch m.step {
			case stepSelectLang, stepSelectInput:
				return m, tea.Quit
			case stepSelectOutput:
				m.step = stepSelectInput
				return m, nil
			case stepSelectStyle, stepSelectStyleFile:
				m.step = stepSelectOutput
				return m, nil
			case stepConfirm:
				m.step = stepSelectStyle
				return m, nil
			case stepConverting:
				if m.done {
					return m, tea.Quit
				}
			}
			return m, nil
		}

		switch m.step {
		case stepSelectLang:
			return m.updateLang(msg)
		case stepSelectInput:
			return m.updateInput(msg)
		case stepSelectOutput:
			return m.updateOutput(msg)
		case stepSelectStyle:
			return m.updateStyle(msg)
		case stepSelectStyleFile:
			return m.updateStyleFile(msg)
		case stepConfirm:
			return m.updateConfirm(msg)
		case stepConverting:
			if msg.String() == "q" || msg.String() == "enter" || msg.String() == "esc" {
				m.done = true
				return m, tea.Quit
			}
		}
	}

	// Delegate to active pickers
	switch m.step {
	case stepSelectInput:
		var cmd tea.Cmd
		m.inputPicker, cmd = m.inputPicker.Update(msg)
		return m, cmd
	case stepSelectOutput:
		var cmd tea.Cmd
		m.outputPicker, cmd = m.outputPicker.Update(msg)
		m.outputInput, _ = m.outputInput.Update(msg)
		return m, cmd
	case stepSelectStyleFile:
		var cmd tea.Cmd
		m.stylePicker, cmd = m.stylePicker.Update(msg)
		return m, cmd
	}

	return m, nil
}

// ---- Language selection ----

func (m Model) updateLang(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.langCursor > 0 {
			m.langCursor--
		}
	case "down", "j":
		if m.langCursor < len(m.langList)-1 {
			m.langCursor++
		}
	case "enter", " ":
		selected := m.langList[m.langCursor]
		cfg, err := config.SetLanguage(selected)
		if err != nil {
			m.err = fmt.Errorf("saving language: %w", err)
			return m, nil
		}
		m.lang = selected
		m.cfg = cfg
		m.styleCursor = presetCursorForLang(selected)
		m.step = stepSelectInput
		return m, nil
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) viewLang() string {
	t := func(k string) string { return i18n.T(m.lang, k) }
	var b strings.Builder
	b.WriteString(titleStyle.Render(t("lang_select_title")))
	b.WriteString("\n\n")
	for i, l := range m.langList {
		cursor := "  "
		if i == m.langCursor {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, i18n.LangName(l)))
	}
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(t("lang_select_help")))
	return b.String()
}

// ---- Input step ----

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		if selected, path := m.inputPicker.DidSelectFile(msg); selected {
			m.inputPath = path
			m.step = stepSelectOutput
			base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			m.outputInput.SetValue(base + ".docx")
			return m, nil
		}
	}
	return m, nil
}

// ---- Output step ----

func (m Model) updateOutput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if fn := m.outputInput.Value(); fn != "" {
			m.outputPath = filepath.Join(m.outputPicker.CurrentDirectory, fn)
			if !strings.HasSuffix(m.outputPath, ".docx") {
				m.outputPath += ".docx"
			}
			m.step = stepSelectStyle
			return m, nil
		}
	case "tab":
		if m.outputInput.Focused() {
			m.outputInput.Blur()
		} else {
			m.outputInput.Focus()
		}
	}
	return m, nil
}

// ---- Style step ----

func (m Model) updateStyle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.styleCursor > 0 {
			m.styleCursor--
		}
	case "down", "j":
		if m.styleCursor < len(m.stylePresets)+1 {
			m.styleCursor++
		}
	case "enter":
		if m.styleCursor < len(m.stylePresets) {
			m.styleSource = "preset:" + m.stylePresets[m.styleCursor]
			m.styleSourceStr = "Preset: " + m.stylePresets[m.styleCursor]
			m.goToConfirm()
		} else if m.styleCursor == len(m.stylePresets) {
			m.styleSource = ""
			m.styleSourceStr = "Preset: default"
			m.goToConfirm()
		} else {
			m.step = stepSelectStyleFile
			var cmd tea.Cmd
			m.stylePicker, cmd = m.stylePicker.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) updateStyleFile(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		if selected, path := m.stylePicker.DidSelectFile(msg); selected {
			m.styleSource = "file:" + path
			m.styleSourceStr = "File: " + filepath.Base(path)
			m.goToConfirm()
			return m, nil
		}
	}
	return m, nil
}

// ---- Confirm step ----

func (m *Model) goToConfirm() {
	t := func(k string) string { return i18n.T(m.lang, k) }
	mermaidTxt := "[ ] " + t("tui_confirm_mermaid")
	if m.useMermaid {
		mermaidTxt = "[x] " + t("tui_confirm_mermaid")
	}
	m.step = stepConfirm
	m.confirmItems = []string{
		fmt.Sprintf("%s:  %s", t("tui_confirm_input"), m.inputPath),
		fmt.Sprintf("%s: %s", t("tui_confirm_output"), m.outputPath),
		fmt.Sprintf("%s:  %s", t("tui_confirm_style"), m.styleSourceStr),
		mermaidTxt,
		"",
		"> " + t("tui_confirm_convert"),
		"  " + t("tui_confirm_back"),
	}
	m.confirmCursor = 4 // "Convert"
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	t := func(k string) string { return i18n.T(m.lang, k) }
	switch msg.String() {
	case "up", "k":
		if m.confirmCursor > 3 {
			m.confirmCursor--
		}
	case "down", "j":
		if m.confirmCursor < 5 {
			m.confirmCursor++
		}
	case " ", "enter":
		if m.confirmCursor == 3 {
			// Toggle mermaid
			m.useMermaid = !m.useMermaid
			txt := "[ ] " + t("tui_confirm_mermaid")
			if m.useMermaid {
				txt = "[x] " + t("tui_confirm_mermaid")
			}
			m.confirmItems[3] = txt
		} else if m.confirmCursor == 4 {
			return m.startConversion()
		} else if m.confirmCursor == 5 {
			m.step = stepSelectStyle
			m.confirmCursor = 0
		}
	}
	return m, nil
}

// ---- Conversion ----

func (m Model) startConversion() (tea.Model, tea.Cmd) {
	m.step = stepConverting
	var styleRef string
	if strings.HasPrefix(m.styleSource, "preset:") {
		styleRef = strings.TrimPrefix(m.styleSource, "preset:")
	} else if strings.HasPrefix(m.styleSource, "file:") {
		styleRef = strings.TrimPrefix(m.styleSource, "file:")
	}
	var opts []converter.ConversionOption
	if m.useMermaid {
		opts = append(opts, converter.WithMermaid(&converter.MermaidInkRenderer{Theme: "default"}))
	}
	result, err := style.ResolveAndConvertWithOptions(m.inputPath, m.outputPath, styleRef, opts...)
	if err != nil {
		m.err = err
	} else {
		m.result = fmt.Sprintf("%s (%d bytes)", result.OutputPath, result.Bytes)
	}
	return m, nil
}

// ---- Views ----

func (m Model) View() string {
	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf(i18n.T(m.lang, "tui_error_msg"), m.err)))
	}
	if m.done {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			successStyle.Render(fmt.Sprintf(i18n.T(m.lang, "tui_done_msg"), m.result)))
	}

	switch m.step {
	case stepSelectLang:
		return m.viewLang()
	case stepSelectInput:
		return m.viewInput()
	case stepSelectOutput:
		return m.viewOutput()
	case stepSelectStyle:
		return m.viewStyleList()
	case stepSelectStyleFile:
		return m.viewStyleFile()
	case stepConfirm:
		return m.viewConfirm()
	case stepConverting:
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			titleStyle.Render(i18n.T(m.lang, "tui_converting")))
	}
	return ""
}

func (m Model) viewInput() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(i18n.T(m.lang, "tui_input_title")),
		m.inputPicker.View(),
		helpStyle.Render(i18n.T(m.lang, "tui_nav_quit")),
	)
}

func (m Model) viewOutput() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(i18n.T(m.lang, "tui_output_title")),
		labelStyle.Render(i18n.T(m.lang, "tui_filename_label")+": ")+m.outputInput.View(),
		"",
		m.outputPicker.View(),
		helpStyle.Render(i18n.T(m.lang, "tui_nav_tab")),
	)
}

func (m Model) viewStyleList() string {
	t := func(k string) string { return i18n.T(m.lang, k) }
	var b strings.Builder
	b.WriteString(titleStyle.Render(t("tui_style_title")))
	b.WriteString("\n\n")
	for i, name := range m.stylePresets {
		cursor := "  "
		if i == m.styleCursor {
			cursor = "> "
		}
		desc := i18n.PresetDescription(m.lang, name)
		b.WriteString(fmt.Sprintf("%s%-20s %s\n", cursor, name, desc))
	}
	cursor := "  "
	if m.styleCursor == len(m.stylePresets) {
		cursor = "> "
	}
	b.WriteString(fmt.Sprintf("\n%s%s\n", cursor, t("tui_style_default")))
	cursor = "  "
	if m.styleCursor == len(m.stylePresets)+1 {
		cursor = "> "
	}
	b.WriteString(fmt.Sprintf("%s%s\n", cursor, t("tui_style_custom")))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(t("tui_nav_help")))
	return b.String()
}

func (m Model) viewStyleFile() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(i18n.T(m.lang, "tui_style_file_title")),
		m.stylePicker.View(),
		helpStyle.Render(i18n.T(m.lang, "tui_nav_help")),
	)
}

func (m Model) viewConfirm() string {
	t := func(k string) string { return i18n.T(m.lang, k) }
	var b strings.Builder
	b.WriteString(titleStyle.Render(t("tui_confirm_title")))
	b.WriteString("\n\n")

	for i, item := range m.confirmItems {
		switch {
		case i < 3:
			// Info lines: Input, Output, Style
			parts := strings.SplitN(item, ":", 2)
			if len(parts) == 2 {
				b.WriteString(labelStyle.Render(parts[0]+": "))
				b.WriteString(valueStyle.Render(strings.TrimSpace(parts[1])))
				b.WriteString("\n")
			} else {
				b.WriteString(item + "\n")
			}
		case i == 3:
			// Mermaid toggle (selectable)
			cursor := "  "
			if i == m.confirmCursor {
				cursor = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, item))
		case item == "":
			b.WriteString("\n")
		default:
			// Action buttons: Convert, Back
			cursor := "  "
			if i == m.confirmCursor+1 { // +1 offset for the empty separator
				cursor = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, strings.TrimSpace(item)))
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render(t("tui_nav_help")))
	return b.String()
}
