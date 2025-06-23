# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

This is a Go project managed with Nix Flakes. The development environment includes Go 1.24, linting, formatting, and pre-commit hooks.

### Essential Commands

- `nix develop` - Enter the development shell with all tools and dependencies
- `nix build` - Build the Go application (outputs binary as `ccugorg`)
- `nix fmt` - Format all code (Go files with gofumpt, Nix files with nixfmt)

### Go Development

Within the Nix development shell, standard Go commands work:
- `go build -o ccugorg` - Build the application with specific binary name
- `go test ./...` - Run all tests
- `go test ./test/plugins/datasource/` - Run specific test package
- `go mod tidy` - Clean up module dependencies

### Code Quality

Pre-commit hooks are automatically installed when entering the development shell:
- `golangci-lint` - Go linting
- `gofumpt` - Go code formatting
- `nixfmt` - Nix code formatting

Manual formatting: `nix fmt` (formats both Go and Nix files)

## Build and Test Commands

**IMPORTANT:** Use these specific commands for building and testing:

### Build
- `nix build` - Build the Go application (outputs binary as `ccugorg`)

### Test
- `nix flake check --no-pure-eval` - Run all tests, verify flake configuration and build

## Development Methodology

**Test-Driven Development (TDD)**: This project follows TDD practices:

1. **Write failing tests first** - Before implementing any new feature or making changes
2. **Make tests pass** - Implement minimal code to make the failing tests pass
3. **Refactor** - Improve code quality while keeping tests green
4. **Repeat** - Continue the Red-Green-Refactor cycle

## Task Completion Verification

**IMPORTANT:** After completing any task, always run these commands in order:
1. `nix fmt` - Format all code (Go and Nix files)
2. `nix flake check --no-pure-eval` - Verify flake configuration and build

These commands ensure code quality and flake integrity are maintained.

## Project Architecture

This is a TUI application that displays ccusage cost data as large ASCII art with rainbow animations. The application follows Clean Architecture principles with a plugin-based system.

### Core Purpose
The application transforms ccusage cost data into visually appealing ASCII art displaying only the total cost in $99.99 format with rainbow color animation. All other UI elements have been removed for a clean, focused display.

### Clean Architecture Layers

**Domain Layer** (`internal/domain/`):
- Core business entities: `CostData`, `AnimationFrame`, `DisplayConfig`
- No external dependencies, contains pure business logic

**Application Layer** (`internal/application/`):
- `interfaces/` - Defines plugin interfaces (`Plugin`, `DataSourcePlugin`, `DisplayPlugin`, `AnimationPlugin`)
- `usecases/` - Business logic orchestration (currently minimal)

**Infrastructure Layer** (`internal/infrastructure/`):
- `tui/` - Bubbletea TUI implementation and models
- External integrations and framework-specific code

### Plugin System Architecture

**Core Plugin Registry** (`internal/core/registry.go`):
- Centralized plugin management and dependency injection
- Runtime plugin registration and lifecycle management

**Plugin Types**:
- **DataSource** (`internal/plugins/datasource/`) - ccusage CLI integration via npx
- **Display** (`internal/plugins/display/`) - ASCII art generation with responsive font sizing
- **Animation** (`internal/plugins/animation/`) - Rainbow color cycling effects

**Plugin Configuration**:
- Viper-based configuration management (`internal/core/config.go`)
- YAML configuration file (`configs/config.yaml`)
- Plugin-specific settings and runtime parameters

### Key Implementation Details

**ASCII Art Display System**:
- Responsive font selection: small (7-row) fonts for screens <100 width or <25 height, large (10-row) fonts otherwise
- Character spacing: 2-space gaps between ASCII letters for readability
- Centered display: both horizontal and vertical centering within terminal dimensions
- Complete character set: digits 0-9, dollar sign ($), and decimal point (.)

**ccusage Integration**:
- Uses `npx ccusage daily --json` for data fetching
- Parses JSON structure: `{daily: [...], totals: {totalCost, inputTokens, outputTokens, modelBreakdowns}}`
- Implements caching and timeout handling for performance

**TUI Framework**:
- Built on Charmbracelet Bubbletea for terminal interface
- Lipgloss for styling and rainbow color application
- HSL-based color cycling for smooth rainbow animations
- Clean startup: no loading text or control instructions displayed

**Testing Strategy**:
- Comprehensive unit tests in `test/` directory mirroring `internal/` structure
- Integration tests for plugin system and end-to-end workflows
- Tests validate ASCII art generation (checking for `█` block characters)
- Test coverage for all plugin implementations and display formats

### Configuration Structure

**Nix Configuration**:
- `flake.nix` - Main flake with `buildGoModule`, renames binary to `ccugorg`
- `nix/pre-commit/default.nix` - Pre-commit hook configuration
- `nix/treefmt/default.nix` - Code formatting rules

**Application Configuration** (`configs/config.yaml`):
- Display settings: format, dimensions, show options
- Animation parameters: speed (100ms), pattern (rainbow), 12-color rainbow spectrum
- Datasource configuration: ccusage path, timeout (30s), caching (10s)
- Plugin selection: specifies active datasource, display, and animation plugins
