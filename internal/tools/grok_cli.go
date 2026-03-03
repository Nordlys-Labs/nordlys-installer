package tools

import (
	"os/exec"
	"path/filepath"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type GrokCLI struct {
	configDir string
}

func NewGrokCLI(configDir string) *GrokCLI {
	return &GrokCLI{configDir: configDir}
}

func (g *GrokCLI) Name() string        { return "grok-cli" }
func (g *GrokCLI) Description() string { return "Grok CLI (xAI)" }
func (g *GrokCLI) RequiresNode() bool  { return true }

func (g *GrokCLI) ConfigPath() (string, error) {
	return filepath.Join(g.configDir, "user-settings.json"), nil
}

func (g *GrokCLI) IsInstalled() bool {
	_, err := exec.LookPath("grok")
	return err == nil
}

func (g *GrokCLI) UpdateConfig(apiKey, model, baseURL string) error {
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

	// Build typed config
	grokConfig := GrokCLIConfig{
		APIKey:       apiKey,
		BaseURL:      baseURL,
		DefaultModel: model,
		Models:       []string{model},
	}

	// Note: Grok CLI doesn't have a public schema, skip validation

	// Convert to map for merging
	updates := structToMap(grokConfig)

	return config.UpdateJSONFields(path, updates)
}

func (g *GrokCLI) GetExistingConfig() (apiKey, model string) {
	path, err := g.ConfigPath()
	if err != nil {
		return "", ""
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		return "", ""
	}
	if k, ok := data["apiKey"].(string); ok {
		apiKey = k
	}
	if m, ok := data["defaultModel"].(string); ok {
		model = m
	}
	return apiKey, model
}

func (g *GrokCLI) Validate() error {
	path, err := g.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	grokConfig := GrokCLIConfig{}
	if apiKey, ok := data["apiKey"].(string); ok {
		grokConfig.APIKey = apiKey
		return config.ValidateAPIKey(grokConfig.APIKey)
	}

	return nil
}

func (g *GrokCLI) Uninstall() error {
	path, err := g.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	delete(data, "apiKey")
	delete(data, "baseURL")
	delete(data, "defaultModel")
	delete(data, "models")

	return config.WriteJSONFile(path, data)
}
