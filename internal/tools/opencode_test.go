package tools

import (
	"path/filepath"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestOpenCode_Name(t *testing.T) {
	t.Parallel()

	o := NewOpenCode(t.TempDir())
	if got := o.Name(); got != "opencode" {
		t.Errorf("Name() = %q, want %q", got, "opencode")
	}
}

func TestOpenCode_RequiresNode(t *testing.T) {
	t.Parallel()

	o := NewOpenCode(t.TempDir())
	if !o.RequiresNode() {
		t.Error("RequiresNode() = false, want true")
	}
}

func TestOpenCode_IsInstalled(t *testing.T) {
	t.Parallel()

	o := NewOpenCode(t.TempDir())
	_ = o.IsInstalled()
}

func TestOpenCode_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	o := NewOpenCode(tmpDir)
	path, err := o.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() should not be empty")
	}

	expected := filepath.Join(tmpDir, "opencode.json")
	if path != expected {
		t.Errorf("ConfigPath() = %q, want %q", path, expected)
	}
}

func TestOpenCode_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	o := NewOpenCode(tmpDir)

	err := o.UpdateConfig("test-api-key", "nordlys/hypernova", "https://api.test.com/v1")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	path, _ := o.ConfigPath()
	data, err := config.ReadJSONFile(path)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	if data["model"] != "nordlys/hypernova" {
		t.Errorf("model = %v, want %q", data["model"], "nordlys/hypernova")
	}

	provider, ok := data["provider"].(map[string]any)
	if !ok {
		t.Fatal("provider should be a map")
	}

	nordlys, ok := provider["nordlys"].(map[string]any)
	if !ok {
		t.Fatal("provider.nordlys should be a map")
	}

	options, ok := nordlys["options"].(map[string]any)
	if !ok {
		t.Fatal("provider.nordlys.options should be a map")
	}

	if options["apiKey"] != "test-api-key" {
		t.Errorf("apiKey = %v, want %q", options["apiKey"], "test-api-key")
	}
}

func TestOpenCode_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	o := NewOpenCode(tmpDir)

	err := o.UpdateConfig("valid-key-1234567890", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err := o.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	err = o.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err := o.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestOpenCode_ConfigPath_JSONC(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	o := NewOpenCode(tmpDir)

	jsoncPath := filepath.Join(tmpDir, "opencode.jsonc")
	if err := config.WriteJSONFile(jsoncPath, map[string]any{}); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	path, err := o.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path != jsoncPath {
		t.Errorf("ConfigPath() should prefer .jsonc, got %q, want %q", path, jsoncPath)
	}
}

func TestOpenCode_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	o := NewOpenCode(tmpDir)
	configPath, _ := o.ConfigPath()

	initialData := map[string]any{
		"provider": map[string]any{
			"nordlys": map[string]any{
				"options": map[string]any{
					"apiKey": "test-key",
				},
			},
			"other": map[string]any{
				"keep": "this",
			},
		},
		"model":       "nordlys/hypernova",
		"userSetting": "keep-this",
	}
	if err := config.WriteJSONFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	if err := o.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	provider, ok := data["provider"].(map[string]any)
	if !ok {
		t.Fatal("provider should still be a map")
	}

	if _, exists := provider["nordlys"]; exists {
		t.Error("Uninstall() should remove nordlys provider")
	}

	if provider["other"] == nil {
		t.Error("Uninstall() should preserve other providers")
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}
}

func TestOpenCode_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	o := NewOpenCode(tmpDir)
	configPath, _ := o.ConfigPath()

	// Create existing config with other providers
	existingConfig := map[string]any{
		"provider": map[string]any{
			"other": map[string]any{
				"keep": "this",
			},
		},
		"userPreference": "preserve",
	}
	if err := config.WriteJSONFile(configPath, existingConfig); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai/v1"

	// Install
	if err := o.UpdateConfig(apiKey, model, baseURL); err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify installation
	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	// User settings preserved
	if data["userPreference"] != "preserve" {
		t.Error("User preference not preserved")
	}

	provider, ok := data["provider"].(map[string]any)
	if !ok {
		t.Fatal("Provider block not found")
	}

	// Other provider preserved
	if _, exists := provider["other"]; !exists {
		t.Error("Other provider not preserved")
	}

	// Nordlys provider added
	nordlys, ok := provider["nordlys"].(map[string]any)
	if !ok {
		t.Fatal("Nordlys provider not found")
	}

	options, ok := nordlys["options"].(map[string]any)
	if !ok {
		t.Fatal("Options not found")
	}

	if options["apiKey"] != apiKey {
		t.Errorf("API key = %v, want %v", options["apiKey"], apiKey)
	}

	// Validate
	if err := o.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err := o.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	// Verify uninstall
	data, err = config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	provider, ok = data["provider"].(map[string]any)
	if !ok {
		t.Fatal("Provider block removed entirely")
	}

	// Other provider preserved
	if _, exists := provider["other"]; !exists {
		t.Error("Other provider not preserved after uninstall")
	}

	// Nordlys provider removed
	if _, exists := provider["nordlys"]; exists {
		t.Error("Nordlys provider not removed after uninstall")
	}
}
