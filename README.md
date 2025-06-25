# ccusage-gorgeous 🌈

A beautiful terminal user interface application that displays Claude Code usage costs obtained from ccusage as stunning ASCII art with animated rainbow colors.

https://github.com/user-attachments/assets/c96a2147-ef59-46bc-9e1b-f064f34a8f64

## 🚀 Quick Start

### Installation & Usage

```bash
# Clone the repository
git clone https://github.com/airRnot1106/ccusage-gorgeous
cd ccusage-gorgeous

# Enter development environment
nix develop

# Build the application
nix build

# Run ccusage-gorgeous
./result/bin/ccugorg
```

### CLI Options

```bash
# View help
ccugorg --help

# Set animation speed
ccugorg --animation-speed 50ms

# Change animation pattern
ccugorg --animation-pattern pulse

# Disable animation
ccugorg --no-animation

# Combine options
ccugorg --animation-speed 200ms --animation-pattern wave
```

<details>
<summary>Demo</summary>

https://github.com/user-attachments/assets/69955a3f-274c-4eaa-a3eb-880f5d724486

https://github.com/user-attachments/assets/c2c65536-1f68-4940-a778-abb8a181b7bd

</details>

## 🛠️ Development

### Development Environment

This project uses **Nix flakes** for reproducible development:

```bash
# Enter development shell with all dependencies
nix develop

# Format code (Go + Nix)
nix fmt

# Run tests and checks
nix flake check --no-pure-eval

# Build application
nix build
```

### Essential Commands

| Command | Description |
|---------|-------------|
| `nix develop` | Enter development shell |
| `nix build` | Build the application |
| `nix fmt` | Format all code |
| `nix flake check --no-pure-eval` | Run all tests |
| `nix build` | Local build |

### Code Quality

The project includes comprehensive code quality tools:

- **golangci-lint**: Go linting
- **gofumpt**: Go code formatting
- **nixfmt**: Nix code formatting
- **Pre-commit hooks**: Automatic quality checks

## 🔧 Plugin Development

### Creating a New Plugin

1. **Implement the interface**:
```go
type MyPlugin struct {
    name        string
    version     string
    description string
    enabled     bool
}

func (p *MyPlugin) Initialize(config map[string]interface{}) error {
    p.enabled = true
    return nil
}
```

2. **Register the plugin**:
```go
plugin := NewMyPlugin()
registry.RegisterDataSource(plugin)
```

3. **Update plugin selection**:
```go
// In main application logic
configManager.UpdateConfig(map[string]interface{}{
    "plugins.datasource": "my-plugin",
})
```

### Plugin Interfaces

All plugins implement the base `Plugin` interface:

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Initialize(config map[string]interface{}) error
    Shutdown() error
    IsEnabled() bool
}
```

Specialized interfaces extend the base:
- `DataSourcePlugin`: Data fetching
- `AnimationPlugin`: Animation generation
- `DisplayPlugin`: Visual rendering

## 📄 License

MIT

## 🙏 Acknowledgments

- [ccusage](https://github.com/ccusage/ccusage) - Claude usage tracking
