package tools

import (
	"os"
	"path/filepath"
)

func GetAllTools() []Tool {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return []Tool{
		NewClaudeCode(filepath.Join(home, ".claude")),
		NewOpenCode(filepath.Join(home, ".config", "opencode")),
		NewCodex(filepath.Join(home, ".codex")),
		NewGeminiCLI(filepath.Join(home, ".gemini")),
		NewGrokCLI(filepath.Join(home, ".grok")),
		NewQwenCode(filepath.Join(home, ".qwen")),
		NewZed(filepath.Join(home, ".config", "zed")),
	}
}

func GetToolByName(name string) Tool {
	for _, tool := range GetAllTools() {
		if tool.Name() == name {
			return tool
		}
	}
	return nil
}
