# Nordlys Installer - Agent Guide

Instructions for AI agents working on this codebase.

## Overview

Nordlys Installer is a CLI that configures AI coding tools (Claude Code, OpenCode, Codex, Gemini CLI, Grok CLI, Qwen Code, Zed) to use Nordlys. It preserves existing configs, creates backups before changes, and supports interactive and non-interactive modes.

## Architecture

```
cmd/nordlys-installer/     # Entry point, Cobra commands
internal/
  config/                  # JSON/TOML read/write, backup, validation
  constants/               # Version, API URLs, defaults
  runtime/                 # Node.js detection, platform helpers
  tools/                   # Tool implementations (Tool interface)
  ui/                      # Bubble Tea TUI
  updater/                 # Self-update from GitHub releases
scripts/                   # Install script
```

## Key Patterns

### Tool Interface

All tools implement `internal/tools/tool.go`:

```go
type Tool interface {
    Name() string
    Description() string
    ConfigPath() (string, error)
    IsInstalled() bool
    RequiresNode() bool
    UpdateConfig(apiKey, model, baseURL string) error
    Validate() error
    Uninstall() error
    GetExistingConfig() (apiKey, model string)
}
```

- **ConfigPath**: Returns path to tool config file (JSON or TOML).
- **UpdateConfig**: Writes Nordlys API key, model, base URL. Preserves other settings.
- **Uninstall**: Removes Nordlys-specific config, keeps user settings.
- **GetExistingConfig**: Used by `update` subcommand to avoid requiring flags when config exists.

### Adding a New Tool

1. Create `internal/tools/<name>.go` implementing `Tool`.
2. Add config types to `internal/tools/types.go` if needed.
3. Register in `internal/tools/registry.go` (`GetAllTools`).
4. Add Cobra subcommand in `cmd/nordlys-installer/main.go` for `<tool> update`.
5. Add tests in `internal/tools/<name>_test.go`.

### Config Handling

- Use `internal/config` for JSON/TOML: `ReadJSONFile`, `WriteJSONFile`, `ReadTOMLFile`, `WriteTOMLFile`, `UpdateJSONFields`.
- Create backups with `config.CreateBackup(path)` before modifying files.
- Validate configs against schemas via `tools.ValidateConfig(schemaURL, data)` when available.

### Error Handling

- Never ignore errors. Check all `err` returns.
- Use `err = x` (not `err := x`) when `err` is already in scope to avoid shadowing.
- Return errors with context: `fmt.Errorf("action: %w", err)`.

## Conventions

- **Go 1.26+**. Standard library and common idioms.
- **Imports**: stdlib, blank line, external packages, blank line, local (`github.com/nordlys-labs/nordlys-installer`).
- **Linting**: `golangci-lint run` must pass. No disabled linters. Fix fieldalignment, shadow, errcheck.
- **Tests**: `go test ./...`. Use `t.Parallel()` where safe. Prefer table-driven tests.
- **No emojis** in code, CLI output, or docs.

## Commands

| Command | Description |
|---------|-------------|
| (default) | Interactive TUI |
| `list` | List supported tools |
| `validate` | Check API key and tool configs |
| `update` | Self-update installer |
| `uninstall [tool...]` | Remove Nordlys from tools |
| `<tool> update` | Re-apply config for one tool |
| `version` | Print version |

## Build & Test

```bash
go build -o nordlys-installer ./cmd/nordlys-installer
go test ./...
golangci-lint run
```

## Links

- [Docs](https://docs.nordlyslabs.com)
- [Issues](https://github.com/nordlys-labs/nordlys-installer/issues)
