package animation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// RainbowAnimationPlugin implements rainbow animation effects
type RainbowAnimationPlugin struct {
	name        string
	version     string
	description string
	enabled     bool
	frameCount  int
}

// NewRainbowAnimationPlugin creates a new rainbow animation plugin
func NewRainbowAnimationPlugin() *RainbowAnimationPlugin {
	return &RainbowAnimationPlugin{
		name:        "rainbow-animation",
		version:     "1.0.0",
		description: "Rainbow animation effects plugin",
		enabled:     false,
		frameCount:  0,
	}
}

// Name returns the plugin name
func (r *RainbowAnimationPlugin) Name() string {
	return r.name
}

// Version returns the plugin version
func (r *RainbowAnimationPlugin) Version() string {
	return r.version
}

// Description returns the plugin description
func (r *RainbowAnimationPlugin) Description() string {
	return r.description
}

// IsEnabled returns whether the plugin is enabled
func (r *RainbowAnimationPlugin) IsEnabled() bool {
	return r.enabled
}

// Initialize initializes the plugin with configuration
func (r *RainbowAnimationPlugin) Initialize(config map[string]interface{}) error {
	r.enabled = true
	r.frameCount = 0
	return nil
}

// Shutdown shuts down the plugin
func (r *RainbowAnimationPlugin) Shutdown() error {
	r.enabled = false
	return nil
}

// GenerateFrame generates an animation frame with rainbow colors
func (r *RainbowAnimationPlugin) GenerateFrame(ctx context.Context, text string, frameNumber int, config *domain.AnimationConfig) (*domain.AnimationFrame, error) {
	if !r.enabled {
		return nil, fmt.Errorf("plugin is not enabled")
	}

	if config == nil {
		return nil, fmt.Errorf("animation config is required")
	}

	if !config.Enabled {
		// Return static frame when animation is disabled
		return &domain.AnimationFrame{
			Colors:    []string{"#FFFFFF"}, // White
			Text:      text,
			Timestamp: time.Now(),
		}, nil
	}

	var colors []string

	switch config.Pattern {
	case domain.PatternRainbow:
		colors = r.generateRainbowColors(frameNumber, len(text), config.Colors)
	case domain.PatternGradient:
		colors = r.generateGradientColors(frameNumber, len(text), config.Colors)
	case domain.PatternPulse:
		colors = r.generatePulseColors(frameNumber, len(text), config.Colors)
	case domain.PatternWave:
		colors = r.generateWaveColors(frameNumber, len(text), config.Colors)
	default:
		colors = r.generateRainbowColors(frameNumber, len(text), config.Colors)
	}

	frame := &domain.AnimationFrame{
		Colors:    colors,
		Text:      text,
		Timestamp: time.Now(),
	}

	r.frameCount++
	return frame, nil
}

// GetSupportedPatterns returns the animation patterns supported by this plugin
func (r *RainbowAnimationPlugin) GetSupportedPatterns() []domain.AnimationPattern {
	return []domain.AnimationPattern{
		domain.PatternRainbow,
		domain.PatternGradient,
		domain.PatternPulse,
		domain.PatternWave,
	}
}

// ValidateAnimationConfig validates the animation configuration
func (r *RainbowAnimationPlugin) ValidateAnimationConfig(config *domain.AnimationConfig) error {
	if config == nil {
		return fmt.Errorf("animation config cannot be nil")
	}

	if config.Speed <= 0 {
		return fmt.Errorf("animation speed must be positive")
	}

	if len(config.Colors) == 0 {
		return fmt.Errorf("at least one color must be specified")
	}

	// Validate color format (basic hex color validation)
	for i, color := range config.Colors {
		if len(color) != 7 || color[0] != '#' {
			return fmt.Errorf("invalid color format at index %d: %s", i, color)
		}
	}

	// Check if pattern is supported
	supported := false
	for _, pattern := range r.GetSupportedPatterns() {
		if config.Pattern == pattern {
			supported = true
			break
		}
	}
	if !supported {
		return fmt.Errorf("unsupported animation pattern: %s", config.Pattern)
	}

	return nil
}

// generateRainbowColors generates rainbow-shifting colors
func (r *RainbowAnimationPlugin) generateRainbowColors(frameNumber, textLength int, baseColors []string) []string {
	if len(baseColors) == 0 {
		return []string{"#FFFFFF"}
	}

	colors := make([]string, textLength)
	for i := 0; i < textLength; i++ {
		colorIndex := (frameNumber + i) % len(baseColors)
		colors[i] = baseColors[colorIndex]
	}
	return colors
}

// generateGradientColors generates smooth gradient colors
func (r *RainbowAnimationPlugin) generateGradientColors(frameNumber, textLength int, baseColors []string) []string {
	if len(baseColors) == 0 {
		return []string{"#FFFFFF"}
	}

	if textLength == 1 {
		return []string{baseColors[frameNumber%len(baseColors)]}
	}

	colors := make([]string, textLength)
	for i := 0; i < textLength; i++ {
		// Create gradient across the text length
		progress := float64(i) / float64(textLength-1)
		gradientPos := progress + float64(frameNumber)*0.01
		colorIndex := int(gradientPos*float64(len(baseColors))) % len(baseColors)
		colors[i] = baseColors[colorIndex]
	}
	return colors
}

// generatePulseColors generates pulsing colors
func (r *RainbowAnimationPlugin) generatePulseColors(frameNumber, textLength int, baseColors []string) []string {
	if len(baseColors) == 0 {
		return []string{"#FFFFFF"}
	}

	// Pulse between first and second color (or just first if only one)
	pulseValue := math.Sin(float64(frameNumber) * 0.2)
	var currentColor string

	if len(baseColors) >= 2 {
		if pulseValue > 0 {
			currentColor = baseColors[0]
		} else {
			currentColor = baseColors[1]
		}
	} else {
		currentColor = baseColors[0]
	}

	colors := make([]string, textLength)
	for i := 0; i < textLength; i++ {
		colors[i] = currentColor
	}
	return colors
}

// generateWaveColors generates wave-like color patterns
func (r *RainbowAnimationPlugin) generateWaveColors(frameNumber, textLength int, baseColors []string) []string {
	if len(baseColors) == 0 {
		return []string{"#FFFFFF"}
	}

	colors := make([]string, textLength)
	for i := 0; i < textLength; i++ {
		// Create wave pattern with sine function
		waveValue := math.Sin(float64(frameNumber)*0.1 + float64(i)*0.5)
		colorIndex := int((waveValue+1)/2*float64(len(baseColors))) % len(baseColors)
		colors[i] = baseColors[colorIndex]
	}
	return colors
}
