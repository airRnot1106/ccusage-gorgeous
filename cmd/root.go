package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/airRnot1106/ccusage-gorgeous/internal/infrastructure/tui"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/animation"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/display"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// Root cobra command
var rootCmd = &cobra.Command{
	Use:   "ccugorg",
	Short: "TUI application for displaying Claude API costs with rainbow animations",
	Long: `ccugorg is a terminal user interface application that displays
Claude API usage costs with beautiful rainbow animations and ASCII art.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runApplication,
}

// Flag variables
var (
	animationSpeed   string
	animationPattern string
	noAnimation      bool
	bankruptcy       bool
)

func init() {
	// Add flags
	rootCmd.Flags().StringVar(&animationSpeed, "animation-speed", "", "Animation speed (e.g., 100ms)")
	rootCmd.Flags().StringVar(&animationPattern, "animation-pattern", "", "Animation pattern (rainbow, gradient, pulse, wave)")
	rootCmd.Flags().BoolVar(&noAnimation, "no-animation", false, "Disable animation")

	// Hidden bankruptcy flag
	rootCmd.Flags().BoolVar(&bankruptcy, "bankruptcy", false, "")
	_ = rootCmd.Flags().MarkHidden("bankruptcy") // Hide bankruptcy flag from help
}

// runApplication executes the main application logic
func runApplication(cmd *cobra.Command, args []string) error {
	// Create context
	ctx := context.Background()

	// Convert cobra flags to our flag config structure
	flagConfig, err := convertCobraFlags()
	if err != nil {
		return fmt.Errorf("failed to convert flags: %w", err)
	}

	// Initialize configuration manager
	configManager := core.NewConfigManager()
	if err := configManager.LoadConfig(""); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Apply command line flags to override configuration
	if err := configManager.ApplyFlagsToConfig(flagConfig); err != nil {
		return fmt.Errorf("failed to apply command line flags: %w", err)
	}

	// Validate configuration
	if err := configManager.ValidateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Update configuration for bankruptcy mode
	if bankruptcy {
		if err := configManager.UpdateConfig(map[string]interface{}{
			"plugins.datasource": "bankruptcy-datasource",
		}); err != nil {
			return fmt.Errorf("failed to update config for bankruptcy mode: %w", err)
		}
	}

	// Initialize plugin registry
	registry := core.NewPluginRegistry(configManager)

	// Register built-in plugins
	if err := registerPlugins(registry, bankruptcy); err != nil {
		return fmt.Errorf("failed to register plugins: %w", err)
	}

	// Initialize plugins
	if err := initializePlugins(registry); err != nil {
		return fmt.Errorf("failed to initialize plugins: %w", err)
	}

	// Verify required plugins are available
	if err := verifyRequiredPlugins(registry); err != nil {
		return fmt.Errorf("required plugins not available: %w", err)
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
		return fmt.Errorf("error running TUI program: %w", err)
	}

	return nil
}

// convertCobraFlags converts cobra flag variables to FlagConfig structure
func convertCobraFlags() (*core.FlagConfig, error) {
	flagConfig := &core.FlagConfig{}

	// Parse animation speed
	if animationSpeed != "" {
		speed, err := time.ParseDuration(animationSpeed)
		if err != nil {
			return nil, fmt.Errorf("invalid animation speed format '%s': %w", animationSpeed, err)
		}
		flagConfig.Animation.Speed = speed
	}

	// Parse animation pattern
	if animationPattern != "" {
		pattern := domain.AnimationPattern(animationPattern)
		// Validate pattern
		validPatterns := []domain.AnimationPattern{
			domain.PatternRainbow, domain.PatternGradient,
			domain.PatternPulse, domain.PatternWave,
		}
		isValid := false
		for _, validPattern := range validPatterns {
			if pattern == validPattern {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, fmt.Errorf("invalid animation pattern '%s'. Valid patterns: rainbow, gradient, pulse, wave", animationPattern)
		}
		flagConfig.Animation.Pattern = pattern
	}

	// Parse no-animation flag
	if noAnimation {
		enabled := false
		flagConfig.Animation.Enabled = &enabled
	}

	// Parse bankruptcy flag
	flagConfig.Bankruptcy = bankruptcy

	return flagConfig, nil
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

	return nil
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
