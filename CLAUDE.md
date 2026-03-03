# Nordlys Installer - Claude Code Instructions

Instructions for Claude when editing this codebase.

## Project Summary

CLI that configures AI coding tools to use Nordlys. Written in Go 1.26. Uses Cobra for CLI, Bubble Tea for TUI.

## Codebase Layout

- `cmd/nordlys-installer/main.go` - Cobra root, commands, flag wiring
- `internal/tools/` - Tool interface + 7 implementations (claude_code, opencode, codex, gemini_cli, grok_cli, qwen_code, zed)
- `internal/config/` - JSON/TOML I/O, backups, API validation
- `internal/ui/` - TUI (Bubble Tea)
- `internal/updater/` - Self-update from GitHub releases
- `internal/runtime/` - Node.js check, platform detection
- `internal/constants/` - Version, API URLs

## Tool Implementation Checklist

When adding or changing a tool:

1. Implement all `Tool` interface methods in `internal/tools/tool.go`.
2. `ConfigPath()` - return config file path; handle `os.UserHomeDir()` errors.
3. `UpdateConfig()` - merge Nordlys settings into existing config; use `config.UpdateJSONFields` or equivalent for JSON.
4. `Uninstall()` - remove Nordlys keys only; preserve user settings.
5. `GetExistingConfig()` - read apiKey/model from config; return `"", ""` on error.
6. `Validate()` - check API key format and optionally connection.
7. Add `runToolUpdate` subcommand in main.go for `<tool> update`.
8. Add `*_test.go` with `Test<Name>_Validate`, `Test<Name>_Uninstall`, `Test<Name>_FullWorkflow`.

## Strict Go Practices

- **Errors**: Never `_ =` or ignore. Always check and return/propagate.
- **Shadow**: Use `err = x` when `err` exists in scope; avoid `if err := x`.
- **Field alignment**: Run `fieldalignment -fix` if govet complains.
- **Imports**: Group: stdlib, external, local (see `.golangci.yml` goimports).
- **Linters**: `golangci-lint run` must pass. No disabling errcheck, fieldalignment, etc.

## Testing

```bash
go test ./...
```

- Use `t.TempDir()` for test files.
- Mock HTTP in validator tests via `config.ValidatorHTTPClient`.
- Tool tests: create temp config, call UpdateConfig/Validate/Uninstall, assert file state.

## Common Tasks

- **Add tool**: New file in `internal/tools/`, types in `types.go`, register in `registry.go`, subcommand in main.
- **Change config format**: Update tool's `UpdateConfig` and `Uninstall`; adjust types.
- **New CLI flag**: Add to `main.go` `var` block, bind in `init()`, use in `Run`.

## References

See `AGENTS.md` for full architecture and patterns.
