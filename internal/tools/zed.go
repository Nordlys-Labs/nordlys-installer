package tools

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type Zed struct {
	configDir string
}

func NewZed(configDir string) *Zed {
	return &Zed{configDir: configDir}
}

func (z *Zed) Name() string        { return "zed" }
func (z *Zed) Description() string { return "Zed Editor" }
func (z *Zed) RequiresNode() bool  { return false }

func (z *Zed) ConfigPath() (string, error) {
	return filepath.Join(z.configDir, "settings.json"), nil
}

func (z *Zed) IsInstalled() bool {
	_, err := exec.LookPath("zed")
	return err == nil
}

func (z *Zed) UpdateConfig(apiKey, model, baseURL string) error {
	path, err := z.ConfigPath()
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
	zedConfig := ZedConfig{
		LanguageModels: ZedLanguageModels{
			OpenAI: ZedOpenAIConfig{
				Version: "1",
				APIURL:  baseURL,
				AvailableModels: []ZedModel{
					{
						Name:            model,
						DisplayName:     "Nordlys Hypernova",
						MaxTokens:       200000,
						MaxOutputTokens: 65536,
					},
				},
			},
		},
		Assistant: ZedAssistant{
			DefaultModel: ZedDefaultModel{
				Provider: "openai",
				Model:    model,
			},
		},
		Env: map[string]string{
			"OPENAI_API_KEY": apiKey,
		},
	}

	// Note: Zed doesn't have a public schema, skip validation

	existing, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		if err := config.CreateBackup(path); err != nil {
			return err
		}
	}

	// Convert to map for merging
	updates := structToMap(zedConfig)

	// Deep merge with existing config
	merged := config.DeepMerge(existing, updates)

	return config.WriteJSONFile(path, merged)
}

func (z *Zed) GetExistingConfig() (apiKey, model string) {
	path, err := z.ConfigPath()
	if err != nil {
		return "", ""
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		return "", ""
	}
	if env, ok := data["env"].(map[string]any); ok {
		if k, ok := env["OPENAI_API_KEY"].(string); ok {
			apiKey = k
		}
	}
	if assistant, ok := data["assistant"].(map[string]any); ok {
		if dm, ok := assistant["default_model"].(map[string]any); ok {
			if m, ok := dm["model"].(string); ok {
				model = m
			}
		}
	}
	return apiKey, model
}

func (z *Zed) Validate() error {
	path, err := z.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	if env, ok := data["env"].(map[string]any); ok {
		if apiKey, ok := env["OPENAI_API_KEY"].(string); ok {
			return config.ValidateAPIKey(apiKey)
		}
	}

	return nil
}

func (z *Zed) Uninstall() error {
	path, err := z.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	delete(data, "language_models")
	delete(data, "assistant")

	if env, ok := data["env"].(map[string]any); ok {
		delete(env, "OPENAI_API_KEY")
		data["env"] = env
	}

	return config.WriteJSONFile(path, data)
}
