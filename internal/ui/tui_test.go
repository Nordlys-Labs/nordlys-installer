package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	t.Parallel()

	m := NewModel()

	if m.state != StateToolSelect {
		t.Errorf("initial state = %v, want StateToolSelect", m.state)
	}

	if len(m.tools) != 7 {
		t.Errorf("tools count = %d, want 7", len(m.tools))
	}

	if m.selected == nil {
		t.Error("selected map should be initialized")
	}

	if m.model == "" {
		t.Error("model should have default value")
	}
}

func TestModel_Init(t *testing.T) {
	t.Parallel()

	m := NewModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init() should return a command for text input blink")
	}
}

func TestModel_View(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state State
	}{
		{"StateToolSelect", StateToolSelect},
		{"StateAPIKeyInput", StateAPIKeyInput},
		{"StateInstalling", StateInstalling},
		{"StateComplete", StateComplete},
		{"StateError", StateError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := NewModel()
			m.state = tt.state

			if tt.state == StateError {
				m.err = errTestError
			}
			if tt.state == StateComplete {
				m.installedDone = []string{"claude-code"}
			}

			view := m.View()
			if view == "" {
				t.Errorf("View() for state %v should not be empty", tt.state)
			}
		})
	}
}

var errTestError = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test error" }

func TestModel_Update_ToolSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		key       string
		wantState State
	}{
		{"quit with q", "q", StateToolSelect},
		{"quit with ctrl+c", "ctrl+c", StateToolSelect},
		{"move up", "up", StateToolSelect},
		{"move down", "down", StateToolSelect},
		{"move with k", "k", StateToolSelect},
		{"move with j", "j", StateToolSelect},
		{"space to select", " ", StateToolSelect},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := NewModel()
			m.state = StateToolSelect

			var msg tea.KeyMsg
			switch tt.key {
			case "ctrl+c":
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			case "up":
				msg = tea.KeyMsg{Type: tea.KeyUp}
			case "down":
				msg = tea.KeyMsg{Type: tea.KeyDown}
			default:
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}

			newModel, _ := m.Update(msg)
			if newModel == nil {
				t.Error("Update() should return a model")
			}
		})
	}
}

func TestModel_Update_APIKeyInput(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateAPIKeyInput

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := m.Update(msg)

	if newModel == nil {
		t.Error("Update() should return a model")
	}
	if cmd == nil {
		t.Error("Update() with ctrl+c should return quit command")
	}
}

func TestModel_Update_Complete(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateComplete

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	newModel, cmd := m.Update(msg)

	if newModel == nil {
		t.Error("Update() should return a model")
	}
	if cmd == nil {
		t.Error("Update() with q should return quit command")
	}
}

func TestModel_Update_Error(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateError
	m.err = errTestError

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	newModel, cmd := m.Update(msg)

	if newModel == nil {
		t.Error("Update() should return a model")
	}
	if cmd == nil {
		t.Error("Update() with q should return quit command")
	}
}

func TestModel_Update_InstallMsg(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateInstalling
	m.selected = map[int]bool{0: true}

	successMsg := installMsg{tool: "claude-code", err: nil}
	newModel, _ := m.Update(successMsg)

	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel)
	}
	if len(updatedModel.installedDone) != 1 {
		t.Error("successful install should add tool to installedDone")
	}

	m2 := NewModel()
	m2.state = StateInstalling

	errorMsg := installMsg{tool: "claude-code", err: errTestError}
	newModel2, _ := m2.Update(errorMsg)

	updatedModel2, ok := newModel2.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel2)
	}
	if updatedModel2.state != StateError {
		t.Error("failed install should set state to StateError")
	}
}

func TestStyles(t *testing.T) {
	t.Parallel()

	styles := []struct {
		style interface{}
		name  string
	}{
		{TitleStyle, "TitleStyle"},
		{SelectedStyle, "SelectedStyle"},
		{UnselectedStyle, "UnselectedStyle"},
		{ErrorStyle, "ErrorStyle"},
		{SuccessStyle, "SuccessStyle"},
		{InfoStyle, "InfoStyle"},
		{PromptStyle, "PromptStyle"},
		{BorderStyle, "BorderStyle"},
	}

	for _, tt := range styles {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.style == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
		})
	}
}

func TestViewToolSelect(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateToolSelect
	m.cursor = 1
	m.selected[0] = true
	m.selected[2] = true

	view := m.viewToolSelect()
	if view == "" {
		t.Error("viewToolSelect() should not return empty string")
	}
}

func TestViewAPIKeyInput(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateAPIKeyInput

	view := m.viewAPIKeyInput()
	if view == "" {
		t.Error("viewAPIKeyInput() should not return empty string")
	}
}

func TestViewInstalling(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateInstalling
	m.selected[0] = true
	m.selected[1] = true
	m.installedDone = []string{"claude-code"}
	m.currentTool = "opencode"

	view := m.viewInstalling()
	if view == "" {
		t.Error("viewInstalling() should not return empty string")
	}
}

func TestViewComplete(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateComplete
	m.installedDone = []string{"claude-code", "opencode"}

	view := m.viewComplete()
	if view == "" {
		t.Error("viewComplete() should not return empty string")
	}
}

func TestViewError(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateError
	m.err = errTestError

	view := m.viewError()
	if view == "" {
		t.Error("viewError() should not return empty string")
	}
}

func TestUpdateToolSelect_Navigation(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateToolSelect
	m.cursor = 3

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := m.Update(upMsg)
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel)
	}
	if updatedModel.cursor != 2 {
		t.Errorf("cursor after up = %d, want 2", updatedModel.cursor)
	}

	m2 := NewModel()
	m2.state = StateToolSelect
	m2.cursor = 0

	upMsg2 := tea.KeyMsg{Type: tea.KeyUp}
	newModel2, _ := m2.Update(upMsg2)
	updatedModel2, ok := newModel2.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel2)
	}
	if updatedModel2.cursor != 0 {
		t.Errorf("cursor should not go below 0, got %d", updatedModel2.cursor)
	}

	m3 := NewModel()
	m3.state = StateToolSelect
	m3.cursor = 5

	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel3, _ := m3.Update(downMsg)
	updatedModel3, ok := newModel3.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel3)
	}
	if updatedModel3.cursor != 6 {
		t.Errorf("cursor after down = %d, want 6", updatedModel3.cursor)
	}

	m4 := NewModel()
	m4.state = StateToolSelect
	m4.cursor = len(m4.tools) - 1

	downMsg2 := tea.KeyMsg{Type: tea.KeyDown}
	newModel4, _ := m4.Update(downMsg2)
	updatedModel4, ok := newModel4.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel4)
	}
	if updatedModel4.cursor != len(m4.tools)-1 {
		t.Errorf("cursor should not exceed max, got %d", updatedModel4.cursor)
	}
}

func TestUpdateToolSelect_Selection(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateToolSelect
	m.cursor = 2

	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	newModel, _ := m.Update(spaceMsg)
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel)
	}

	if !updatedModel.selected[2] {
		t.Error("space should toggle selection on")
	}

	newModel2, _ := updatedModel.Update(spaceMsg)
	updatedModel2, ok := newModel2.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel2)
	}

	if updatedModel2.selected[2] {
		t.Error("space should toggle selection off")
	}
}

func TestUpdateToolSelect_Enter(t *testing.T) {
	t.Parallel()

	m := NewModel()
	m.state = StateToolSelect

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(enterMsg)
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel)
	}

	if updatedModel.state != StateToolSelect {
		t.Error("enter with no selection should stay in StateToolSelect")
	}

	m2 := NewModel()
	m2.state = StateToolSelect
	m2.selected[0] = true

	newModel2, _ := m2.Update(enterMsg)
	updatedModel2, ok := newModel2.(Model)
	if !ok {
		t.Fatalf("Update should return Model, got %T", newModel2)
	}

	if updatedModel2.state != StateAPIKeyInput {
		t.Errorf("enter with selection should go to StateAPIKeyInput, got %v", updatedModel2.state)
	}
}
