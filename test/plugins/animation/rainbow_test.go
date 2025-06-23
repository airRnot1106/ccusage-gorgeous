package animation_test

import (
	"context"
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/animation"
	"github.com/stretchr/testify/assert"
)

func TestNewRainbowAnimationPlugin(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	assert.NotNil(t, plugin)
	assert.Equal(t, "rainbow-animation", plugin.Name())
	assert.Equal(t, "1.0.0", plugin.Version())
	assert.Equal(t, "Rainbow animation effects plugin", plugin.Description())
	assert.False(t, plugin.IsEnabled()) // Should be disabled initially
}

func TestRainbowAnimationPlugin_Initialize(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()

	// Test with empty config
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())
}

func TestRainbowAnimationPlugin_Shutdown(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()

	// Initialize first
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())

	// Shutdown
	err = plugin.Shutdown()
	assert.NoError(t, err)
	assert.False(t, plugin.IsEnabled())
}

func TestRainbowAnimationPlugin_GetSupportedPatterns(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()

	patterns := plugin.GetSupportedPatterns()
	assert.Len(t, patterns, 4)
	assert.Contains(t, patterns, domain.PatternRainbow)
	assert.Contains(t, patterns, domain.PatternGradient)
	assert.Contains(t, patterns, domain.PatternPulse)
	assert.Contains(t, patterns, domain.PatternWave)
}

func TestRainbowAnimationPlugin_ValidateAnimationConfig(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()

	// Test valid config
	validConfig := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	err := plugin.ValidateAnimationConfig(validConfig)
	assert.NoError(t, err)

	// Test nil config
	err = plugin.ValidateAnimationConfig(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")

	// Test zero speed
	invalidConfig := &domain.AnimationConfig{
		Speed:   0,
		Colors:  []string{"#FF0000"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	err = plugin.ValidateAnimationConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "speed must be positive")

	// Test empty colors
	invalidConfig.Speed = 100 * time.Millisecond
	invalidConfig.Colors = []string{}

	err = plugin.ValidateAnimationConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one color")

	// Test invalid color format
	invalidConfig.Colors = []string{"invalid-color"}

	err = plugin.ValidateAnimationConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid color format")

	// Test unsupported pattern
	invalidConfig.Colors = []string{"#FF0000"}
	invalidConfig.Pattern = domain.AnimationPattern("unsupported")

	err = plugin.ValidateAnimationConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported animation pattern")
}

func TestRainbowAnimationPlugin_GenerateFrame_NotEnabled(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	// Should fail when plugin is not enabled
	_, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin is not enabled")
}

func TestRainbowAnimationPlugin_GenerateFrame_NilConfig(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	// Should fail when config is nil
	_, err = plugin.GenerateFrame(ctx, "test", 0, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animation config is required")
}

func TestRainbowAnimationPlugin_GenerateFrame_DisabledAnimation(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: false, // Animation disabled
		Pattern: domain.PatternRainbow,
	}

	frame, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "test", frame.Text)
	assert.Len(t, frame.Colors, 1)
	assert.Equal(t, "#FFFFFF", frame.Colors[0]) // Should be white
}

func TestRainbowAnimationPlugin_GenerateFrame_RainbowPattern(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	frame, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "test", frame.Text)
	assert.Len(t, frame.Colors, 4) // Length of "test"
	assert.False(t, frame.Timestamp.IsZero())

	// Test different frame number produces different colors
	frame2, err := plugin.GenerateFrame(ctx, "test", 1, config)
	assert.NoError(t, err)
	assert.NotEqual(t, frame.Colors, frame2.Colors)
}

func TestRainbowAnimationPlugin_GenerateFrame_GradientPattern(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternGradient,
	}

	frame, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "test", frame.Text)
	assert.Len(t, frame.Colors, 4)
}

func TestRainbowAnimationPlugin_GenerateFrame_PulsePattern(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00"},
		Enabled: true,
		Pattern: domain.PatternPulse,
	}

	frame, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "test", frame.Text)
	assert.Len(t, frame.Colors, 4)

	// All colors should be the same for pulse pattern
	firstColor := frame.Colors[0]
	for _, color := range frame.Colors {
		assert.Equal(t, firstColor, color)
	}
}

func TestRainbowAnimationPlugin_GenerateFrame_WavePattern(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternWave,
	}

	frame, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "test", frame.Text)
	assert.Len(t, frame.Colors, 4)
}

func TestRainbowAnimationPlugin_GenerateFrame_EmptyText(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	frame, err := plugin.GenerateFrame(ctx, "", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "", frame.Text)
	assert.Len(t, frame.Colors, 0)
}

func TestRainbowAnimationPlugin_GenerateFrame_SingleColor(t *testing.T) {
	plugin := animation.NewRainbowAnimationPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000"}, // Single color
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	frame, err := plugin.GenerateFrame(ctx, "test", 0, config)
	assert.NoError(t, err)
	assert.Equal(t, "test", frame.Text)
	assert.Len(t, frame.Colors, 4)

	// All colors should be the same
	for _, color := range frame.Colors {
		assert.Equal(t, "#FF0000", color)
	}
}
