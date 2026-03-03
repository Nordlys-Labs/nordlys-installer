package tools

import (
	"path/filepath"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestGrokCLI_Name(t *testing.T) {
	t.Parallel()

	g := NewGrokCLI(t.TempDir())
	if got := g.Name(); got != "grok-cli" {
		t.Errorf("Name() = %q, want %q", got, "grok-cli")
	}
}

func TestGrokCLI_Description(t *testing.T) {
	t.Parallel()

	g := NewGrokCLI(t.TempDir())
	if got := g.Description(); got == "" {
		t.Error("Description() should not be empty")
	}
}

func TestGrokCLI_RequiresNode(t *testing.T) {
	t.Parallel()

	g := NewGrokCLI(t.TempDir())
	if !g.RequiresNode() {
		t.Error("RequiresNode() = false, want true")
	}
}

func TestGrokCLI_IsInstalled(t *testing.T) {
	t.Parallel()

	g := NewGrokCLI(t.TempDir())
	_ = g.IsInstalled()
}

func TestGrokCLI_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGrokCLI(tmpDir)
	path, err := g.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() should not be empty")
	}

	expected := filepath.Join(tmpDir, "user-settings.json")
	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func TestGrokCLI_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGrokCLI(tmpDir)

	err := g.UpdateConfig("test-api-key", "nordlys/hypernova", "https://api.test.com")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	path, err := g.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	if data["defaultModel"] != "nordlys/hypernova" {
		t.Errorf("defaultModel = %v, want %q", data["defaultModel"], "nordlys/hypernova")
	}

	if data["apiKey"] != "test-api-key" {
		t.Errorf("apiKey = %v, want %q", data["apiKey"], "test-api-key")
	}

	models, ok := data["models"].([]any)
	if !ok || len(models) == 0 {
		t.Error("models should be a non-empty array")
	}
}

func TestGrokCLI_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGrokCLI(tmpDir)

	err := g.UpdateConfig("valid-key-1234567890", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = g.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	err = g.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = g.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestGrokCLI_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGrokCLI(tmpDir)
	configPath, err := g.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	initialData := map[string]any{
		"defaultModel": "nordlys/hypernova",
		"apiKey":       "test-key",
		"baseURL":      "https://api.test.com",
		"models":       []string{"nordlys/hypernova"},
		"userSetting":  "keep-this",
	}
	if err = config.WriteJSONFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	if err = g.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	if _, exists := data["defaultModel"]; exists {
		t.Error("Uninstall() should remove defaultModel field")
	}

	if _, exists := data["apiKey"]; exists {
		t.Error("Uninstall() should remove apiKey field")
	}

	if _, exists := data["models"]; exists {
		t.Error("Uninstall() should remove models field")
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}
}

func TestGrokCLI_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGrokCLI(tmpDir)

	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai"

	// Install
	err := g.UpdateConfig(apiKey, model, baseURL)
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify
	configPath, err := g.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	if data["defaultModel"] != model {
		t.Errorf("defaultModel = %v, want %v", data["defaultModel"], model)
	}
	if data["apiKey"] != apiKey {
		t.Errorf("apiKey = %v, want %v", data["apiKey"], apiKey)
	}

	// Validate
	if err = g.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err = g.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
}
