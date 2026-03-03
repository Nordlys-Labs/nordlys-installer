package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type GeminiCLI struct {
	configDir string
}

func NewGeminiCLI(configDir string) *GeminiCLI {
	return &GeminiCLI{configDir: configDir}
}

func (g *GeminiCLI) Name() string        { return "gemini-cli" }
func (g *GeminiCLI) Description() string { return "Gemini CLI (Google)" }
func (g *GeminiCLI) RequiresNode() bool  { return true }

func (g *GeminiCLI) ConfigPath() (string, error) {
	return filepath.Join(g.configDir, "settings.json"), nil
}

func (g *GeminiCLI) IsInstalled() bool {
	_, err := exec.LookPath("gemini")
	return err == nil
}

func (g *GeminiCLI) UpdateConfig(apiKey, model, baseURL string) error {
	path, err := g.ConfigPath()
	if err != nil {
		return err
	}

	if model == "" {
		model = constants.DefaultModel
	}
	if baseURL == "" {
		baseURL = constants.APIBaseURL
	}

	// Build Gemini CLI config matching official structure from bash script
	geminiConfig := GeminiCLIConfig{
		Model: &GeminiModelConfig{
			Name: model,
		},
		Privacy: &GeminiPrivacyConfig{
			UsageStatisticsEnabled: false,
		},
		ModelConfigs: &GeminiModelConfigsConfig{
			CustomAliases: map[string]any{
				"summarizer-default": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
						"generateContentConfig": map[string]any{
							"maxOutputTokens": 2000,
							"temperature":     0.2,
						},
					},
				},
				"summarizer-shell": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
						"generateContentConfig": map[string]any{
							"maxOutputTokens": 2000,
							"temperature":     0,
						},
					},
				},
				"classifier": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
					},
				},
				"prompt-completion": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
					},
				},
				"edit-corrector": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
					},
				},
				"web-search": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
					},
				},
				"web-fetch": map[string]any{
					"modelConfig": map[string]any{
						"model": model,
					},
				},
			},
		},
	}

	// Validate against schema
	if err := ValidateConfig(GeminiCLISchemaURL, geminiConfig); err != nil {
		return err
	}

	// Convert to map for merging
	updates := structToMap(geminiConfig)

	// Write environment file for credentials using correct Gemini CLI env vars
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent := "GEMINI_API_KEY=" + apiKey + "\nGOOGLE_GEMINI_BASE_URL=" + baseURL + "\n"
	if err := os.WriteFile(envPath, []byte(envContent), 0o600); err != nil {
		return err
	}

	return config.UpdateJSONFields(path, updates)
}

func (g *GeminiCLI) GetExistingConfig() (apiKey, model string) {
	path, err := g.ConfigPath()
	if err != nil {
		return "", ""
	}
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		return "", ""
	}
	for _, line := range strings.Split(string(envContent), "\n") {
		if strings.HasPrefix(line, "GEMINI_API_KEY=") {
			apiKey = strings.TrimPrefix(line, "GEMINI_API_KEY=")
			break
		}
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		return apiKey, ""
	}
	if m, ok := data["model"].(map[string]any); ok {
		if name, ok := m["name"].(string); ok {
			model = name
		}
	}
	return apiKey, model
}

func (g *GeminiCLI) Validate() error {
	path, err := g.ConfigPath()
	if err != nil {
		return err
	}

	// Check if .env file exists with API key
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		return nil // .env file is optional
	}

	// Extract GEMINI_API_KEY from .env
	for _, line := range strings.Split(string(envContent), "\n") {
		if strings.HasPrefix(line, "GEMINI_API_KEY=") {
			apiKey := strings.TrimPrefix(line, "GEMINI_API_KEY=")
			return config.ValidateAPIKey(apiKey)
		}
	}

	return nil
}

func (g *GeminiCLI) Uninstall() error {
	path, err := g.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	delete(data, "model")
	delete(data, "privacy")
	delete(data, "modelConfigs")

	// Remove .env file
	envPath := filepath.Join(filepath.Dir(path), ".env")
	_ = os.Remove(envPath)

	return config.WriteJSONFile(path, data)
}
