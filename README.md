# Nordlys Installer

CLI to configure AI coding tools (Claude Code, OpenCode, Codex, Gemini CLI, Grok CLI, Qwen Code, Zed) with Nordlys. Preserves existing configs and creates backups before changes.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/nordlys-labs/nordlys-installer/main/scripts/install.sh | bash
```

## Usage

Interactive (default):

```bash
nordlys-installer
```

Non-interactive:

```bash
nordlys-installer --api-key "your-key" --tools claude-code,opencode
# or: export NORDLYS_API_KEY="your-key" && nordlys-installer --non-interactive
```

## Commands

| Command | Description |
|---------|-------------|
| `nordlys-installer` | Interactive setup |
| `list` | List supported tools |
| `validate` | Check API key and configs |
| `update` | Self-update installer |
| `uninstall [tool...]` | Remove Nordlys from tools |
| `<tool> update` | Re-apply config for one tool |
| `version` | Print version |

## Tools

| Tool | Config | Node.js |
|------|--------|---------|
| Claude Code | `~/.claude/settings.json` | yes |
| OpenCode | `~/.config/opencode/opencode.json` | yes |
| Codex | `~/.codex/config.json` | yes |
| Gemini CLI | `~/.gemini/settings.json` | yes |
| Grok CLI | `~/.grok/config.json` | yes |
| Qwen Code | `~/.qwen/config.json` | yes |
| Zed | `~/.config/zed/settings.json` | no |

## Config

- API key: `--api-key` flag, then `NORDLYS_API_KEY` env, then interactive prompt
- Model: `--model` (default: `nordlys/hypernova`)
- Use `<tool> update` to refresh config when formats change

## Build

```bash
go build -o nordlys-installer ./cmd/nordlys-installer
go test ./...
```

Go 1.22+. Linux, macOS, Windows (amd64, arm64).

## Links

- [Docs](https://docs.nordlyslabs.com)
- [Issues](https://github.com/nordlys-labs/nordlys-installer/issues)
