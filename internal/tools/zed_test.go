package tools

import (
	"path/filepath"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestZed_Name(t *testing.T) {
	t.Parallel()

	z := NewZed(t.TempDir())
	if got := z.Name(); got != "zed" {
		t.Errorf("Name() = %q, want %q", got, "zed")
	}
}

func TestZed_Description(t *testing.T) {
	t.Parallel()

	z := NewZed(t.TempDir())
	if got := z.Description(); got == "" {
		t.Error("Description() should not be empty")
	}
}

func TestZed_RequiresNode(t *testing.T) {
	t.Parallel()

	z := NewZed(t.TempDir())
	if z.RequiresNode() {
		t.Error("RequiresNode() = true, want false for Zed")
	}
}

func TestZed_IsInstalled(t *testing.T) {
	t.Parallel()

	z := NewZed(t.TempDir())
	_ = z.IsInstalled()
}

func TestZed_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	z := NewZed(tmpDir)
	path, err := z.ConfigPath()
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

func TestZed_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	z := NewZed(tmpDir)
	configPath, err := z.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	existingData := map[string]any{
		"userSetting": "keep-this",
	}
	if err = config.WriteJSONFile(configPath, existingData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	err = z.UpdateConfig("test-api-key", "nordlys/hypernova", "https://api.test.com/v1")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	if data["userSetting"] != "keep-this" {
		t.Error("UpdateConfig() should preserve existing user settings")
	}

	langModels, ok := data["language_models"].(map[string]any)
	if !ok {
		t.Fatal("language_models should be a map")
	}

	openai, ok := langModels["openai"].(map[string]any)
	if !ok {
		t.Fatal("language_models.openai should be a map")
	}

	if openai["api_url"] != "https://api.test.com/v1" {
		t.Errorf("api_url = %v, want %q", openai["api_url"], "https://api.test.com/v1")
	}

	env, ok := data["env"].(map[string]any)
	if !ok {
		t.Fatal("env should be a map")
	}

	if env["OPENAI_API_KEY"] != "test-api-key" {
		t.Errorf("OPENAI_API_KEY = %v, want %q", env["OPENAI_API_KEY"], "test-api-key")
	}
}

func TestZed_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	z := NewZed(tmpDir)

	err := z.UpdateConfig("valid-key-1234567890", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = z.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	err = z.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = z.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestZed_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	z := NewZed(tmpDir)
	configPath, err := z.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	initialData := map[string]any{
		"language_models": map[string]any{
			"openai": map[string]any{
				"api_url": "https://api.test.com",
			},
		},
		"assistant": map[string]any{
			"default_model": map[string]any{
				"provider": "openai",
			},
		},
		"env": map[string]any{
			"OPENAI_API_KEY": "test-key",
			"OTHER_ENV_VAR":  "keep-this",
		},
		"userSetting": "keep-this",
	}
	if err = config.WriteJSONFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	if err = z.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	if _, exists := data["language_models"]; exists {
		t.Error("Uninstall() should remove language_models field")
	}

	if _, exists := data["assistant"]; exists {
		t.Error("Uninstall() should remove assistant field")
	}

	if env, ok := data["env"].(map[string]any); ok {
		if _, exists := env["OPENAI_API_KEY"]; exists {
			t.Error("Uninstall() should remove OPENAI_API_KEY from env")
		}
		if env["OTHER_ENV_VAR"] != "keep-this" {
			t.Error("Uninstall() should preserve other env vars")
		}
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}
}

func TestZed_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	z := NewZed(tmpDir)
	configPath, err := z.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	// Create existing config with user settings
	existingConfig := map[string]any{
		"userSetting": "preserve-this",
		"env": map[string]any{
			"OTHER_VAR": "keep-this",
		},
	}
	if err = config.WriteJSONFile(configPath, existingConfig); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai/v1"

	// Install
	if err = z.UpdateConfig(apiKey, model, baseURL); err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify installation
	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	// User settings preserved
	if data["userSetting"] != "preserve-this" {
		t.Error("User setting not preserved")
	}

	// Language models configured
	langModels, ok := data["language_models"].(map[string]any)
	if !ok {
		t.Fatal("Language models not found")
	}

	openai, ok := langModels["openai"].(map[string]any)
	if !ok {
		t.Fatal("OpenAI config not found")
	}

	if openai["api_url"] != baseURL {
		t.Errorf("API URL = %v, want %v", openai["api_url"], baseURL)
	}

	// Env configured
	env, ok := data["env"].(map[string]any)
	if !ok {
		t.Fatal("Env not found")
	}

	if env["OPENAI_API_KEY"] != apiKey {
		t.Errorf("API key = %v, want %v", env["OPENAI_API_KEY"], apiKey)
	}

	// Other env vars preserved
	if env["OTHER_VAR"] != "keep-this" {
		t.Error("Other env var not preserved")
	}

	// Validate
	if err = z.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err = z.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	// Verify uninstall
	data, err = config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	// User settings preserved
	if data["userSetting"] != "preserve-this" {
		t.Error("User setting not preserved after uninstall")
	}

	// Nordlys config removed
	if _, exists := data["language_models"]; exists {
		t.Error("language_models not removed after uninstall")
	}
	if _, exists := data["assistant"]; exists {
		t.Error("assistant not removed after uninstall")
	}

	// Env vars
	env, ok = data["env"].(map[string]any)
	if !ok {
		t.Fatal("Env removed entirely")
	}
	if env["OTHER_VAR"] != "keep-this" {
		t.Error("Other env var not preserved after uninstall")
	}
	if _, exists := env["OPENAI_API_KEY"]; exists {
		t.Error("API key not removed from env after uninstall")
	}
}
