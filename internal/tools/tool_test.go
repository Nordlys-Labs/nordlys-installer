package tools

import (
	"testing"
)

func TestToolInterface(t *testing.T) {
	t.Parallel()

	allTools := GetAllTools()

	if len(allTools) != 7 {
		t.Errorf("GetAllTools() returned %d tools, want 7", len(allTools))
	}

	for _, tool := range allTools {
		t.Run(tool.Name(), func(t *testing.T) {
			t.Parallel()

			if tool.Name() == "" {
				t.Error("Tool.Name() should not be empty")
			}

			if tool.Description() == "" {
				t.Error("Tool.Description() should not be empty")
			}

			path, err := tool.ConfigPath()
			if err != nil {
				t.Errorf("Tool.ConfigPath() error = %v", err)
			}
			if path == "" {
				t.Error("Tool.ConfigPath() should not be empty")
			}
		})
	}
}

func TestAllTools_ValidateNoConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	tools := []Tool{
		NewClaudeCode(tmpDir),
		NewOpenCode(tmpDir),
		NewCodex(tmpDir),
		NewGeminiCLI(tmpDir),
		NewGrokCLI(tmpDir),
		NewQwenCode(tmpDir),
		NewZed(tmpDir),
	}

	for _, tool := range tools {
		t.Run(tool.Name(), func(t *testing.T) {
			t.Parallel()
			_ = tool.Validate()
		})
	}
}

func TestAllTools_UninstallNoConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	tools := []Tool{
		NewClaudeCode(tmpDir),
		NewOpenCode(tmpDir),
		NewCodex(tmpDir),
		NewGeminiCLI(tmpDir),
		NewGrokCLI(tmpDir),
		NewQwenCode(tmpDir),
		NewZed(tmpDir),
	}

	for _, tool := range tools {
		t.Run(tool.Name(), func(t *testing.T) {
			t.Parallel()
			_ = tool.Uninstall()
		})
	}
}

func TestGetToolByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		toolName string
		wantNil  bool
	}{
		{
			name:     "claude-code exists",
			toolName: "claude-code",
			wantNil:  false,
		},
		{
			name:     "opencode exists",
			toolName: "opencode",
			wantNil:  false,
		},
		{
			name:     "unknown tool",
			toolName: "unknown-tool",
			wantNil:  true,
		},
		{
			name:     "empty name",
			toolName: "",
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetToolByName(tt.toolName)
			if (got == nil) != tt.wantNil {
				t.Errorf("GetToolByName(%q) nil = %v, wantNil %v", tt.toolName, got == nil, tt.wantNil)
			}

			if !tt.wantNil && got.Name() != tt.toolName {
				t.Errorf("GetToolByName(%q).Name() = %q", tt.toolName, got.Name())
			}
		})
	}
}
