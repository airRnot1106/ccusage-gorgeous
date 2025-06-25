package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/infrastructure/tui"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/animation"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/display"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags using our custom parser
	// Check for help flag first
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			printUsage()
			return
		}
	}

	flagConfig, err := core.ParseFlagsFromArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to parse command line flags: %v", err)
	}

	// Get bankruptcy mode from parsed flags
	bankruptcyMode := flagConfig.Bankruptcy

	// Create context
	ctx := context.Background()

	// Initialize configuration manager
	configManager := core.NewConfigManager()
	if err := configManager.LoadConfig(""); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Apply command line flags to override configuration
	if err := configManager.ApplyFlagsToConfig(flagConfig); err != nil {
		log.Fatalf("Failed to apply command line flags: %v", err)
	}

	// Validate configuration
	if err := configManager.ValidateConfig(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Update configuration for bankruptcy mode
	if bankruptcyMode {
		if err := configManager.UpdateConfig(map[string]interface{}{
			"plugins.datasource": "bankruptcy-datasource",
		}); err != nil {
			log.Fatalf("Failed to update config for bankruptcy mode: %v", err)
		}
	}

	// Initialize plugin registry
	registry := core.NewPluginRegistry(configManager)

	// Register built-in plugins
	if err := registerPlugins(registry, bankruptcyMode); err != nil {
		log.Fatalf("Failed to register plugins: %v", err)
	}

	// Initialize plugins
	if err := initializePlugins(registry); err != nil {
		log.Fatalf("Failed to initialize plugins: %v", err)
	}

	// Verify required plugins are available
	if err := verifyRequiredPlugins(registry); err != nil {
		log.Fatalf("Required plugins not available: %v", err)
	}

	// Create TUI model
	model := tui.NewModel(ctx, registry, configManager)

	// Create TUI program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Setup cleanup
	defer func() {
		if err := registry.ShutdownAll(); err != nil {
			log.Printf("Warning: Error during plugin shutdown: %v", err)
		}
	}()

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI program: %v", err)
	}
}

// registerPlugins registers all built-in plugins
func registerPlugins(registry *core.PluginRegistry, bankruptcyMode bool) error {
	// Register appropriate data source plugin based on bankruptcy mode
	if bankruptcyMode {
		bankruptcyPlugin := datasource.NewBankruptcyDataSourcePlugin()
		if err := registry.RegisterDataSource(bankruptcyPlugin); err != nil {
			return fmt.Errorf("failed to register bankruptcy data source plugin: %w", err)
		}
	} else {
		ccusagePlugin := datasource.NewCcusageCliPlugin()
		if err := registry.RegisterDataSource(ccusagePlugin); err != nil {
			return fmt.Errorf("failed to register ccusage CLI plugin: %w", err)
		}
	}

	// Register animation plugins
	rainbowAnimationPlugin := animation.NewRainbowAnimationPlugin()
	if err := registry.RegisterAnimation(rainbowAnimationPlugin); err != nil {
		return fmt.Errorf("failed to register rainbow animation plugin: %w", err)
	}

	// Register display plugins
	rainbowDisplayPlugin := display.NewRainbowTUIPlugin()
	if err := registry.RegisterDisplay(rainbowDisplayPlugin); err != nil {
		return fmt.Errorf("failed to register rainbow display plugin: %w", err)
	}

	return nil
}

// initializePlugins initializes all registered plugins
func initializePlugins(registry *core.PluginRegistry) error {
	plugins := registry.ListPlugins()

	for _, plugin := range plugins {
		if err := registry.InitializePlugin(plugin); err != nil {
			return fmt.Errorf("failed to initialize plugin '%s': %w", plugin.Name(), err)
		}

		// Plugin initialized silently
	}

	return nil
}

// verifyRequiredPlugins verifies that all required plugins are available
func verifyRequiredPlugins(registry *core.PluginRegistry) error {
	// Check data source plugin
	if _, err := registry.GetActiveDataSource(); err != nil {
		return fmt.Errorf("active data source plugin not available: %w", err)
	}

	// Check display plugin
	if _, err := registry.GetActiveDisplay(); err != nil {
		return fmt.Errorf("active display plugin not available: %w", err)
	}

	// Check animation plugin
	if _, err := registry.GetActiveAnimation(); err != nil {
		return fmt.Errorf("active animation plugin not available: %w", err)
	}

	// Plugin verification completed silently

	return nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", "ccugorg")
	fmt.Fprintf(os.Stderr, "  --animation-speed string\n")
	fmt.Fprintf(os.Stderr, "        Animation speed (e.g., 100ms)\n")
	fmt.Fprintf(os.Stderr, "  --animation-pattern string\n")
	fmt.Fprintf(os.Stderr, "        Animation pattern (rainbow, gradient, pulse, wave)\n")
	fmt.Fprintf(os.Stderr, "  --no-animation\n")
	fmt.Fprintf(os.Stderr, "        Disable animation\n")
}
