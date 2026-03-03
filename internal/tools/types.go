package tools

// This file contains all exported configuration types for AI coding tools

// =============================================================================
// OpenCode Configuration Types
// =============================================================================

// OpenCodeConfig represents the OpenCode JSON configuration structure
type OpenCodeConfig struct {
	Schema   string                      `json:"$schema,omitzero"`
	Provider map[string]OpenCodeProvider `json:"provider,omitzero"`
	Model    string                      `json:"model,omitzero"`
}

// OpenCodeProvider represents a provider configuration in OpenCode
type OpenCodeProvider struct {
	Options OpenCodeProviderOptions  `json:"options,omitzero"`
	Models  map[string]OpenCodeModel `json:"models,omitzero"`
	NPM     string                   `json:"npm,omitzero"`
	Name    string                   `json:"name,omitzero"`
}

// OpenCodeProviderOptions contains provider-specific options
type OpenCodeProviderOptions struct {
	Headers map[string]string `json:"headers,omitzero"`
	BaseURL string            `json:"baseURL,omitzero"`
	APIKey  string            `json:"apiKey,omitzero"`
}

// OpenCodeModel represents a model configuration
type OpenCodeModel struct {
	Name  string        `json:"name,omitzero"`
	Limit OpenCodeLimit `json:"limit,omitzero"`
}

// OpenCodeLimit defines token limits for a model
type OpenCodeLimit struct {
	Context int `json:"context,omitzero"`
	Output  int `json:"output,omitzero"`
}

// =============================================================================
// Claude Code Configuration Types
// =============================================================================

// ClaudeCodeConfig represents Claude Code JSON configuration
type ClaudeCodeConfig struct {
	Model string                `json:"model,omitzero"`
	Env   ClaudeCodeEnvironment `json:"env,omitzero"`
}

// ClaudeCodeEnvironment contains environment variables for Claude Code
// See: https://docs.anthropic.com/en/docs/claude-code/settings#environment-variables
type ClaudeCodeEnvironment struct {
	// Authentication
	AnthropicAuthToken string `json:"ANTHROPIC_AUTH_TOKEN,omitzero"`
	AnthropicBaseURL   string `json:"ANTHROPIC_BASE_URL,omitzero"`

	// Timeouts
	APITimeoutMS string `json:"API_TIMEOUT_MS,omitzero"`

	// Model configuration - route all model aliases through Nordlys
	AnthropicModel          string `json:"ANTHROPIC_MODEL,omitzero"`
	ClaudeCodeSubagentModel string `json:"CLAUDE_CODE_SUBAGENT_MODEL,omitzero"`
	AnthropicDefaultHaiku   string `json:"ANTHROPIC_DEFAULT_HAIKU_MODEL,omitzero"`
	AnthropicDefaultSonnet  string `json:"ANTHROPIC_DEFAULT_SONNET_MODEL,omitzero"`
	AnthropicDefaultOpus    string `json:"ANTHROPIC_DEFAULT_OPUS_MODEL,omitzero"`
}

// =============================================================================
// Zed Configuration Types
// =============================================================================

// ZedConfig represents Zed editor JSON configuration
type ZedConfig struct {
	Env            map[string]string `json:"env,omitzero"`
	Assistant      ZedAssistant      `json:"assistant,omitzero"`
	LanguageModels ZedLanguageModels `json:"language_models,omitzero"`
}

// ZedLanguageModels contains language model configurations
type ZedLanguageModels struct {
	OpenAI ZedOpenAIConfig `json:"openai,omitzero"`
}

// ZedOpenAIConfig represents OpenAI-compatible configuration in Zed
type ZedOpenAIConfig struct {
	Version         string     `json:"version,omitzero"`
	APIURL          string     `json:"api_url,omitzero"`
	AvailableModels []ZedModel `json:"available_models,omitzero"`
}

// ZedModel represents a model configuration in Zed
type ZedModel struct {
	Name            string `json:"name,omitzero"`
	DisplayName     string `json:"display_name,omitzero"`
	MaxTokens       int    `json:"max_tokens,omitzero"`
	MaxOutputTokens int    `json:"max_output_tokens,omitzero"`
}

// ZedAssistant contains assistant configuration
type ZedAssistant struct {
	DefaultModel ZedDefaultModel `json:"default_model,omitzero"`
}

// ZedDefaultModel specifies the default model
type ZedDefaultModel struct {
	Provider string `json:"provider,omitzero"`
	Model    string `json:"model,omitzero"`
}

// =============================================================================
// Codex Configuration Types (TOML)
// =============================================================================

// CodexConfig represents Codex TOML configuration
type CodexConfig struct {
	ModelProviders map[string]CodexModelProvider `toml:"model_providers,omitzero"`
	Model          string                        `toml:"model,omitzero"`
	ModelProvider  string                        `toml:"model_provider,omitzero"`
}

// CodexModelProvider represents a model provider in Codex
type CodexModelProvider struct {
	HTTPHeaders         map[string]string `toml:"http_headers,omitzero"`
	EnvHTTPHeaders      map[string]string `toml:"env_http_headers,omitzero"`
	QueryParams         map[string]string `toml:"query_params,omitzero"`
	Name                string            `toml:"name,omitzero"`
	BaseURL             string            `toml:"base_url,omitzero"`
	EnvKey              string            `toml:"env_key,omitzero"`
	EnvKeyInstructions  string            `toml:"env_key_instructions,omitzero"`
	WireAPI             string            `toml:"wire_api,omitzero"`
	RequestMaxRetries   int               `toml:"request_max_retries,omitzero"`
	StreamIdleTimeoutMS int               `toml:"stream_idle_timeout_ms,omitzero"`
	StreamMaxRetries    int               `toml:"stream_max_retries,omitzero"`
	RequiresOpenAIAuth  bool              `toml:"requires_openai_auth,omitzero"`
}

// =============================================================================
// Grok CLI Configuration Types
// =============================================================================

// GrokCLIConfig represents Grok CLI JSON configuration
type GrokCLIConfig struct {
	APIKey       string   `json:"apiKey,omitzero"`
	BaseURL      string   `json:"baseURL,omitzero"`
	DefaultModel string   `json:"defaultModel,omitzero"`
	Models       []string `json:"models,omitzero"`
}

// =============================================================================
// Qwen Code Configuration Types
// =============================================================================

// QwenCodeConfig represents Qwen Code JSON configuration
type QwenCodeConfig struct {
	Model          QwenCodeModelConfig           `json:"model,omitzero"`
	ModelProviders map[string][]QwenCodeProvider `json:"modelProviders,omitzero"`
	Security       QwenCodeSecurityConfig        `json:"security,omitzero"`
}

// QwenCodeModelConfig contains model-specific settings
type QwenCodeModelConfig struct {
	Name string `json:"name,omitzero"`
}

// QwenCodeProvider represents a model provider in Qwen Code
type QwenCodeProvider struct {
	GenerationConfig *QwenGenerationConfig `json:"generationConfig,omitzero"`
	ID               string                `json:"id"`
	Name             string                `json:"name,omitzero"`
	EnvKey           string                `json:"envKey"`
	BaseURL          string                `json:"baseUrl,omitzero"`
}

// QwenGenerationConfig contains generation parameters
type QwenGenerationConfig struct {
	CustomHeaders  map[string]string   `json:"customHeaders,omitzero"`
	ExtraBody      map[string]any      `json:"extra_body,omitzero"`
	SamplingParams *QwenSamplingParams `json:"samplingParams,omitzero"`
	Timeout        int                 `json:"timeout,omitzero"`
	MaxRetries     int                 `json:"maxRetries,omitzero"`
}

// QwenSamplingParams contains sampling parameters
type QwenSamplingParams struct {
	Temperature float64 `json:"temperature,omitzero"`
	TopP        float64 `json:"top_p,omitzero"`
	MaxTokens   int     `json:"max_tokens,omitzero"`
}

// QwenCodeSecurityConfig contains security settings
type QwenCodeSecurityConfig struct {
	Auth QwenCodeAuthConfig `json:"auth,omitzero"`
}

// QwenCodeAuthConfig contains authentication settings
type QwenCodeAuthConfig struct {
	SelectedType string `json:"selectedType,omitzero"`
}

// This file contains the complete Gemini CLI configuration types

// GeminiCLIConfig represents the complete Gemini CLI settings.json structure
type GeminiCLIConfig struct {
	Tools         *GeminiToolsConfig         `json:"tools,omitzero"`
	MCP           *GeminiMCPConfig           `json:"mcp,omitzero"`
	Output        *GeminiOutputConfig        `json:"output,omitzero"`
	UI            *GeminiUIConfig            `json:"ui,omitzero"`
	IDE           *GeminiIDEConfig           `json:"ide,omitzero"`
	Privacy       *GeminiPrivacyConfig       `json:"privacy,omitzero"`
	Model         *GeminiModelConfig         `json:"model,omitzero"`
	ModelConfigs  *GeminiModelConfigsConfig  `json:"modelConfigs,omitzero"`
	Agents        *GeminiAgentsConfig        `json:"agents,omitzero"`
	Telemetry     *GeminiTelemetryConfig     `json:"telemetry,omitzero"`
	General       *GeminiGeneralConfig       `json:"general,omitzero"`
	MCPServers    map[string]GeminiMCPServer `json:"mcpServers,omitzero"`
	Context       *GeminiContextConfig       `json:"context,omitzero"`
	Security      *GeminiSecurityConfig      `json:"security,omitzero"`
	Advanced      *GeminiAdvancedConfig      `json:"advanced,omitzero"`
	Experimental  *GeminiExperimentalConfig  `json:"experimental,omitzero"`
	Skills        *GeminiSkillsConfig        `json:"skills,omitzero"`
	HooksConfig   *GeminiHooksConfigConfig   `json:"hooksConfig,omitzero"`
	Hooks         *GeminiHooksConfig         `json:"hooks,omitzero"`
	Admin         *GeminiAdminConfig         `json:"admin,omitzero"`
	PolicyPaths   []string                   `json:"policyPaths,omitzero"`
	UseWriteTodos bool                       `json:"useWriteTodos,omitzero"`
}

type GeminiGeneralConfig struct {
	Checkpointing                *GeminiCheckpointingConfig    `json:"checkpointing,omitzero"`
	SessionRetention             *GeminiSessionRetentionConfig `json:"sessionRetention,omitzero"`
	PreferredEditor              string                        `json:"preferredEditor,omitzero"`
	DefaultApprovalMode          string                        `json:"defaultApprovalMode,omitzero"`
	VimMode                      bool                          `json:"vimMode,omitzero"`
	Devtools                     bool                          `json:"devtools,omitzero"`
	EnableAutoUpdate             bool                          `json:"enableAutoUpdate,omitzero"`
	EnableAutoUpdateNotification bool                          `json:"enableAutoUpdateNotification,omitzero"`
	EnablePromptCompletion       bool                          `json:"enablePromptCompletion,omitzero"`
	RetryFetchErrors             bool                          `json:"retryFetchErrors,omitzero"`
	DebugKeystrokeLogging        bool                          `json:"debugKeystrokeLogging,omitzero"`
}

type GeminiCheckpointingConfig struct {
	Enabled bool `json:"enabled,omitzero"`
}

type GeminiSessionRetentionConfig struct {
	MaxAge              string `json:"maxAge,omitzero"`
	MinRetention        string `json:"minRetention,omitzero"`
	MaxCount            int    `json:"maxCount,omitzero"`
	Enabled             bool   `json:"enabled,omitzero"`
	WarningAcknowledged bool   `json:"warningAcknowledged,omitzero"`
}

type GeminiOutputConfig struct {
	Format string `json:"format,omitzero"`
}

type GeminiUIConfig struct {
	Accessibility                     *GeminiAccessibilityConfig `json:"accessibility,omitzero"`
	CustomThemes                      map[string]any             `json:"customThemes,omitzero"`
	Footer                            *GeminiUIFooterConfig      `json:"footer,omitzero"`
	Theme                             string                     `json:"theme,omitzero"`
	InlineThinkingMode                string                     `json:"inlineThinkingMode,omitzero"`
	CustomWittyPhrases                []string                   `json:"customWittyPhrases,omitzero"`
	TerminalBackgroundPollingInterval int                        `json:"terminalBackgroundPollingInterval,omitzero"`
	HideContextSummary                bool                       `json:"hideContextSummary,omitzero"`
	ShowLineNumbers                   bool                       `json:"showLineNumbers,omitzero"`
	HideTips                          bool                       `json:"hideTips,omitzero"`
	ShowShortcutsHint                 bool                       `json:"showShortcutsHint,omitzero"`
	HideBanner                        bool                       `json:"hideBanner,omitzero"`
	DynamicWindowTitle                bool                       `json:"dynamicWindowTitle,omitzero"`
	ShowStatusInTitle                 bool                       `json:"showStatusInTitle,omitzero"`
	HideFooter                        bool                       `json:"hideFooter,omitzero"`
	ShowMemoryUsage                   bool                       `json:"showMemoryUsage,omitzero"`
	ShowHomeDirectoryWarning          bool                       `json:"showHomeDirectoryWarning,omitzero"`
	ShowCitations                     bool                       `json:"showCitations,omitzero"`
	ShowModelInfoInChat               bool                       `json:"showModelInfoInChat,omitzero"`
	ShowUserIdentity                  bool                       `json:"showUserIdentity,omitzero"`
	UseAlternateBuffer                bool                       `json:"useAlternateBuffer,omitzero"`
	UseBackgroundColor                bool                       `json:"useBackgroundColor,omitzero"`
	IncrementalRendering              bool                       `json:"incrementalRendering,omitzero"`
	ShowSpinner                       bool                       `json:"showSpinner,omitzero"`
	HideWindowTitle                   bool                       `json:"hideWindowTitle,omitzero"`
	AutoThemeSwitching                bool                       `json:"autoThemeSwitching,omitzero"`
}

type GeminiUIFooterConfig struct {
	HideCWD               bool `json:"hideCWD,omitzero"`
	HideSandboxStatus     bool `json:"hideSandboxStatus,omitzero"`
	HideModelInfo         bool `json:"hideModelInfo,omitzero"`
	HideContextPercentage bool `json:"hideContextPercentage,omitzero"`
}

type GeminiAccessibilityConfig struct {
	EnableLoadingPhrases bool `json:"enableLoadingPhrases,omitzero"`
	ScreenReader         bool `json:"screenReader,omitzero"`
}

type GeminiIDEConfig struct {
	Enabled      bool `json:"enabled,omitzero"`
	HasSeenNudge bool `json:"hasSeenNudge,omitzero"`
}

type GeminiPrivacyConfig struct {
	UsageStatisticsEnabled bool `json:"usageStatisticsEnabled,omitzero"`
}

type GeminiModelConfig struct {
	SummarizeToolOutput  map[string]any `json:"summarizeToolOutput,omitzero"`
	Name                 string         `json:"name,omitzero"`
	MaxSessionTurns      int            `json:"maxSessionTurns,omitzero"`
	CompressionThreshold float64        `json:"compressionThreshold,omitzero"`
	DisableLoopDetection bool           `json:"disableLoopDetection,omitzero"`
	SkipNextSpeakerCheck bool           `json:"skipNextSpeakerCheck,omitzero"`
}

type GeminiModelConfigsConfig struct {
	Aliases         map[string]any `json:"aliases,omitzero"`
	CustomAliases   map[string]any `json:"customAliases,omitzero"`
	CustomOverrides []any          `json:"customOverrides,omitzero"`
	Overrides       []any          `json:"overrides,omitzero"`
}

type GeminiAgentsConfig struct {
	Overrides map[string]any `json:"overrides,omitzero"`
}

type GeminiContextConfig struct {
	FileName                         any                        `json:"fileName,omitzero"`
	FileFiltering                    *GeminiFileFilteringConfig `json:"fileFiltering,omitzero"`
	ImportFormat                     string                     `json:"importFormat,omitzero"`
	IncludeDirectories               []string                   `json:"includeDirectories,omitzero"`
	DiscoveryMaxDirs                 int                        `json:"discoveryMaxDirs,omitzero"`
	IncludeDirectoryTree             bool                       `json:"includeDirectoryTree,omitzero"`
	LoadMemoryFromIncludeDirectories bool                       `json:"loadMemoryFromIncludeDirectories,omitzero"`
}

type GeminiFileFilteringConfig struct {
	CustomIgnoreFilePaths     []string `json:"customIgnoreFilePaths,omitzero"`
	RespectGitIgnore          bool     `json:"respectGitIgnore,omitzero"`
	RespectGeminiIgnore       bool     `json:"respectGeminiIgnore,omitzero"`
	EnableRecursiveFileSearch bool     `json:"enableRecursiveFileSearch,omitzero"`
	EnableFuzzySearch         bool     `json:"enableFuzzySearch,omitzero"`
}

type GeminiToolsConfig struct {
	Sandbox                     any                `json:"sandbox,omitzero"`
	Shell                       *GeminiShellConfig `json:"shell,omitzero"`
	DiscoveryCommand            string             `json:"discoveryCommand,omitzero"`
	CallCommand                 string             `json:"callCommand,omitzero"`
	Core                        []string           `json:"core,omitzero"`
	Allowed                     []string           `json:"allowed,omitzero"`
	Exclude                     []string           `json:"exclude,omitzero"`
	TruncateToolOutputThreshold int                `json:"truncateToolOutputThreshold,omitzero"`
	UseRipgrep                  bool               `json:"useRipgrep,omitzero"`
	DisableLLMCorrection        bool               `json:"disableLLMCorrection,omitzero"`
}

type GeminiShellConfig struct {
	Pager                       string `json:"pager,omitzero"`
	InactivityTimeout           int    `json:"inactivityTimeout,omitzero"`
	EnableInteractiveShell      bool   `json:"enableInteractiveShell,omitzero"`
	ShowColor                   bool   `json:"showColor,omitzero"`
	EnableShellOutputEfficiency bool   `json:"enableShellOutputEfficiency,omitzero"`
}

type GeminiMCPConfig struct {
	ServerCommand string   `json:"serverCommand,omitzero"`
	Allowed       []string `json:"allowed,omitzero"`
	Excluded      []string `json:"excluded,omitzero"`
}

type GeminiSecurityConfig struct {
	FolderTrust                  *GeminiFolderTrustConfig  `json:"folderTrust,omitzero"`
	EnvironmentVariableRedaction *GeminiEnvRedactionConfig `json:"environmentVariableRedaction,omitzero"`
	Auth                         *GeminiAuthConfig         `json:"auth,omitzero"`
	AllowedExtensions            []string                  `json:"allowedExtensions,omitzero"`
	DisableYoloMode              bool                      `json:"disableYoloMode,omitzero"`
	EnablePermanentToolApproval  bool                      `json:"enablePermanentToolApproval,omitzero"`
	BlockGitExtensions           bool                      `json:"blockGitExtensions,omitzero"`
}

type GeminiFolderTrustConfig struct {
	Enabled bool `json:"enabled,omitzero"`
}

type GeminiEnvRedactionConfig struct {
	Allowed []string `json:"allowed,omitzero"`
	Blocked []string `json:"blocked,omitzero"`
	Enabled bool     `json:"enabled,omitzero"`
}

type GeminiAuthConfig struct {
	SelectedType string `json:"selectedType,omitzero"`
	EnforcedType string `json:"enforcedType,omitzero"`
	UseExternal  bool   `json:"useExternal,omitzero"`
}

type GeminiAdvancedConfig struct {
	BugCommand          map[string]any `json:"bugCommand,omitzero"`
	DNSResolutionOrder  string         `json:"dnsResolutionOrder,omitzero"`
	ExcludedEnvVars     []string       `json:"excludedEnvVars,omitzero"`
	AutoConfigureMemory bool           `json:"autoConfigureMemory,omitzero"`
}

type GeminiExperimentalConfig struct {
	ToolOutputMasking   *GeminiToolOutputMaskingConfig `json:"toolOutputMasking,omitzero"`
	EnableAgents        bool                           `json:"enableAgents,omitzero"`
	ExtensionManagement bool                           `json:"extensionManagement,omitzero"`
	ExtensionConfig     bool                           `json:"extensionConfig,omitzero"`
	ExtensionRegistry   bool                           `json:"extensionRegistry,omitzero"`
	ExtensionReloading  bool                           `json:"extensionReloading,omitzero"`
	JitContext          bool                           `json:"jitContext,omitzero"`
	UseOSC52Paste       bool                           `json:"useOSC52Paste,omitzero"`
	Plan                bool                           `json:"plan,omitzero"`
}

type GeminiToolOutputMaskingConfig struct {
	ToolProtectionThreshold    int  `json:"toolProtectionThreshold,omitzero"`
	MinPrunableTokensThreshold int  `json:"minPrunableTokensThreshold,omitzero"`
	Enabled                    bool `json:"enabled,omitzero"`
	ProtectLatestTurn          bool `json:"protectLatestTurn,omitzero"`
}

type GeminiSkillsConfig struct {
	Disabled []string `json:"disabled,omitzero"`
	Enabled  bool     `json:"enabled,omitzero"`
}

type GeminiHooksConfigConfig struct {
	Disabled      []string `json:"disabled,omitzero"`
	Enabled       bool     `json:"enabled,omitzero"`
	Notifications bool     `json:"notifications,omitzero"`
}

type GeminiHooksConfig struct {
	BeforeTool          []any `json:"BeforeTool,omitzero"`
	AfterTool           []any `json:"AfterTool,omitzero"`
	BeforeAgent         []any `json:"BeforeAgent,omitzero"`
	AfterAgent          []any `json:"AfterAgent,omitzero"`
	Notification        []any `json:"Notification,omitzero"`
	SessionStart        []any `json:"SessionStart,omitzero"`
	SessionEnd          []any `json:"SessionEnd,omitzero"`
	PreCompress         []any `json:"PreCompress,omitzero"`
	BeforeModel         []any `json:"BeforeModel,omitzero"`
	AfterModel          []any `json:"AfterModel,omitzero"`
	BeforeToolSelection []any `json:"BeforeToolSelection,omitzero"`
}

type GeminiAdminConfig struct {
	Extensions        *GeminiAdminExtConfig    `json:"extensions,omitzero"`
	MCP               *GeminiAdminMCPConfig    `json:"mcp,omitzero"`
	Skills            *GeminiAdminSkillsConfig `json:"skills,omitzero"`
	SecureModeEnabled bool                     `json:"secureModeEnabled,omitzero"`
}

type GeminiAdminExtConfig struct {
	Enabled bool `json:"enabled,omitzero"`
}

type GeminiAdminMCPConfig struct {
	Config  map[string]any `json:"config,omitzero"`
	Enabled bool           `json:"enabled,omitzero"`
}

type GeminiAdminSkillsConfig struct {
	Enabled bool `json:"enabled,omitzero"`
}

type GeminiMCPServer struct {
	Env          map[string]string `json:"env,omitzero"`
	Headers      map[string]string `json:"headers,omitzero"`
	Command      string            `json:"command,omitzero"`
	Cwd          string            `json:"cwd,omitzero"`
	URL          string            `json:"url,omitzero"`
	HTTPURL      string            `json:"httpUrl,omitzero"`
	Description  string            `json:"description,omitzero"`
	Args         []string          `json:"args,omitzero"`
	IncludeTools []string          `json:"includeTools,omitzero"`
	ExcludeTools []string          `json:"excludeTools,omitzero"`
	Timeout      int               `json:"timeout,omitzero"`
	Trust        bool              `json:"trust,omitzero"`
}

type GeminiTelemetryConfig struct {
	Target       string `json:"target,omitzero"`
	OTLPEndpoint string `json:"otlpEndpoint,omitzero"`
	OTLPProtocol string `json:"otlpProtocol,omitzero"`
	Outfile      string `json:"outfile,omitzero"`
	Enabled      bool   `json:"enabled,omitzero"`
	LogPrompts   bool   `json:"logPrompts,omitzero"`
	UseCollector bool   `json:"useCollector,omitzero"`
}
