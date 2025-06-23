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
	assert.Len(t, capabilities.SupportedFormats, 4)
	assert.Contains(t, capabilities.SupportedFormats, domain.FormatLarge)
	assert.Contains(t, capabilities.SupportedFormats, domain.FormatMedium)
	assert.Contains(t, capabilities.SupportedFormats, domain.FormatSmall)
	assert.Contains(t, capabilities.SupportedFormats, domain.FormatMinimal)
	assert.Equal(t, 200, capabilities.MaxWidth)
	assert.Equal(t, 50, capabilities.MaxHeight)
	assert.True(t, capabilities.SupportsColor)
	assert.True(t, capabilities.SupportsUnicode)
}

func TestRainbowTUIPlugin_ValidateDisplayConfig(t *testing.T) {
	plugin := display.NewRainbowTUIPlugin()

	// Test valid config
	validConfig := &domain.DisplayConfig{
		RefreshRate:   1 * time.Second,
		ShowTimestamp: true,
		ShowBreakdown: true,
		Format:        domain.FormatLarge,
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

	// Test unsupported format
	invalidConfig := &domain.DisplayConfig{
		Format: domain.DisplayFormat("unsupported"),
		Size:   domain.DisplaySize{Width: 80, Height: 24},
	}

	err = plugin.ValidateDisplayConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported display format")

	// Test width exceeds maximum
	invalidConfig.Format = domain.FormatLarge
	invalidConfig.Size.Width = 300 // Exceeds max of 200

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
			Format: domain.FormatLarge,
			Size:   domain.DisplaySize{Width: 80, Height: 24},
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

func TestRainbowTUIPlugin_Render_LargeFormat(t *testing.T) {
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
			ModelBreakdown: map[string]float64{
				"claude-3-opus":   15.25,
				"claude-3-sonnet": 10.50,
			},
		},
		Animation: &domain.AnimationFrame{
			Colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			Text:      "$25.75",
			Timestamp: now,
		},
		Config: &domain.DisplayConfig{
			RefreshRate:   1 * time.Second,
			ShowTimestamp: true,
			ShowBreakdown: true,
			Format:        domain.FormatLarge,
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
	// Check for ASCII art block characters instead of literal cost values
	assert.Contains(t, output, "█") // Should contain ASCII block characters
	// Note: Model breakdown display is not implemented in current ASCII art rendering
}

func TestRainbowTUIPlugin_Render_MediumFormat(t *testing.T) {
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
			RefreshRate:   1 * time.Second,
			ShowTimestamp: true,
			ShowBreakdown: false,
			Format:        domain.FormatMedium,
			Size: domain.DisplaySize{
				Width:  60,
				Height: 20,
			},
		},
		LastUpdated: now,
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	// Check for ASCII art block characters instead of literal cost values
	assert.Contains(t, output, "█") // Should contain ASCII block characters
}

func TestRainbowTUIPlugin_Render_SmallFormat(t *testing.T) {
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
			Format: domain.FormatSmall,
			Size: domain.DisplaySize{
				Width:  40,
				Height: 10,
			},
		},
		LastUpdated: now,
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	// Check for ASCII art block characters instead of literal cost values
	assert.Contains(t, output, "█") // Should contain ASCII block characters
}

func TestRainbowTUIPlugin_Render_MinimalFormat(t *testing.T) {
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
			Format: domain.FormatMinimal,
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
	// Check for ASCII art block characters instead of literal cost values
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
			Format: domain.FormatSmall,
			Size: domain.DisplaySize{
				Width:  40,
				Height: 10,
			},
		},
		LastUpdated: time.Now(),
	}

	output, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.Empty(t, output) // Should be empty for small format with no cost data
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
			Format: domain.FormatMinimal,
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
	// Check for ASCII art block characters instead of literal cost values
	assert.Contains(t, output, "█") // Should contain ASCII block characters
}

func TestRainbowTUIPlugin_Render_NoTimestamp(t *testing.T) {
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
			RefreshRate:   1 * time.Second,
			ShowTimestamp: false, // Don't show timestamp
			ShowBreakdown: false,
			Format:        domain.FormatLarge,
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
	// Check for ASCII art block characters instead of literal cost values
	assert.Contains(t, output, "█") // Should contain ASCII block characters
	// Should not contain timestamp info when disabled
	assert.NotContains(t, output, "Last Updated")
}
