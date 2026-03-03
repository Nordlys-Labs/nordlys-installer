package tools

import (
	"os/exec"
	"path/filepath"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type OpenCode struct {
	configDir string
}

func NewOpenCode(configDir string) *OpenCode {
	return &OpenCode{configDir: configDir}
}

func (o *OpenCode) Name() string        { return "opencode" }
func (o *OpenCode) Description() string { return "OpenCode" }
func (o *OpenCode) RequiresNode() bool  { return true }

func (o *OpenCode) ConfigPath() (string, error) {
	jsoncPath := filepath.Join(o.configDir, "opencode.jsonc")
	if config.FileExists(jsoncPath) {
		return jsoncPath, nil
	}
	return filepath.Join(o.configDir, "opencode.json"), nil
}

func (o *OpenCode) IsInstalled() bool {
	_, err := exec.LookPath("opencode")
	return err == nil
}

func (o *OpenCode) UpdateConfig(apiKey, model, baseURL string) error {
	path, err := o.ConfigPath()
	if err != nil {
		return err
	}

	if model == "" {
		model = constants.DefaultModel
	}
	if baseURL == "" {
		baseURL = constants.APIBaseURL + "/v1"
	}

	// Build typed config
	nordlysConfig := OpenCodeConfig{
		Schema: OpenCodeSchemaURL,
		Provider: map[string]OpenCodeProvider{
			"nordlys": {
				NPM:  "@ai-sdk/openai-compatible",
				Name: "Nordlys",
				Options: OpenCodeProviderOptions{
					BaseURL: baseURL,
					APIKey:  apiKey,
					Headers: map[string]string{
						"User-Agent": "opencode-nordlys-integration",
					},
				},
				Models: map[string]OpenCodeModel{
					model: {
						Name: "Hypernova",
						Limit: OpenCodeLimit{
							Context: 200000,
							Output:  65536,
						},
					},
				},
			},
		},
		Model: model,
	}

	// Validate against schema
	if err := ValidateConfig(OpenCodeSchemaURL, nordlysConfig); err != nil {
		return err
	}

	// Convert to map for merging
	updates := structToMap(nordlysConfig)

	// Merge with existing config (update in place)
	return config.UpdateJSONFields(path, updates)
}

func (o *OpenCode) GetExistingConfig() (apiKey, model string) {
	path, err := o.ConfigPath()
	if err != nil {
		return "", ""
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		return "", ""
	}
	if provider, ok := data["provider"].(map[string]any); ok {
		if nordlys, ok := provider["nordlys"].(map[string]any); ok {
			if options, ok := nordlys["options"].(map[string]any); ok {
				if k, ok := options["apiKey"].(string); ok {
					apiKey = k
				}
			}
		}
	}
	if m, ok := data["model"].(string); ok {
		model = m
	}
	return apiKey, model
}

func (o *OpenCode) Validate() error {
	path, err := o.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	if provider, ok := data["provider"].(map[string]any); ok {
		if nordlys, ok := provider["nordlys"].(map[string]any); ok {
			if options, ok := nordlys["options"].(map[string]any); ok {
				if apiKey, ok := options["apiKey"].(string); ok {
					return config.ValidateAPIKey(apiKey)
				}
			}
		}
	}

	return nil
}

func (o *OpenCode) Uninstall() error {
	path, err := o.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	if provider, ok := data["provider"].(map[string]any); ok {
		delete(provider, "nordlys")
		data["provider"] = provider
	}

	if data["model"] == constants.DefaultModel {
		delete(data, "model")
	}

	return config.WriteJSONFile(path, data)
}
