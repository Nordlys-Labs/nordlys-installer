package tools

import (
	"path/filepath"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestCodex_Name(t *testing.T) {
	t.Parallel()

	c := NewCodex(t.TempDir())
	if got := c.Name(); got != "codex" {
		t.Errorf("Name() = %q, want %q", got, "codex")
	}
}

func TestCodex_Description(t *testing.T) {
	t.Parallel()

	c := NewCodex(t.TempDir())
	if got := c.Description(); got == "" {
		t.Error("Description() should not be empty")
	}
}

func TestCodex_RequiresNode(t *testing.T) {
	t.Parallel()

	c := NewCodex(t.TempDir())
	if !c.RequiresNode() {
		t.Error("RequiresNode() = false, want true")
	}
}

func TestCodex_IsInstalled(t *testing.T) {
	t.Parallel()

	c := NewCodex(t.TempDir())
	_ = c.IsInstalled()
}

func TestCodex_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewCodex(tmpDir)
	path, err := c.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() should not be empty")
	}

	expected := filepath.Join(tmpDir, "config.toml")
	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func TestCodex_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewCodex(tmpDir)

	err := c.UpdateConfig("test-api-key", "nordlys/hypernova", "https://api.test.com/v1")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	path, err := c.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadTOMLFile(path)
	if err != nil {
		t.Fatalf("ReadTOMLFile() error = %v", err)
	}

	if data["model"] != "nordlys/hypernova" {
		t.Errorf("model = %v, want %q", data["model"], "nordlys/hypernova")
	}

	if data["model_provider"] != "nordlys" {
		t.Errorf("model_provider = %v, want %q", data["model_provider"], "nordlys")
	}

	providers, ok := data["model_providers"].(map[string]any)
	if !ok {
		t.Fatal("model_providers not found or not a map")
	}
	nordlys, ok := providers["nordlys"].(map[string]any)
	if !ok {
		t.Fatal("nordlys provider not found or not a map")
	}
	if nordlys["base_url"] != "https://api.test.com/v1" {
		t.Errorf("base_url = %v, want %q", nordlys["base_url"], "https://api.test.com/v1")
	}
}

func TestCodex_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewCodex(tmpDir)

	err := c.UpdateConfig("valid-key-1234567890", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = c.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	err = c.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = c.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestCodex_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewCodex(tmpDir)
	configPath, err := c.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	initialData := map[string]any{
		"model":          "nordlys/hypernova",
		"model_provider": "nordlys",
		"model_providers": map[string]any{
			"nordlys": map[string]any{
				"name":     "Nordlys",
				"base_url": "https://api.test.com",
			},
		},
		"userSetting": "keep-this",
	}
	if err = config.WriteTOMLFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteTOMLFile() error = %v", err)
	}

	if err = c.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadTOMLFile(configPath)
	if err != nil {
		t.Fatalf("ReadTOMLFile() after uninstall error = %v", err)
	}

	if _, exists := data["model"]; exists {
		t.Error("Uninstall() should remove model field")
	}

	if _, exists := data["model_provider"]; exists {
		t.Error("Uninstall() should remove model_provider field")
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}
}

func TestCodex_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	c := NewCodex(tmpDir)

	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai/v1"

	// Install
	err := c.UpdateConfig(apiKey, model, baseURL)
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify
	configPath, err := c.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadTOMLFile(configPath)
	if err != nil {
		t.Fatalf("ReadTOMLFile() error = %v", err)
	}

	if data["model"] != model {
		t.Errorf("model = %v, want %v", data["model"], model)
	}
	if data["model_provider"] != "nordlys" {
		t.Errorf("model_provider = %v, want %v", data["model_provider"], "nordlys")
	}

	// Validate
	if err = c.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err = c.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	// Verify uninstall
	data, err = config.ReadTOMLFile(configPath)
	if err != nil {
		t.Fatalf("ReadTOMLFile() after uninstall error = %v", err)
	}

	if _, exists := data["model"]; exists {
		t.Error("model not removed after uninstall")
	}
	if _, exists := data["model_provider"]; exists {
		t.Error("model_provider not removed after uninstall")
	}
}
