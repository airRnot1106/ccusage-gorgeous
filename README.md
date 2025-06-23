# ccusage-gorgeous 🌈

A beautiful terminal user interface (TUI) application that displays Claude API usage costs as stunning ASCII art with animated rainbow colors.

## ✨ Features

- **ASCII Art Display**: Transform cost data into large, eye-catching ASCII art
- **Rainbow Animations**: 4 different animation patterns (rainbow, gradient, pulse, wave)
- **Responsive Design**: Automatically adapts font size to terminal dimensions
- **Real-time Updates**: Live cost data fetching from ccusage CLI
- **Clean Architecture**: Plugin-based system with dependency inversion
- **Modern Development**: Built with Nix flakes for reproducible development

## 🚀 Quick Start

### Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled
- [ccusage](https://github.com/ccusage/ccusage) CLI tool accessible via npx

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

### Controls

- `r` - Refresh cost data
- `q` or `Ctrl+C` - Quit application

## 🏗️ Architecture

ccusage-gorgeous follows **Clean Architecture** principles with a sophisticated plugin system:

```
┌──────────────────────────────────────────────────────┐
│                   Presentation Layer                 │
│  ┌─────────────────────────────────────────────────┐ │
│  │          TUI (Bubbletea + Lipgloss)             │ │
│  └─────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│                 Application Layer                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │
│  │  Plugin     │ │ Animation   │ │    Display      │ │
│  │ Interfaces  │ │ Interfaces  │ │   Interfaces    │ │
│  └─────────────┘ └─────────────┘ └─────────────────┘ │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│                   Domain Layer                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │
│  │  CostData   │ │ Animation   │ │ DisplayConfig   │ │
│  │   Models    │ │   Models    │ │    Models       │ │
│  └─────────────┘ └─────────────┘ └─────────────────┘ │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│               Infrastructure Layer                   │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐ │
│  │  ccusage    │ │  Rainbow    │ │   ASCII Art     │ │
│  │    CLI      │ │ Animation   │ │   Rendering     │ │
│  └─────────────┘ └─────────────┘ └─────────────────┘ │
└──────────────────────────────────────────────────────┘
```

### Plugin System

The application uses a **plugin-based architecture** with three types of plugins:

1. **DataSource Plugins**: Fetch cost data from various sources
   - `ccusage-cli`: Integrates with ccusage CLI tool
   - `bankruptcy`: Mock data source for testing

2. **Animation Plugins**: Generate animated color effects
   - `rainbow`: Cycling rainbow colors
   - `gradient`: Smooth color transitions
   - `pulse`: Pulsing color effects
   - `wave`: Sine wave patterns

3. **Display Plugins**: Render visual output
   - `rainbow-tui`: ASCII art with rainbow colors

## 📁 Project Structure

```
ccusage-gorgeous/
├── main.go                    # Application entry point
├── configs/
│   └── config.yaml           # Configuration file
├── internal/
│   ├── domain/               # Business entities (no external deps)
│   │   ├── cost.go          # Cost data models
│   │   ├── animation.go     # Animation models
│   │   ├── display.go       # Display models
│   │   └── errors.go        # Domain errors
│   ├── application/         # Use cases and interfaces
│   │   └── interfaces/      # Plugin interfaces
│   ├── core/               # Core infrastructure
│   │   ├── config.go       # Configuration management
│   │   └── registry.go     # Plugin registry
│   ├── infrastructure/     # External concerns
│   │   └── tui/           # Terminal UI implementation
│   └── plugins/           # Plugin implementations
│       ├── datasource/    # Data source plugins
│       ├── animation/     # Animation plugins
│       └── display/       # Display plugins
├── test/                  # Comprehensive test suite
├── nix/                   # Nix configuration
└── flake.nix             # Nix flake definition
```

## ⚙️ Configuration

Configuration is managed through a hierarchical system:

1. **Command-line flags** (highest priority)
2. **Environment variables** (`CCUSAGE_*`)
3. **Configuration files** (`./configs/config.yaml`, `~/.ccusage-gorgeous/config.yaml`)
4. **Built-in defaults** (lowest priority)

### Configuration Options

```yaml
app:
  log_level: "info"
  refresh_rate: "1s"

display:
  format: "large"          # large, medium, small, minimal
  width: 80               # Terminal width
  height: 24              # Terminal height
  show_timestamp: true
  show_breakdown: true

animation:
  enabled: true
  speed: "100ms"          # Animation speed
  pattern: "rainbow"      # rainbow, gradient, pulse, wave
  colors:                 # Custom color palette
    - "#FF0000"          # Red
    - "#FF8000"          # Orange
    - "#FFFF00"          # Yellow
    # ... 12 colors total

datasource:
  ccusage_path: "ccusage" # Path to ccusage CLI
  timeout: "30s"         # Command timeout
  cache_time: "10s"      # Cache duration

plugins:
  datasource: "ccusage-cli"
  display: "rainbow-display"
  animation: "rainbow-animation"
```

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
| `go test ./...` | Run all tests |
| `go build -o ccugorg` | Local build |

### Code Quality

The project includes comprehensive code quality tools:

- **golangci-lint**: Go linting
- **gofumpt**: Go code formatting
- **goimports**: Import organization
- **nixfmt**: Nix code formatting
- **Pre-commit hooks**: Automatic quality checks

### Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run specific test package
go test ./test/plugins/display/

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test ./test/integration/
```

### Test-Driven Development (TDD)

This project follows TDD practices:

1. **Write failing tests first**
2. **Implement minimal code to pass tests**
3. **Refactor while keeping tests green**
4. **Repeat the Red-Green-Refactor cycle**

## 🎨 ASCII Art System

The ASCII art generation system supports:

### Font Sizes

- **Large fonts** (10 rows): For terminals ≥100x25
- **Small fonts** (7 rows): For smaller terminals

### Character Support

Complete numeric display with:
- Digits: `0-9`
- Currency: `$`
- Decimal: `.`

### Features

- **Responsive sizing**: Automatic font selection
- **Character spacing**: 2-space gaps between characters
- **Perfect centering**: Both horizontal and vertical
- **Unicode blocks**: Solid appearance using `█` characters

## 🌈 Animation System

### Animation Patterns

1. **Rainbow**: Classic cycling rainbow colors
2. **Gradient**: Smooth positional color transitions
3. **Pulse**: Sine-wave based color pulsing
4. **Wave**: Traveling wave pattern across text

### Technical Implementation

Colors are applied per-character using mathematical functions:

```go
// Wave pattern example
waveValue := math.Sin(float64(frameNumber)*0.1 + float64(i)*0.5)
colorIndex := int((waveValue+1)/2*float64(len(colors))) % len(colors)
```

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

3. **Add configuration**:
```yaml
plugins:
  datasource: "my-plugin"
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

## 📋 Advanced Usage

### Environment Variables

```bash
# Set log level
export CCUSAGE_LOG_LEVEL=debug

# Custom ccusage path
export CCUSAGE_CLI_PATH=/usr/local/bin/ccusage

# Animation speed
export CCUSAGE_ANIMATION_SPEED=50ms
```

### Custom Configuration

Create `~/.ccusage-gorgeous/config.yaml`:

```yaml
animation:
  pattern: "pulse"
  speed: "50ms"
  colors:
    - "#FF6B6B"
    - "#4ECDC4"
    - "#45B7D1"
    - "#96CEB4"
```

## 📄 License

MIT

## 🙏 Acknowledgments

- [ccusage](https://github.com/ccusage/ccusage) - Claude usage tracking

**ccusage-gorgeous** - Making Claude API cost tracking beautiful, one rainbow at a time! 🌈✨
