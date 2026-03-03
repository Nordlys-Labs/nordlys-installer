package ui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
	"github.com/nordlys-labs/nordlys-installer/internal/runtime"
	"github.com/nordlys-labs/nordlys-installer/internal/tools"
)

type State int

const (
	StateToolSelect State = iota
	StateAPIKeyInput
	StateInstalling
	StateComplete
	StateError
)

type Model struct {
	state         State
	tools         []tools.Tool
	selected      map[int]bool
	cursor        int
	apiKeyInput   textinput.Model
	apiKey        string
	model         string
	err           error
	installedDone []string
	currentTool   string
}

type installMsg struct {
	tool string
	err  error
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your Nordlys API key"
	ti.EchoMode = textinput.EchoPassword
	ti.Focus()
	ti.Width = 50

	return Model{
		state:       StateToolSelect,
		tools:       tools.GetAllTools(),
		selected:    make(map[int]bool),
		apiKeyInput: ti,
		model:       constants.DefaultModel,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateToolSelect:
			return m.updateToolSelect(msg)
		case StateAPIKeyInput:
			return m.updateAPIKeyInput(msg)
		case StateComplete, StateError:
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}

	case installMsg:
		if msg.err != nil {
			m.err = fmt.Errorf("failed to configure %s: %w", msg.tool, msg.err)
			m.state = StateError
			return m, nil
		}
		m.installedDone = append(m.installedDone, msg.tool)
		return m, m.installNext()
	}

	return m, nil
}

func (m Model) updateToolSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.tools)-1 {
			m.cursor++
		}
	case " ":
		m.selected[m.cursor] = !m.selected[m.cursor]
	case "enter":
		if len(m.selected) == 0 {
			return m, nil
		}
		m.state = StateAPIKeyInput
		return m, nil
	}
	return m, nil
}

func (m Model) updateAPIKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.apiKey = m.apiKeyInput.Value()
		if m.apiKey == "" {
			return m, nil
		}
		m.state = StateInstalling
		return m, m.installNext()
	}

	var cmd tea.Cmd
	m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
	return m, cmd
}

func (m Model) installNext() tea.Cmd {
	for i, tool := range m.tools {
		if !m.selected[i] {
			continue
		}

		alreadyDone := slices.Contains(m.installedDone, tool.Name())
		if alreadyDone {
			continue
		}

		m.currentTool = tool.Name()

		if tool.RequiresNode() {
			if err := runtime.EnsureNodeJS(); err != nil {
				return func() tea.Msg {
					return installMsg{tool: tool.Name(), err: err}
				}
			}
		}

		return func() tea.Msg {
			err := tool.UpdateConfig(m.apiKey, m.model, constants.APIBaseURL)
			return installMsg{tool: tool.Name(), err: err}
		}
	}

	m.state = StateComplete
	return nil
}

func (m Model) View() string {
	switch m.state {
	case StateToolSelect:
		return m.viewToolSelect()
	case StateAPIKeyInput:
		return m.viewAPIKeyInput()
	case StateInstalling:
		return m.viewInstalling()
	case StateComplete:
		return m.viewComplete()
	case StateError:
		return m.viewError()
	}
	return ""
}

func (m Model) viewToolSelect() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("Nordlys Installer"))
	b.WriteString("\n\n")
	b.WriteString("Select tools to configure with Nordlys:\n\n")

	for i, tool := range m.tools {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected[i] {
			checked = "x"
		}

		line := fmt.Sprintf("%s [%s] %s - %s", cursor, checked, tool.Name(), tool.Description())
		if m.cursor == i {
			line = SelectedStyle.Render(line)
		} else {
			line = UnselectedStyle.Render(line)
		}
		b.WriteString(line + "\n")
	}

	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("↑/↓: Navigate  Space: Select  Enter: Continue  Q: Quit"))

	return BorderStyle.Render(b.String())
}

func (m Model) viewAPIKeyInput() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("API Key Configuration"))
	b.WriteString("\n\n")
	b.WriteString(PromptStyle.Render(fmt.Sprintf("Get your API key from: %s", constants.APIKeyURL)))
	b.WriteString("\n\n")
	b.WriteString(m.apiKeyInput.View())
	b.WriteString("\n\n")
	b.WriteString(InfoStyle.Render("Enter: Continue  Ctrl+C: Quit"))

	return BorderStyle.Render(b.String())
}

func (m Model) viewInstalling() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("Installing..."))
	b.WriteString("\n\n")

	for i, tool := range m.tools {
		if !m.selected[i] {
			continue
		}

		done := slices.Contains(m.installedDone, tool.Name())

		status := "waiting"
		if done {
			status = "done"
		} else if tool.Name() == m.currentTool {
			status = "configuring"
		}

		b.WriteString(fmt.Sprintf("%s %s\n", status, tool.Name()))
	}

	return BorderStyle.Render(b.String())
}

func (m Model) viewComplete() string {
	var b strings.Builder

	b.WriteString(SuccessStyle.Render("Installation complete"))
	b.WriteString("\n\n")
	b.WriteString("Configured tools:\n\n")

	for _, toolName := range m.installedDone {
		b.WriteString(fmt.Sprintf("  - %s\n", toolName))
	}

	b.WriteString("\n")
	b.WriteString(InfoStyle.Render("Press Q to quit"))

	return BorderStyle.Render(b.String())
}

func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(ErrorStyle.Render("Error"))
	b.WriteString("\n\n")
	b.WriteString(m.err.Error())
	b.WriteString("\n\n")
	b.WriteString(InfoStyle.Render("Press Q to quit"))

	return BorderStyle.Render(b.String())
}

func RunTUI() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}
