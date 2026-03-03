package tools

import (
	"path/filepath"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestClaudeCode_Name(t *testing.T) {
	t.Parallel()

	c := NewClaudeCode(t.TempDir())
	if got := c.Name(); got != "claude-code" {
		t.Errorf("Name() = %q, want %q", got, "claude-code")
	}
}

func TestClaudeCode_Description(t *testing.T) {
	t.Parallel()

	c := NewClaudeCode(t.TempDir())
	if got := c.Description(); got == "" {
		t.Error("Description() should not be empty")
	}
}

func TestClaudeCode_RequiresNode(t *testing.T) {
	t.Parallel()

	c := NewClaudeCode(t.TempDir())
	if !c.RequiresNode() {
		t.Error("RequiresNode() = false, want true")
	}
}

func TestClaudeCode_IsInstalled(t *testing.T) {
	t.Parallel()

	c := NewClaudeCode(t.TempDir())
	_ = c.IsInstalled()
}

func TestClaudeCode_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewClaudeCode(tmpDir)
	path, err := c.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() should not be empty")
	}

	expected := filepath.Join(tmpDir, "settings.json")
	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func TestClaudeCode_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewClaudeCode(tmpDir)

	err := c.UpdateConfig("test-api-key", "nordlys/hypernova", "https://api.test.com")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	path, _ := c.ConfigPath()
	data, err := config.ReadJSONFile(path)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	if data["model"] != "nordlys/hypernova" {
		t.Errorf("model = %v, want %q", data["model"], "nordlys/hypernova")
	}

	env, ok := data["env"].(map[string]any)
	if !ok {
		t.Fatal("env should be a map")
	}

	if env["ANTHROPIC_AUTH_TOKEN"] != "test-api-key" {
		t.Errorf("ANTHROPIC_AUTH_TOKEN = %v, want %q", env["ANTHROPIC_AUTH_TOKEN"], "test-api-key")
	}

	if env["ANTHROPIC_BASE_URL"] != "https://api.test.com" {
		t.Errorf("ANTHROPIC_BASE_URL = %v, want %q", env["ANTHROPIC_BASE_URL"], "https://api.test.com")
	}
}

func TestClaudeCode_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewClaudeCode(tmpDir)

	err := c.UpdateConfig("valid-key-1234567890", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err := c.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	err = c.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err := c.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestClaudeCode_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewClaudeCode(tmpDir)
	configPath, _ := c.ConfigPath()

	initialData := map[string]any{
		"model": "nordlys/hypernova",
		"env": map[string]any{
			"ANTHROPIC_AUTH_TOKEN": "test-key",
		},
		"userSetting": "keep-this",
	}
	if err := config.WriteJSONFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	if err := c.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	if _, exists := data["model"]; exists {
		t.Error("Uninstall() should remove model field")
	}

	if _, exists := data["env"]; exists {
		t.Error("Uninstall() should remove env field")
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}
}

func TestClaudeCode_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewClaudeCode(tmpDir)
	configPath, _ := c.ConfigPath()

	// Create existing config with user settings
	existingConfig := map[string]any{
		"userSetting": "preserve-this",
		"theme":       "dark",
		"fontSize":    14,
	}
	if err := config.WriteJSONFile(configPath, existingConfig); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	// Test full install -> validate -> uninstall workflow
	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai"

	// Install
	if err := c.UpdateConfig(apiKey, model, baseURL); err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify installation
	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	// User settings preserved
	if data["userSetting"] != "preserve-this" {
		t.Error("User setting not preserved after install")
	}
	if data["theme"] != "dark" {
		t.Error("Theme not preserved after install")
	}

	// Nordlys config added
	if data["model"] != model {
		t.Errorf("model = %v, want %v", data["model"], model)
	}

	env, ok := data["env"].(map[string]any)
	if !ok {
		t.Fatal("env block not found after install")
	}
	if env["ANTHROPIC_AUTH_TOKEN"] != apiKey {
		t.Errorf("API key = %v, want %v", env["ANTHROPIC_AUTH_TOKEN"], apiKey)
	}
	if env["ANTHROPIC_BASE_URL"] != baseURL {
		t.Errorf("base URL = %v, want %v", env["ANTHROPIC_BASE_URL"], baseURL)
	}

	// Validate
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err := c.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	// Verify uninstall
	data, err = config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	// User settings still preserved
	if data["userSetting"] != "preserve-this" {
		t.Error("User setting not preserved after uninstall")
	}

	// Nordlys config removed
	if _, exists := data["model"]; exists {
		t.Error("model not removed after uninstall")
	}
	if _, exists := data["env"]; exists {
		t.Error("env not removed after uninstall")
	}
}
