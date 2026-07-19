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

	"github.com/md2docx/cli/internal/converter"
	"github.com/md2docx/cli/internal/style"
)

type step int

const (
	stepSelectInput step = iota
	stepSelectOutput
	stepSelectStyle
	stepSelectStyleFile // browsing for a custom JSON template file
	stepConfirm
	stepConverting
)

// lipgloss styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F3F4F6"))
)

// Model is the Bubble Tea model for the TUI.
type Model struct {
	step   step
	width  int
	height int

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
	styleSource    string // "preset:<name>" or "file:<path>"
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

// NewModel creates the initial TUI model.
func NewModel() Model {
	// Input file picker
	inputPicker := filepicker.New()
	inputPicker.CurrentDirectory, _ = os.Getwd()
	inputPicker.AllowedTypes = []string{".md", ".markdown", ".txt"}
	inputPicker.DirAllowed = true

	// Output file picker (but we also have a text input for naming)
	outputPicker := filepicker.New()
	outputPicker.CurrentDirectory, _ = os.Getwd()
	outputPicker.DirAllowed = true
	outputPicker.FileAllowed = false

	// Output text input
	oi := textinput.New()
	oi.Placeholder = "output.docx"
	oi.CharLimit = 256

	// Style file picker (for custom JSON)
	stylePicker := filepicker.New()
	stylePicker.CurrentDirectory, _ = os.Getwd()
	stylePicker.AllowedTypes = []string{".json"}

	return Model{
		step:         stepSelectInput,
		inputPicker:  inputPicker,
		outputPicker: outputPicker,
		outputInput:  oi,
		stylePicker:  stylePicker,
		stylePresets: style.AllPresetNames(),
		styleCursor:  len(style.AllPresetNames()) - 1, // default to "default" preset
		confirmCursor: 0,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.inputPicker.Init(),
		textinput.Blink,
	)
}

// Update implements tea.Model.
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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Esc goes back one step, or quits from the first step
			switch m.step {
			case stepSelectInput:
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
		case stepSelectInput:
			return m.updateStepInput(msg)
		case stepSelectOutput:
			return m.updateStepOutput(msg)
		case stepSelectStyle:
			return m.updateStepStyle(msg)
		case stepSelectStyleFile:
			return m.updateStepStyleFile(msg)
		case stepConfirm:
			return m.updateStepConfirm(msg)
		case stepConverting:
			if msg.String() == "q" || msg.String() == "enter" || msg.String() == "esc" {
				m.done = true
				return m, tea.Quit
			}
		}
	}

	// Delegate to active file picker
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

func (m Model) updateStepInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if path, err := m.inputPicker.DidSelectFile(msg); err == nil {
			m.inputPath = path
			m.step = stepSelectOutput
			base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			m.outputInput.SetValue(base + ".docx")
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateStepOutput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		filename := m.outputInput.Value()
		if filename != "" {
			dir := m.outputPicker.CurrentDirectory
			m.outputPath = filepath.Join(dir, filename)
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

func (m Model) updateStepStyle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
			m.styleSourceStr = fmt.Sprintf("Preset: %s", m.stylePresets[m.styleCursor])
			m.goToConfirm()
		} else if m.styleCursor == len(m.stylePresets) {
			m.styleSource = ""
			m.styleSourceStr = "Preset: default"
			m.goToConfirm()
		} else if m.styleCursor == len(m.stylePresets)+1 {
			// Transition to file picker for custom JSON template
			m.step = stepSelectStyleFile
			var cmd tea.Cmd
			m.stylePicker, cmd = m.stylePicker.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) updateStepStyleFile(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		if path, err := m.stylePicker.DidSelectFile(msg); err == nil {
			m.styleSource = "file:" + path
			m.styleSourceStr = fmt.Sprintf("File: %s", filepath.Base(path))
			m.goToConfirm()
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) goToConfirm() {
	m.step = stepConfirm
	mermaidStatus := "[ ] Render mermaid diagrams"
	if m.useMermaid {
		mermaidStatus = "[x] Render mermaid diagrams"
	}
	m.confirmItems = []string{
		fmt.Sprintf("Input:   %s", m.inputPath),
		fmt.Sprintf("Output:  %s", m.outputPath),
		fmt.Sprintf("Style:   %s", m.styleSourceStr),
		mermaidStatus,
		"",
		"> Convert",
		"  Back",
	}
	m.confirmCursor = 5 // "Convert"
}

func (m Model) updateStepConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
			mermaidStatus := "[ ] Render mermaid diagrams"
			if m.useMermaid {
				mermaidStatus = "[x] Render mermaid diagrams"
			}
			m.confirmItems[3] = mermaidStatus
		} else if m.confirmCursor == 5 {
			return m.startConversion()
		} else if m.confirmCursor == 4 {
			m.step = stepSelectStyle
			m.confirmCursor = 0
		}
	}
	return m, nil
}

func (m Model) startConversion() (tea.Model, tea.Cmd) {
	m.step = stepConverting

	var styleRef string
	if strings.HasPrefix(m.styleSource, "preset:") {
		styleRef = strings.TrimPrefix(m.styleSource, "preset:")
	} else if strings.HasPrefix(m.styleSource, "file:") {
		styleRef = strings.TrimPrefix(m.styleSource, "file:")
	}

	var convertOpts []converter.ConversionOption
	if m.useMermaid {
		convertOpts = append(convertOpts, converter.WithMermaid(&converter.MermaidInkRenderer{
			Theme: "default",
		}))
	}

	result, err := style.ResolveAndConvertWithOptions(m.inputPath, m.outputPath, styleRef, convertOpts...)
	if err != nil {
		m.err = err
	} else {
		m.result = fmt.Sprintf("%s (%d bytes)", result.OutputPath, result.Bytes)
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.err != nil {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %v\n\nPress esc to exit.", m.err)),
		)
	}

	if m.done {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			successStyle.Render(fmt.Sprintf("Done: %s\n\nPress any key to exit.", m.result)),
		)
	}

	switch m.step {
	case stepSelectInput:
		return m.viewInputPicker()
	case stepSelectOutput:
		return m.viewOutputPicker()
	case stepSelectStyle:
		return m.viewStylePicker()
	case stepSelectStyleFile:
		return m.viewStyleFilePicker()
	case stepConfirm:
		return m.viewConfirm()
	case stepConverting:
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			titleStyle.Render("Converting..."),
		)
	}
	return ""
}

func (m Model) viewInputPicker() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("md2docx – Select Markdown Input"),
		m.inputPicker.View(),
		helpStyle.Render("↑/↓ navigate • enter select • esc quit"),
	)
}

func (m Model) viewOutputPicker() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("md2docx – Choose Output"),
		labelStyle.Render("Filename: ")+m.outputInput.View(),
		"",
		m.outputPicker.View(),
		helpStyle.Render("↑/↓ navigate • enter confirm • tab toggle input • esc back"),
	)
}

func (m Model) viewStylePicker() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("md2docx – Choose Style"))
	b.WriteString("\n\n")

	for i, name := range m.stylePresets {
		cursor := "  "
		if i == m.styleCursor {
			cursor = "> "
		}
		desc := style.PresetDescriptions()[name]
		b.WriteString(fmt.Sprintf("%s%-20s %s\n", cursor, name, desc))
	}

	cursor := "  "
	if m.styleCursor == len(m.stylePresets) {
		cursor = "> "
	}
	b.WriteString(fmt.Sprintf("\n%sUse default\n", cursor))

	cursor = "  "
	if m.styleCursor == len(m.stylePresets)+1 {
		cursor = "> "
	}
	b.WriteString(fmt.Sprintf("%sCustom JSON file...\n", cursor))

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter select • esc back"))

	return b.String()
}

func (m Model) viewStyleFilePicker() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("md2docx – Select Style Template (JSON)"),
		m.stylePicker.View(),
		helpStyle.Render("↑/↓ navigate • enter select • esc back"),
	)
}

func (m Model) viewConfirm() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("md2docx – Confirm Conversion"))
	b.WriteString("\n\n")

	for i, item := range m.confirmItems {
		if i < 3 {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) == 2 {
				b.WriteString(labelStyle.Render(parts[0]+": "))
				b.WriteString(valueStyle.Render(strings.TrimSpace(parts[1])))
				b.WriteString("\n")
			} else {
				b.WriteString(item + "\n")
			}
		} else if item == "" {
			b.WriteString("\n")
		} else {
			cursor := "  "
			if i == m.confirmCursor {
				cursor = "> "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, strings.TrimSpace(item)))
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter select • esc back"))

	return b.String()
}
