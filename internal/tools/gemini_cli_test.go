package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
)

func TestGeminiCLI_Name(t *testing.T) {
	t.Parallel()

	g := NewGeminiCLI(t.TempDir())
	if got := g.Name(); got != "gemini-cli" {
		t.Errorf("Name() = %q, want %q", got, "gemini-cli")
	}
}

func TestGeminiCLI_Description(t *testing.T) {
	t.Parallel()

	g := NewGeminiCLI(t.TempDir())
	if got := g.Description(); got == "" {
		t.Error("Description() should not be empty")
	}
}

func TestGeminiCLI_RequiresNode(t *testing.T) {
	t.Parallel()

	g := NewGeminiCLI(t.TempDir())
	if !g.RequiresNode() {
		t.Error("RequiresNode() = false, want true")
	}
}

func TestGeminiCLI_IsInstalled(t *testing.T) {
	t.Parallel()

	g := NewGeminiCLI(t.TempDir())
	_ = g.IsInstalled()
}

func TestGeminiCLI_ConfigPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGeminiCLI(tmpDir)
	path, err := g.ConfigPath()
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

func TestGeminiCLI_UpdateConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGeminiCLI(tmpDir)

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

	// Check nested model structure
	modelObj, ok := data["model"].(map[string]any)
	if !ok {
		t.Fatalf("model field should be an object, got %T", data["model"])
	}
	if modelObj["name"] != "nordlys/hypernova" {
		t.Errorf("model.name = %v, want %q", modelObj["name"], "nordlys/hypernova")
	}

	// Check .env file for credentials (correct Gemini CLI env vars)
	envPath := filepath.Join(tmpDir, ".env")
	envBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("ReadFile(.env) error = %v", err)
	}
	if string(envBytes) != "GEMINI_API_KEY=test-api-key\nGOOGLE_GEMINI_BASE_URL=https://api.test.com\n" {
		t.Errorf(".env = %q, want GEMINI_API_KEY and GOOGLE_GEMINI_BASE_URL", string(envBytes))
	}
}

func TestGeminiCLI_Validate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGeminiCLI(tmpDir)

	// Valid key should pass
	err := g.UpdateConfig("valid-key-1234567890-long-enough", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = g.Validate(); err != nil {
		t.Errorf("Validate() should accept valid key, got error: %v", err)
	}

	// Short key should fail
	err = g.UpdateConfig("short", "", "")
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	if err = g.Validate(); err == nil {
		t.Error("Validate() should reject invalid key")
	}
}

func TestGeminiCLI_Uninstall(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGeminiCLI(tmpDir)
	configPath, err := g.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	initialData := map[string]any{
		"model": map[string]any{
			"name": "nordlys/hypernova",
		},
		"privacy": map[string]any{
			"usageStatisticsEnabled": false,
		},
		"modelConfigs": map[string]any{
			"customAliases": map[string]any{
				"summarizer-default": map[string]any{},
			},
		},
		"userSetting": "keep-this",
	}
	if err = config.WriteJSONFile(configPath, initialData); err != nil {
		t.Fatalf("setup WriteJSONFile() error = %v", err)
	}

	// Create .env file to test removal
	envPath := filepath.Join(tmpDir, ".env")
	if err = os.WriteFile(envPath, []byte("GEMINI_API_KEY=test\n"), 0o600); err != nil {
		t.Fatalf("setup WriteFile(.env) error = %v", err)
	}

	if err = g.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	if _, exists := data["model"]; exists {
		t.Error("Uninstall() should remove model field")
	}

	if _, exists := data["privacy"]; exists {
		t.Error("Uninstall() should remove privacy field")
	}

	if _, exists := data["modelConfigs"]; exists {
		t.Error("Uninstall() should remove modelConfigs field")
	}

	if data["userSetting"] != "keep-this" {
		t.Error("Uninstall() should preserve user settings")
	}

	// Check .env file was removed
	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		t.Error("Uninstall() should remove .env file")
	}
}

func TestGeminiCLI_FullWorkflow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	g := NewGeminiCLI(tmpDir)

	apiKey := "test-api-key-1234567890-valid"
	model := "nordlys/hypernova"
	baseURL := "https://api.nordlys.ai"

	// Install
	err := g.UpdateConfig(apiKey, model, baseURL)
	if err != nil {
		t.Fatalf("UpdateConfig() error = %v", err)
	}

	// Verify config structure
	configPath, err := g.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}
	data, err := config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() error = %v", err)
	}

	modelObj, ok := data["model"].(map[string]any)
	if !ok {
		t.Fatalf("model field should be an object")
	}
	if modelObj["name"] != model {
		t.Errorf("model.name = %v, want %v", modelObj["name"], model)
	}

	// Verify .env file
	envPath := filepath.Join(tmpDir, ".env")
	envBytes, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("ReadFile(.env) error = %v", err)
	}
	envContent := string(envBytes)
	if !strings.Contains(envContent, "GEMINI_API_KEY="+apiKey) {
		t.Errorf(".env missing GEMINI_API_KEY")
	}

	// Validate
	if err = g.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Uninstall
	if err = g.Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}

	// Verify uninstall
	data, err = config.ReadJSONFile(configPath)
	if err != nil {
		t.Fatalf("ReadJSONFile() after uninstall error = %v", err)
	}

	if _, exists := data["model"]; exists {
		t.Error("model not removed after uninstall")
	}
	if _, exists := data["privacy"]; exists {
		t.Error("privacy not removed after uninstall")
	}
	if _, exists := data["modelConfigs"]; exists {
		t.Error("modelConfigs not removed after uninstall")
	}
}
