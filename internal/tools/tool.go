package tools

type Tool interface {
	Name() string
	Description() string
	ConfigPath() (string, error)
	IsInstalled() bool
	RequiresNode() bool
	UpdateConfig(apiKey, model, baseURL string) error
	Validate() error
	Uninstall() error
	// GetExistingConfig returns apiKey and model from the tool's config if present.
	// Returns empty strings if not configured. Used by update command to avoid requiring flags.
	GetExistingConfig() (apiKey, model string)
}

type BaseConfig struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout int
}
