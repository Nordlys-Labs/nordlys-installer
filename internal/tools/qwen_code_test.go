package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestQwenCode_Name(t *testing.T) {
	t.Parallel()

	q := NewQwenCode(t.TempDir())
	if got := q.Name(); got != "qwen-code" {
		t.Errorf("Name() = %q, want %q", got, "qwen-code")
	}
}

func TestQwenCode_Description(t *testing.T) {
	t.Parallel()

	q := NewQwenCode(t.TempDir())
	if got := q.Description(); got == "" {
		t.Error("Description() should not be empty")
	}
}

func TestQwenCode_RequiresNode(t *testing.T) {
	t.Parallel()

	q := NewQwenCode(t.TempDir())
	if !q.RequiresNode() {
		t.Error("RequiresNode() = false, want true")
	}
}

func TestQwenCode_IsInstalled(t *testing.T) {
	t.Parallel()

	q := NewQwenCode(t.TempDir())
	_ = q.IsInstalled()
}

func TestQwenCode_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	q := NewQwenCode(tmpDir)
	path, err := q.ConfigPath()
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

func TestQwenCode_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	q := NewQwenCode(tmpDir)

	err := q.UpdateConfig("test-api-key", "nordlys/hypernova", "https://api.test.com/v1")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	path, err := q.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	model, ok := data["model"].(map[string]any)
	if !ok {
		t.Fatal("model not found or not a map")
	}
	if model["name"] != "nordlys/hypernova" {
		t.Errorf("model.name = %v, want %q", model["name"], "nordlys/hypernova")
	}

	providers, ok := data["modelProviders"].(map[string]any)
	if !ok {
		t.Fatal("modelProviders not found or not a map")
	}
	openaiProviders, ok := providers["openai"].([]any)
	if !ok || len(openaiProviders) == 0 {
		t.Fatal("openai providers not found or empty")
	}
}

func TestQwenCode_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	q := NewQwenCode(tmpDir)

	// Valid key should pass
	err := q.UpdateConfig("valid-key-1234567890-long-enough", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = q.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	// Short key should fail
	err = q.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = q.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestQwenCode_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	q := NewQwenCode(tmpDir)
	configPath, err := q.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	initialData := map[string]any{
		"model": map[string]any{
			"name": "nordlys/hypernova",
		},
		"modelProviders": map[string]any{
			"openai": []any{
				map[string]any{
					"id":      "nordlys/hypernova",
					"baseUrl": "https://api.test.com/v1",
				},
			},
		},
		"security": map[string]any{
			"auth": map[string]any{
				"selectedType": "openai",
			},
		},
		"userSetting": "keep-this",
	}
	if err = config.WriteJSONFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	// Write env file
	envPath := filepath.Join(filepath.Dir(configPath), ".env")
	if err = os.WriteFile(envPath, []byte("OPENAI_API_KEY=test-key\n"), 0o600); err != nil {
		t.Fatalf("setup WriteFile() error = %v", err)
	}

	if err = q.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	if _, exists := data["model"]; exists {
		t.Error("Uninstall() should remove model field")
	}

	if _, exists := data["modelProviders"]; exists {
		t.Error("Uninstall() should remove modelProviders field")
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}
}

func TestQwenCode_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	q := NewQwenCode(tmpDir)

	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai/v1"

	// Install
	err := q.UpdateConfig(apiKey, model, baseURL)
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify
	configPath, err := q.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	modelConfig, ok := data["model"].(map[string]any)
	if !ok {
		t.Fatal("model not found or not a map")
	}
	if modelConfig["name"] != model {
		t.Errorf("model.name = %v, want %v", modelConfig["name"], model)
	}

	// Validate
	if err = q.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err = q.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
}
