package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/animation"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/display"
	"github.com/stretchr/testify/assert"
)

// TestApplicationIntegration tests the full application integration
func TestApplicationIntegration(t *testing.T) {
	ctx := context.Background()

	// Initialize configuration manager
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Validate configuration
	err = configManager.ValidateConfig()
	assert.NoError(t, err)

	// Initialize plugin registry
	registry := core.NewPluginRegistry(configManager)

	// Register plugins
	ccusagePlugin := datasource.NewCcusageCliPlugin()
	rainbowAnimationPlugin := animation.NewRainbowAnimationPlugin()
	rainbowDisplayPlugin := display.NewRainbowTUIPlugin()

	err = registry.RegisterDataSource(ccusagePlugin)
	assert.NoError(t, err)

	err = registry.RegisterAnimation(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.RegisterDisplay(rainbowDisplayPlugin)
	assert.NoError(t, err)

	// Initialize plugins
	err = registry.InitializePlugin(ccusagePlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowDisplayPlugin)
	assert.NoError(t, err)

	// Verify active plugins are available
	activeDataSource, err := registry.GetActiveDataSource()
	assert.NoError(t, err)
	assert.Equal(t, ccusagePlugin, activeDataSource)

	activeAnimation, err := registry.GetActiveAnimation()
	assert.NoError(t, err)
	assert.Equal(t, rainbowAnimationPlugin, activeAnimation)

	activeDisplay, err := registry.GetActiveDisplay()
	assert.NoError(t, err)
	assert.Equal(t, rainbowDisplayPlugin, activeDisplay)

	// Test end-to-end workflow with mock data
	t.Run("EndToEndWorkflow", func(t *testing.T) {
		// Create mock cost data
		mockCostData := &domain.CostData{
			TotalCost: 42.50,
			Currency:  "USD",
			Timestamp: time.Now(),
			ModelBreakdown: map[string]float64{
				"claude-3-opus":   25.00,
				"claude-3-sonnet": 17.50,
			},
		}

		// Generate animation frame
		animationConfig := configManager.GetAnimationConfig()
		assert.NotNil(t, animationConfig)

		frame, err := activeAnimation.GenerateFrame(ctx, "$42.50", 5, animationConfig)
		assert.NoError(t, err)
		assert.NotNil(t, frame)
		assert.Equal(t, "$42.50", frame.Text)
		assert.NotEmpty(t, frame.Colors)

		// Create display data
		displayConfig := configManager.GetDisplayConfig()
		assert.NotNil(t, displayConfig)

		displayData := &domain.DisplayData{
			Cost:        mockCostData,
			Animation:   frame,
			Config:      displayConfig,
			LastUpdated: time.Now(),
		}

		// Render display
		output, err := activeDisplay.Render(ctx, displayData)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		// Check that ASCII art is generated (contains ASCII block characters)
		assert.Contains(t, output, "█")
	})

	// Test different animation patterns
	t.Run("AnimationPatterns", func(t *testing.T) {
		patterns := []domain.AnimationPattern{
			domain.PatternRainbow,
			domain.PatternGradient,
			domain.PatternPulse,
			domain.PatternWave,
		}

		for _, pattern := range patterns {
			animationConfig := &domain.AnimationConfig{
				Speed:   100 * time.Millisecond,
				Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
				Enabled: true,
				Pattern: pattern,
			}

			frame, err := activeAnimation.GenerateFrame(ctx, "Test", 0, animationConfig)
			assert.NoError(t, err, "Pattern: %s", pattern)
			assert.NotNil(t, frame, "Pattern: %s", pattern)
			assert.Equal(t, "Test", frame.Text, "Pattern: %s", pattern)
		}
	})

	// Test different display configurations
	t.Run("DisplayConfigurations", func(t *testing.T) {
		configs := []struct {
			name   string
			width  int
			height int
		}{
			{"Large", 120, 40},
			{"Medium", 80, 24},
			{"Small", 60, 20},
			{"Minimal", 40, 10},
		}

		mockCostData := &domain.CostData{
			TotalCost: 15.75,
			Currency:  "USD",
			Timestamp: time.Now(),
		}

		for _, config := range configs {
			displayConfig := &domain.DisplayConfig{
				RefreshRate: 1 * time.Second,
				Size: domain.DisplaySize{
					Width:  config.width,
					Height: config.height,
				},
			}

			displayData := &domain.DisplayData{
				Cost:        mockCostData,
				Config:      displayConfig,
				LastUpdated: time.Now(),
			}

			output, err := activeDisplay.Render(ctx, displayData)
			assert.NoError(t, err, "Config: %s", config.name)

			if config.width > 40 && config.height > 10 || mockCostData != nil {
				assert.NotEmpty(t, output, "Config: %s", config.name)
			}
		}
	})

	// Test plugin shutdown
	t.Run("PluginShutdown", func(t *testing.T) {
		// All plugins should be enabled
		assert.True(t, ccusagePlugin.IsEnabled())
		assert.True(t, rainbowAnimationPlugin.IsEnabled())
		assert.True(t, rainbowDisplayPlugin.IsEnabled())

		// Shutdown all plugins
		err := registry.ShutdownAll()
		assert.NoError(t, err)

		// All plugins should be disabled
		assert.False(t, ccusagePlugin.IsEnabled())
		assert.False(t, rainbowAnimationPlugin.IsEnabled())
		assert.False(t, rainbowDisplayPlugin.IsEnabled())
	})
}

// TestConfigurationEdgeCases tests edge cases in configuration
func TestConfigurationEdgeCases(t *testing.T) {
	t.Run("InvalidConfigValues", func(t *testing.T) {
		configManager := core.NewConfigManager()
		err := configManager.LoadConfig("")
		assert.NoError(t, err)

		// UpdateConfig now only supports plugin datasource updates
		// So we manually set invalid values for testing validation
		config := configManager.GetConfig()
		config.Animation.Pattern = "invalid-pattern"
		config.Display.Width = -1
		config.Animation.Speed = -100 * time.Millisecond

		// Validation should fail
		err = configManager.ValidateConfig()
		assert.Error(t, err)
	})

	t.Run("DefaultConfiguration", func(t *testing.T) {
		configManager := core.NewConfigManager()

		// Config is now loaded with defaults in NewConfigManager
		config := configManager.GetConfig()
		assert.NotNil(t, config)

		displayConfig := configManager.GetDisplayConfig()
		assert.NotNil(t, displayConfig)

		animationConfig := configManager.GetAnimationConfig()
		assert.NotNil(t, animationConfig)
	})
}

// TestPluginInteraction tests plugin interactions
func TestPluginInteraction(t *testing.T) {
	ctx := context.Background()

	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)

	// Register only animation and display plugins (no data source)
	rainbowAnimationPlugin := animation.NewRainbowAnimationPlugin()
	rainbowDisplayPlugin := display.NewRainbowTUIPlugin()

	err = registry.RegisterAnimation(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.RegisterDisplay(rainbowDisplayPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowDisplayPlugin)
	assert.NoError(t, err)

	t.Run("MissingDataSource", func(t *testing.T) {
		// Should fail to get active data source
		_, err := registry.GetActiveDataSource()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("AnimationWithoutData", func(t *testing.T) {
		// Animation should still work without data source
		animationConfig := &domain.AnimationConfig{
			Speed:   100 * time.Millisecond,
			Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
			Enabled: true,
			Pattern: domain.PatternRainbow,
		}

		frame, err := rainbowAnimationPlugin.GenerateFrame(ctx, "No Data", 0, animationConfig)
		assert.NoError(t, err)
		assert.Equal(t, "No Data", frame.Text)
	})

	t.Run("DisplayWithoutCostData", func(t *testing.T) {
		// Display should handle missing cost data gracefully
		displayConfig := &domain.DisplayConfig{
			RefreshRate: 1 * time.Second,
			Size:        domain.DisplaySize{Width: 40, Height: 10},
		}

		displayData := &domain.DisplayData{
			Cost:        nil, // No cost data
			Config:      displayConfig,
			LastUpdated: time.Now(),
		}

		output, err := rainbowDisplayPlugin.Render(ctx, displayData)
		assert.NoError(t, err)
		// Small display with no cost data should return empty string
		assert.Empty(t, output)
	})
}

// TestErrorHandling tests error handling scenarios
func TestErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("PluginNotEnabled", func(t *testing.T) {
		// Test animation plugin not enabled
		animationPlugin := animation.NewRainbowAnimationPlugin()
		assert.False(t, animationPlugin.IsEnabled())

		config := &domain.AnimationConfig{
			Speed:   100 * time.Millisecond,
			Colors:  []string{"#FF0000"},
			Enabled: true,
			Pattern: domain.PatternRainbow,
		}

		_, err := animationPlugin.GenerateFrame(ctx, "test", 0, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not enabled")

		// Test display plugin not enabled
		displayPlugin := display.NewRainbowTUIPlugin()
		assert.False(t, displayPlugin.IsEnabled())

		displayData := &domain.DisplayData{
			Config: &domain.DisplayConfig{
				RefreshRate: 1 * time.Second,
				Size:        domain.DisplaySize{Width: 40, Height: 10},
			},
		}

		_, err = displayPlugin.Render(ctx, displayData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not enabled")
	})

	t.Run("InvalidPluginConfiguration", func(t *testing.T) {
		animationPlugin := animation.NewRainbowAnimationPlugin()
		err := animationPlugin.Initialize(map[string]interface{}{})
		assert.NoError(t, err)

		// Test with invalid config
		invalidConfigs := []*domain.AnimationConfig{
			nil, // Nil config
			{
				Speed:   0, // Invalid speed
				Colors:  []string{"#FF0000"},
				Enabled: true,
				Pattern: domain.PatternRainbow,
			},
			{
				Speed:   100 * time.Millisecond,
				Colors:  []string{}, // Empty colors
				Enabled: true,
				Pattern: domain.PatternRainbow,
			},
			{
				Speed:   100 * time.Millisecond,
				Colors:  []string{"invalid-color"}, // Invalid color format
				Enabled: true,
				Pattern: domain.PatternRainbow,
			},
		}

		for i, config := range invalidConfigs {
			err := animationPlugin.ValidateAnimationConfig(config)
			assert.Error(t, err, "Config %d should be invalid", i)
		}
	})
}

// TestBankruptcyDataSourceIntegration tests bankruptcy data source functionality
func TestBankruptcyDataSourceIntegration(t *testing.T) {
	ctx := context.Background()

	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Update configuration to use bankruptcy data source
	err = configManager.UpdateConfig(map[string]interface{}{
		"plugins.datasource": "bankruptcy-datasource",
	})
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)

	// Register and initialize plugins including bankruptcy data source
	bankruptcyPlugin := datasource.NewBankruptcyDataSourcePlugin()
	rainbowAnimationPlugin := animation.NewRainbowAnimationPlugin()
	rainbowDisplayPlugin := display.NewRainbowTUIPlugin()

	err = registry.RegisterDataSource(bankruptcyPlugin)
	assert.NoError(t, err)

	err = registry.RegisterAnimation(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.RegisterDisplay(rainbowDisplayPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(bankruptcyPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowDisplayPlugin)
	assert.NoError(t, err)

	t.Run("BankruptcyDataSourceEndToEnd", func(t *testing.T) {
		// Verify active data source is bankruptcy plugin
		activeDataSource, err := registry.GetActiveDataSource()
		assert.NoError(t, err)
		assert.Equal(t, "bankruptcy-datasource", activeDataSource.Name())

		// Fetch cost data from bankruptcy source
		costData, err := activeDataSource.FetchCostData(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, costData)
		assert.Equal(t, 9999.99, costData.TotalCost)

		// Generate animation frame for bankruptcy amount
		animationConfig := configManager.GetAnimationConfig()
		assert.NotNil(t, animationConfig)

		frame, err := rainbowAnimationPlugin.GenerateFrame(ctx, "$9999.99", 5, animationConfig)
		assert.NoError(t, err)
		assert.NotNil(t, frame)
		assert.Equal(t, "$9999.99", frame.Text)
		assert.NotEmpty(t, frame.Colors)

		// Create display data for bankruptcy mode
		displayConfig := configManager.GetDisplayConfig()
		assert.NotNil(t, displayConfig)

		displayData := &domain.DisplayData{
			Cost:        costData,
			Animation:   frame,
			Config:      displayConfig,
			LastUpdated: time.Now(),
		}

		// Render display
		output, err := rainbowDisplayPlugin.Render(ctx, displayData)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		// Check that ASCII art is generated for bankruptcy amount
		assert.Contains(t, output, "█")
	})
}
