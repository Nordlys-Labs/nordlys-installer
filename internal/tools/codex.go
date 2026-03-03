package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type Codex struct {
	configDir string
}

func NewCodex(configDir string) *Codex {
	return &Codex{configDir: configDir}
}

func (c *Codex) Name() string        { return "codex" }
func (c *Codex) Description() string { return "OpenAI Codex CLI" }
func (c *Codex) RequiresNode() bool  { return true }

func (c *Codex) ConfigPath() (string, error) {
	return filepath.Join(c.configDir, "config.toml"), nil
}

func (c *Codex) IsInstalled() bool {
	_, err := exec.LookPath("codex")
	return err == nil
}

func (c *Codex) UpdateConfig(apiKey, model, baseURL string) error {
	path, err := c.ConfigPath()
	if err != nil {
		return err
	}

	if model == "" {
		model = constants.DefaultModel
	}
	if baseURL == "" {
		baseURL = constants.APIBaseURL + "/v1"
	}

	// Build typed config using Codex's model_providers structure
	codexConfig := CodexConfig{
		Model:         model,
		ModelProvider: "nordlys",
		ModelProviders: map[string]CodexModelProvider{
			"nordlys": {
				Name:               "Nordlys",
				BaseURL:            baseURL,
				EnvKey:             "NORDLYS_API_KEY",
				EnvKeyInstructions: "Get your API key from https://nordlyslabs.com/api-platform/orgs",
				WireAPI:            "responses",
				RequiresOpenAIAuth: false,
			},
		},
	}

	// Validate against schema (convert to JSON for validation)
	if err := ValidateConfig(CodexSchemaURL, codexConfig); err != nil {
		return err
	}

	// Write environment file for credentials
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent := "NORDLYS_API_KEY=" + apiKey + "\n"
	if err := os.WriteFile(envPath, []byte(envContent), 0o600); err != nil {
		return err
	}

	// Marshal the config directly to TOML format
	updates := map[string]any{
		"model":          codexConfig.Model,
		"model_provider": codexConfig.ModelProvider,
		"model_providers": map[string]any{
			"nordlys": map[string]any{
				"name":                 codexConfig.ModelProviders["nordlys"].Name,
				"base_url":             codexConfig.ModelProviders["nordlys"].BaseURL,
				"env_key":              codexConfig.ModelProviders["nordlys"].EnvKey,
				"env_key_instructions": codexConfig.ModelProviders["nordlys"].EnvKeyInstructions,
				"wire_api":             codexConfig.ModelProviders["nordlys"].WireAPI,
				"requires_openai_auth": codexConfig.ModelProviders["nordlys"].RequiresOpenAIAuth,
			},
		},
	}

	return config.UpdateTOMLFields(path, updates)
}

func (c *Codex) GetExistingConfig() (apiKey, model string) {
	path, err := c.ConfigPath()
	if err != nil {
		return "", ""
	}
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envData, err := os.ReadFile(envPath)
	if err != nil {
		return "", ""
	}
	for _, line := range strings.Split(string(envData), "\n") {
		if strings.HasPrefix(line, "NORDLYS_API_KEY=") {
			apiKey = strings.TrimPrefix(line, "NORDLYS_API_KEY=")
			break
		}
	}
	data, err := config.ReadTOMLFile(path)
	if err != nil {
		return apiKey, ""
	}
	if m, ok := data["model"].(string); ok {
		model = m
	}
	return apiKey, model
}

func (c *Codex) Validate() error {
	path, err := c.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadTOMLFile(path)
	if err != nil {
		return err
	}

	// Check if nordlys provider is configured
	if providers, ok := data["model_providers"].(map[string]any); ok {
		if _, ok := providers["nordlys"]; ok {
			// Read API key from env file
			envPath := filepath.Join(filepath.Dir(path), ".env")
			envData, err := os.ReadFile(envPath)
			if err == nil {
				lines := strings.Split(string(envData), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "NORDLYS_API_KEY=") {
						apiKey := strings.TrimPrefix(line, "NORDLYS_API_KEY=")
						return config.ValidateAPIKey(apiKey)
					}
				}
			}
		}
	}

	return nil
}

func (c *Codex) Uninstall() error {
	path, err := c.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadTOMLFile(path)
	if err != nil {
		return err
	}

	// Remove nordlys provider
	if providers, ok := data["model_providers"].(map[string]any); ok {
		delete(providers, "nordlys")
		data["model_providers"] = providers
	}

	// Reset model provider if it was nordlys
	if data["model_provider"] == "nordlys" {
		delete(data, "model_provider")
		delete(data, "model")
	}

	// Remove env file (ignore if not exists)
	envPath := filepath.Join(filepath.Dir(path), ".env")
	if err := os.Remove(envPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove env file: %w", err)
	}

	return config.WriteTOMLFile(path, data)
}
