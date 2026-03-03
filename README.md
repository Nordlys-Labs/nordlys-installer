# Nordlys Installer

A cross-platform CLI tool to configure developer tools with Nordlys's Mixture of Models.

## Features

- 🚀 **Interactive TUI** - Beautiful terminal UI powered by Bubbletea
- 🔧 **7 Supported Tools** - Claude Code, OpenCode, Codex, Gemini CLI, Grok CLI, Qwen Code, Zed
- 🔄 **Auto-Update** - Self-update capability to latest version
- 💾 **Safe Updates** - Automatic backups before modifying configs
- 🎯 **Non-Interactive Mode** - Perfect for CI/CD and automation
- ✅ **Validation** - API key and connection validation
- 🗑️ **Uninstall** - Clean removal of Nordlys configuration

## Quick Start

### Installation

```bash
curl -fsSL https://raw.githubusercontent.com/nordlys-labs/nordlys-installer/main/scripts/install.sh | bash
```

### Usage

#### Interactive Mode (Recommended)
```bash
nordlys-installer
```

#### Non-Interactive Mode
```bash
# Configure specific tools
nordlys-installer --api-key "your-key" --tools claude-code,opencode

# Configure all tools
nordlys-installer --api-key "your-key"

# Use environment variable
export NORDLYS_API_KEY="your-key"
nordlys-installer --non-interactive
```

## Commands

```bash
nordlys-installer              # Interactive mode
nordlys-installer list         # List all supported tools
nordlys-installer validate     # Validate configuration
nordlys-installer update       # Update nordlys-installer to latest version
nordlys-installer uninstall    # Remove Nordlys config from tools
nordlys-installer version      # Show version
nordlys-installer --help       # Show help

# Update Nordlys config for a specific tool (re-applies config when formats change)
nordlys-installer claude-code update
nordlys-installer zed update
# ... same for opencode, codex, gemini-cli, grok-cli, qwen-code
```

## Supported Tools

| Tool | Config Path | Requires Node.js |
|------|-------------|------------------|
| Claude Code | `~/.claude/settings.json` | ✅ |
| OpenCode | `~/.config/opencode/opencode.json` | ✅ |
| Codex | `~/.codex/config.json` | ✅ |
| Gemini CLI | `~/.gemini/settings.json` | ✅ |
| Grok CLI | `~/.grok/config.json` | ✅ |
| Qwen Code | `~/.qwen/config.json` | ✅ |
| Zed | `~/.config/zed/settings.json` | ❌ |

## Configuration

The installer only updates Nordlys-specific fields and preserves your existing configuration. Automatic backups are created with timestamps before any modifications.

Use `nordlys-installer <tool> update` to re-apply Nordlys configuration to an already-configured tool. This is useful when config formats change in new versions.

### API Key Priority

1. `--api-key` CLI flag
2. `NORDLYS_API_KEY` environment variable
3. Interactive prompt

### Model Override

By default, the installer uses `nordlys/hypernova`. To use a different model:

```bash
nordlys-installer --api-key "key" --model "nordlys/custom-model"
```

## Development

### Prerequisites

- Go 1.22+
- Make (optional)

### Build

```bash
go build -o nordlys-installer ./cmd/nordlys-installer
```

### Test

```bash
go test ./...
```

## Platform Support

- ✅ Linux (amd64, arm64)
- ✅ macOS (Intel, Apple Silicon)
- ✅ Windows (amd64, arm64)

## License

See [LICENSE](LICENSE) file for details.

## Support

- 📖 Documentation: https://docs.nordlyslabs.com
- 🐛 Issues: https://github.com/nordlys-labs/nordlys-installer/issues
- 💬 Support: support@llmadaptive.uk
