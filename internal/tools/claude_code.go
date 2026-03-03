package tools

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
)

type ClaudeCode struct {
	configDir string
}

func NewClaudeCode(configDir string) *ClaudeCode {
	return &ClaudeCode{configDir: configDir}
}

func (c *ClaudeCode) Name() string        { return "claude-code" }
func (c *ClaudeCode) Description() string { return "Claude Code (Anthropic)" }
func (c *ClaudeCode) RequiresNode() bool  { return true }

func (c *ClaudeCode) ConfigPath() (string, error) {
	return filepath.Join(c.configDir, "settings.json"), nil
}

func (c *ClaudeCode) IsInstalled() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

func (c *ClaudeCode) UpdateConfig(apiKey, model, baseURL string) error {
	path, err := c.ConfigPath()
	if err != nil {
		return err
	}

	if model == "" {
		model = constants.DefaultModel
	}
	if baseURL == "" {
		baseURL = constants.APIBaseURL
	}

	// Build typed config with comprehensive model overrides
	// This ensures ALL Claude Code operations route through Nordlys
	claudeConfig := ClaudeCodeConfig{
		Model: model,
		Env: ClaudeCodeEnvironment{
			// Authentication
			AnthropicAuthToken: apiKey,
			AnthropicBaseURL:   baseURL,
			APITimeoutMS:       "3000000",

			// Model configuration - override ALL model aliases to route through Nordlys
			AnthropicModel:          model, // Default model
			ClaudeCodeSubagentModel: model, // Subagents (e.g., Task tool spawned agents)
			AnthropicDefaultHaiku:   model, // Background tasks, fast operations
			AnthropicDefaultSonnet:  model, // Sonnet alias
			AnthropicDefaultOpus:    model, // Opus alias
		},
	}

	// Note: Claude Code doesn't have a public schema, skip validation

	// Convert to map for merging
	updates := structToMap(claudeConfig)

	// Update settings.json
	if err := config.UpdateJSONFields(path, updates); err != nil {
		return err
	}

	// Also update .claude.json with hasCompletedOnboarding (matches bash script)
	claudeJSONPath := filepath.Join(filepath.Dir(c.configDir), ".claude.json")
	claudeJSON, err := config.ReadJSONFile(claudeJSONPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if claudeJSON == nil {
		claudeJSON = make(map[string]any)
	}
	claudeJSON["hasCompletedOnboarding"] = true
	return config.WriteJSONFile(claudeJSONPath, claudeJSON)
}

func (c *ClaudeCode) GetExistingConfig() (apiKey, model string) {
	path, err := c.ConfigPath()
	if err != nil {
		return "", ""
	}
	data, err := config.ReadJSONFile(path)
	if err != nil {
		return "", ""
	}
	if env, ok := data["env"].(map[string]any); ok {
		if token, ok := env["ANTHROPIC_AUTH_TOKEN"].(string); ok {
			apiKey = token
		}
	}
	if m, ok := data["model"].(string); ok {
		model = m
	}
	return apiKey, model
}

func (c *ClaudeCode) Validate() error {
	path, err := c.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	if env, ok := data["env"].(map[string]any); ok {
		if token, ok := env["ANTHROPIC_AUTH_TOKEN"].(string); ok {
			return config.ValidateAPIKey(token)
		}
	}

	return nil
}

func (c *ClaudeCode) Uninstall() error {
	path, err := c.ConfigPath()
	if err != nil {
		return err
	}

	data, err := config.ReadJSONFile(path)
	if err != nil {
		return err
	}

	delete(data, "model")
	delete(data, "env")

	return config.WriteJSONFile(path, data)
}
