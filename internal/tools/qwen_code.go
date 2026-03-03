package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type QwenCode struct {
	configDir string
}

func NewQwenCode(configDir string) *QwenCode {
	return &QwenCode{configDir: configDir}
}

func (q *QwenCode) Name() string        { return "qwen-code" }
func (q *QwenCode) Description() string { return "Qwen Code (Alibaba)" }
func (q *QwenCode) RequiresNode() bool  { return true }

func (q *QwenCode) ConfigPath() (string, error) {
	return filepath.Join(q.configDir, "settings.json"), nil
}

func (q *QwenCode) IsInstalled() bool {
	_, err := exec.LookPath("qwen")
	return err == nil
}

func (q *QwenCode) UpdateConfig(apiKey, model, baseURL string) error {
	path, err := q.ConfigPath()
	if err != nil {
		return err
	}

	if model == "" {
		model = constants.DefaultModel
	}
	if baseURL == "" {
		baseURL = constants.APIBaseURL + "/v1"
	}

	// Build typed config using Qwen Code's modelProviders structure
	qwenConfig := QwenCodeConfig{
		Model: QwenCodeModelConfig{
			Name: model,
		},
		ModelProviders: map[string][]QwenCodeProvider{
			"openai": {
				{
					ID:      model,
					Name:    "Nordlys Hypernova",
					EnvKey:  "OPENAI_API_KEY",
					BaseURL: baseURL,
				},
			},
		},
		Security: QwenCodeSecurityConfig{
			Auth: QwenCodeAuthConfig{
				SelectedType: "openai",
			},
		},
	}

	// Note: Qwen Code doesn't have a public schema, skip validation

	// Convert to map for merging
	updates := structToMap(qwenConfig)

	// Write environment file for credentials
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent := "OPENAI_API_KEY=" + apiKey + "\n"
	if err := os.WriteFile(envPath, []byte(envContent), 0o600); err != nil {
		return err
	}

	return config.UpdateJSONFields(path, updates)
}

func (q *QwenCode) GetExistingConfig() (apiKey, model string) {
	path, err := q.ConfigPath()
	if err != nil {
		return "", ""
	}
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		return "", ""
	}
	for _, line := range strings.Split(string(envContent), "\n") {
		if strings.HasPrefix(line, "OPENAI_API_KEY=") {
			apiKey = strings.TrimPrefix(line, "OPENAI_API_KEY=")
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

func (q *QwenCode) Validate() error {
	path, err := q.ConfigPath()
	if err != nil {
		return err
	}

	// Check if .env file exists with API key
	envPath := filepath.Join(filepath.Dir(path), ".env")
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		return nil // .env file is optional
	}

	// Extract OPENAI_API_KEY from .env
	for _, line := range strings.Split(string(envContent), "\n") {
		if strings.HasPrefix(line, "OPENAI_API_KEY=") {
			apiKey := strings.TrimPrefix(line, "OPENAI_API_KEY=")
			return config.ValidateAPIKey(apiKey)
		}
	}

	return nil
}

func (q *QwenCode) Uninstall() error {
	path, err := q.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	delete(data, "model")
	delete(data, "modelProviders")
	delete(data, "security")

	// Remove .env file
	envPath := filepath.Join(filepath.Dir(path), ".env")
	_ = os.Remove(envPath)

	return config.WriteJSONFile(path, data)
}
