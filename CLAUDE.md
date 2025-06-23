# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

This is a Go project managed with Nix Flakes. The development environment includes Go 1.24.3, linting, formatting, and pre-commit hooks.

### Essential Commands

- `nix develop` - Enter the development shell with all tools and dependencies
- `nix build` - Build the Go application (outputs to `./result/bin/ccusage-gorgeous`)
- `nix fmt` - Format all code (Go files with gofumpt, Nix files with nixfmt)

### Go Development

Within the Nix development shell, standard Go commands work:
- `go build` - Build the application
- `go test ./...` - Run all tests
- `go mod tidy` - Clean up module dependencies

### Code Quality

Pre-commit hooks are automatically installed when entering the development shell:
- `golangci-lint` - Go linting
- `gofumpt` - Go code formatting  
- `nixfmt` - Nix code formatting

Manual formatting: `nix fmt` (formats both Go and Nix files)

## Task Completion Verification

**IMPORTANT:** After completing any task, always run these commands in order:
1. `nix fmt` - Format all code (Go and Nix files)
2. `nix flake check --no-pure-eval` - Verify flake configuration and build

These commands ensure code quality and flake integrity are maintained.

## Project Architecture

**Nix Configuration Structure:**
- `flake.nix` - Main flake configuration with package definition and development shell
- `nix/pre-commit/default.nix` - Pre-commit hook configuration
- `nix/treefmt/default.nix` - Code formatting rules

**Package Configuration:**
The project uses `pkgs.buildGoModule` with vendor hash for reproducible builds. Dependencies are managed via Go modules.

**Development Shell:**
Includes Git, Nil (Nix LSP), and Go toolchain. Pre-commit hooks are automatically activated on shell entry.

**ccusage Integration:**
The application integrates with the ccusage CLI tool via npx:
- Default configuration uses `npx ccusage` for automatic package resolution
- Custom paths can be specified in configuration for alternative installations
- Supports JSON output parsing and caching for performance