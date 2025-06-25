package display_test

import (
	"context"
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/display"
	"github.com/stretchr/testify/assert"
)

func TestNewRainbowTUIPlugin(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	assert.NotNil(t, plugin)
	assert.Equal(t, "rainbow-display", plugin.Name())
	assert.Equal(t, "1.0.0", plugin.Version())
	assert.Equal(t, "Rainbow TUI display plugin", plugin.Description())
	assert.False(t, plugin.IsEnabled()) // Should be disabled initially
}

func TestRainbowTUIPlugin_Initialize(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()

	// Test with empty config
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())
}

func TestRainbowTUIPlugin_Shutdown(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()

	// Initialize first
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())

	// Shutdown
	err = plugin.Shutdown()
	assert.NoError(t, err)
	assert.False(t, plugin.IsEnabled())
}

func TestRainbowTUIPlugin_GetCapabilities(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()

	capabilities := plugin.GetCapabilities()
	assert.Equal(t, 200, capabilities.MaxWidth)
	assert.Equal(t, 50, capabilities.MaxHeight)
	assert.True(t, capabilities.SupportsColor)
	assert.True(t, capabilities.SupportsUnicode)
}

func TestRainbowTUIPlugin_ValidateDisplayConfig(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()

	// Test valid config
	validConfig := &domain.DisplayConfig{
		RefreshRate: 1 * time.Second,
		Size: domain.DisplaySize{
			Width:  80,
			Height: 24,
		},
	}

	err := plugin.ValidateDisplayConfig(validConfig)
	assert.NoError(t, err)

	// Test nil config
	err = plugin.ValidateDisplayConfig(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")

	// Test width exceeds maximum
	invalidConfig := &domain.DisplayConfig{
		Size: domain.DisplaySize{Width: 300, Height: 24}, // Exceeds max of 200
	}

	err = plugin.ValidateDisplayConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "width")
	assert.Contains(t, err.Error(), "exceeds maximum")

	// Test height exceeds maximum
	invalidConfig.Size.Width = 80
	invalidConfig.Size.Height = 60 // Exceeds max of 50

	err = plugin.ValidateDisplayConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "height")
	assert.Contains(t, err.Error(), "exceeds maximum")
}

func TestRainbowTUIPlugin_Render_NotEnabled(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 25.75,
			Currency:  "USD",
			Timestamp: time.Now(),
		},
		Config: &domain.DisplayConfig{
			Size: domain.DisplaySize{Width: 80, Height: 24},
		},
	}

	// Should fail when plugin is not enabled
	_, err := plugin.Render(ctx, displayData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin is not enabled")
}

func TestRainbowTUIPlugin_Render_NilData(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	// Should fail when data is nil
	_, err = plugin.Render(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "display data cannot be nil")
}

func TestRainbowTUIPlugin_Render_WithCostData(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	now := time.Now()
	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 25.75,
			Currency:  "USD",
			Timestamp: now,
		},
		Animation: &domain.AnimationFrame{
			Colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			Text:      "$25.75",
			Timestamp: now,
		},
		Config: &domain.DisplayConfig{
			RefreshRate: 1 * time.Second,
			Size: domain.DisplaySize{
				Width:  80,
				Height: 24,
			},
		},
		LastUpdated: now,
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	// Check for ASCII art block characters
	assert.Contains(t, output, "█") // Should contain ASCII block characters
}

func TestRainbowTUIPlugin_Render_NoCostData(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	displayData := &domain.DisplayData{
		Cost: nil, // No cost data
		Config: &domain.DisplayConfig{
			Size: domain.DisplaySize{
				Width:  40,
				Height: 10,
			},
		},
		LastUpdated: time.Now(),
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.Empty(t, output) // Should be empty when no cost data
}

func TestRainbowTUIPlugin_Render_NoAnimation(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	now := time.Now()
	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 25.75,
			Currency:  "USD",
			Timestamp: now,
		},
		Animation: nil, // No animation
		Config: &domain.DisplayConfig{
			Size: domain.DisplaySize{
				Width:  20,
				Height: 5,
			},
		},
		LastUpdated: now,
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	// Check for ASCII art block characters
	assert.Contains(t, output, "█") // Should contain ASCII block characters
}

func TestRainbowTUIPlugin_Render_SmallDisplay(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	now := time.Now()
	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 99.99,
			Currency:  "USD",
			Timestamp: now,
		},
		Config: &domain.DisplayConfig{
			Size: domain.DisplaySize{
				Width:  30, // Small width should trigger small patterns
				Height: 8,
			},
		},
		LastUpdated: now,
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	// Small displays should still generate ASCII art
	assert.Contains(t, output, "█")
}

func TestRainbowTUIPlugin_Render_LargeDisplay(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	now := time.Now()
	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 1234.56,
			Currency:  "USD",
			Timestamp: now,
		},
		Config: &domain.DisplayConfig{
			Size: domain.DisplaySize{
				Width:  120, // Large width should trigger large patterns
				Height: 30,
			},
		},
		LastUpdated: now,
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	// Large displays should generate more detailed ASCII art
	assert.Contains(t, output, "█")
	// Large display output should be longer than small display output
	assert.True(t, len(output) > 100, "Large display should generate substantial output")
}
