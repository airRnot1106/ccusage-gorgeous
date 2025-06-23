package display

import (
	"context"
	"fmt"
	"strings"

	"github.com/airRnot1106/ccusage-gorgeous/internal/application/interfaces"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

// RainbowTUIPlugin implements the DisplayPlugin interface for rainbow TUI display
type RainbowTUIPlugin struct {
	name        string
	version     string
	description string
	enabled     bool
}

// NewRainbowTUIPlugin creates a new rainbow TUI display plugin
func NewRainbowTUIPlugin() *RainbowTUIPlugin {
	return &RainbowTUIPlugin{
		name:        "rainbow-display",
		version:     "1.0.0",
		description: "Rainbow TUI display plugin",
		enabled:     false,
	}
}

// Name returns the plugin name
func (r *RainbowTUIPlugin) Name() string {
	return r.name
}

// Version returns the plugin version
func (r *RainbowTUIPlugin) Version() string {
	return r.version
}

// Description returns the plugin description
func (r *RainbowTUIPlugin) Description() string {
	return r.description
}

// IsEnabled returns whether the plugin is enabled
func (r *RainbowTUIPlugin) IsEnabled() bool {
	return r.enabled
}

// Initialize initializes the plugin with configuration
func (r *RainbowTUIPlugin) Initialize(config map[string]interface{}) error {
	r.enabled = true
	return nil
}

// Shutdown shuts down the plugin
func (r *RainbowTUIPlugin) Shutdown() error {
	r.enabled = false
	return nil
}

// Render renders the display data with rainbow animation
func (r *RainbowTUIPlugin) Render(ctx context.Context, data *domain.DisplayData) (string, error) {
	if !r.enabled {
		return "", fmt.Errorf("plugin is not enabled")
	}

	if data == nil {
		return "", fmt.Errorf("display data cannot be nil")
	}

	// Render based on format
	switch data.Config.Format {
	case domain.FormatLarge:
		return r.renderLarge(data)
	case domain.FormatMedium:
		return r.renderMedium(data)
	case domain.FormatSmall:
		return r.renderSmall(data)
	case domain.FormatMinimal:
		return r.renderMinimal(data)
	default:
		return r.renderLarge(data)
	}
}

// GetCapabilities returns the display capabilities
func (r *RainbowTUIPlugin) GetCapabilities() interfaces.DisplayCapabilities {
	return interfaces.DisplayCapabilities{
		SupportedFormats: []domain.DisplayFormat{
			domain.FormatLarge,
			domain.FormatMedium,
			domain.FormatSmall,
			domain.FormatMinimal,
		},
		MaxWidth:        200,
		MaxHeight:       50,
		SupportsColor:   true,
		SupportsUnicode: true,
	}
}

// ValidateDisplayConfig validates the display configuration
func (r *RainbowTUIPlugin) ValidateDisplayConfig(config *domain.DisplayConfig) error {
	if config == nil {
		return fmt.Errorf("display config cannot be nil")
	}

	capabilities := r.GetCapabilities()

	// Check if format is supported
	formatSupported := false
	for _, format := range capabilities.SupportedFormats {
		if config.Format == format {
			formatSupported = true
			break
		}
	}
	if !formatSupported {
		return fmt.Errorf("unsupported display format: %s", config.Format)
	}

	// Check dimensions
	if config.Size.Width > capabilities.MaxWidth {
		return fmt.Errorf("width %d exceeds maximum %d", config.Size.Width, capabilities.MaxWidth)
	}
	if config.Size.Height > capabilities.MaxHeight {
		return fmt.Errorf("height %d exceeds maximum %d", config.Size.Height, capabilities.MaxHeight)
	}

	return nil
}

// renderLarge renders the large format display
func (r *RainbowTUIPlugin) renderLarge(data *domain.DisplayData) (string, error) {
	if data.Cost == nil {
		return "", nil
	}

	asciiArt := r.generateASCIIArtWithSize(data.Cost.TotalCost, data.Config.Size.Width, data.Config.Size.Height)
	centeredAsciiArt := r.centerASCIIArt(asciiArt, data.Config.Size.Width, data.Config.Size.Height)
	rainbowAsciiArt := r.applyRainbowColors(centeredAsciiArt, data.Animation)

	return rainbowAsciiArt, nil
}

// renderMedium renders the medium format display
func (r *RainbowTUIPlugin) renderMedium(data *domain.DisplayData) (string, error) {
	if data.Cost == nil {
		return "", nil
	}

	asciiArt := r.generateASCIIArtWithSize(data.Cost.TotalCost, data.Config.Size.Width, data.Config.Size.Height)
	centeredAsciiArt := r.centerASCIIArt(asciiArt, data.Config.Size.Width, data.Config.Size.Height)
	rainbowAsciiArt := r.applyRainbowColors(centeredAsciiArt, data.Animation)

	return rainbowAsciiArt, nil
}

// renderSmall renders the small format display
func (r *RainbowTUIPlugin) renderSmall(data *domain.DisplayData) (string, error) {
	if data.Cost == nil {
		return "", nil
	}

	asciiArt := r.generateASCIIArtWithSize(data.Cost.TotalCost, data.Config.Size.Width, data.Config.Size.Height)
	centeredAsciiArt := r.centerASCIIArt(asciiArt, data.Config.Size.Width, data.Config.Size.Height)
	rainbowAsciiArt := r.applyRainbowColors(centeredAsciiArt, data.Animation)

	return rainbowAsciiArt, nil
}

// renderMinimal renders the minimal format display
func (r *RainbowTUIPlugin) renderMinimal(data *domain.DisplayData) (string, error) {
	if data.Cost == nil {
		return "", nil
	}

	asciiArt := r.generateASCIIArtWithSize(data.Cost.TotalCost, data.Config.Size.Width, data.Config.Size.Height)
	centeredAsciiArt := r.centerASCIIArt(asciiArt, data.Config.Size.Width, data.Config.Size.Height)
	rainbowAsciiArt := r.applyRainbowColors(centeredAsciiArt, data.Animation)

	return rainbowAsciiArt, nil
}

// centerASCIIArt centers ASCII art both horizontally and vertically within given dimensions
func (r *RainbowTUIPlugin) centerASCIIArt(asciiArt string, width, height int) string {
	lines := strings.Split(asciiArt, "\n")
	if len(lines) == 0 {
		return ""
	}

	// Find the maximum line width
	maxLineWidth := 0
	for _, line := range lines {
		lineWidth := len([]rune(line)) // Use runes to handle Unicode properly
		if lineWidth > maxLineWidth {
			maxLineWidth = lineWidth
		}
	}

	// Calculate horizontal padding for centering
	horizontalPadding := 0
	if width > maxLineWidth {
		horizontalPadding = (width - maxLineWidth) / 2
	}

	// Calculate vertical padding for centering
	verticalPadding := 0
	if height > len(lines) {
		verticalPadding = (height - len(lines)) / 2
	}

	// Create centered output
	var result strings.Builder

	// Add top vertical padding
	for i := 0; i < verticalPadding; i++ {
		result.WriteString("\n")
	}

	// Center each line horizontally
	for i, line := range lines {
		if horizontalPadding > 0 {
			result.WriteString(strings.Repeat(" ", horizontalPadding))
		}
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	// Add bottom vertical padding
	for i := 0; i < verticalPadding; i++ {
		result.WriteString("\n")
	}

	return result.String()
}

// getSmallLetterPatterns returns small ASCII art patterns for small screens
func (r *RainbowTUIPlugin) getSmallLetterPatterns() map[rune][]string {
	return map[rune][]string{
		'$': {
			"    ███  ",
			" ███████ ",
			"███ ███  ",
			" ███████ ",
			"  ███ ███",
			" ███████ ",
			"   ███   ",
		},
		'0': {
			" ███████ ",
			"███   ███",
			"███   ███",
			"███   ███",
			"███   ███",
			"███   ███",
			" ███████ ",
		},
		'1': {
			"   ███   ",
			" █████   ",
			"   ███   ",
			"   ███   ",
			"   ███   ",
			"   ███   ",
			" ███████ ",
		},
		'2': {
			" ███████ ",
			"███   ███",
			"      ███",
			" ███████ ",
			"███      ",
			"███      ",
			"█████████",
		},
		'3': {
			" ███████ ",
			"███   ███",
			"      ███",
			"   █████ ",
			"      ███",
			"███   ███",
			" ███████ ",
		},
		'4': {
			"███   ███",
			"███   ███",
			"███   ███",
			"█████████",
			"      ███",
			"      ███",
			"      ███",
		},
		'5': {
			"█████████",
			"███      ",
			"███      ",
			"████████ ",
			"      ███",
			"███   ███",
			" ███████ ",
		},
		'6': {
			" ███████ ",
			"███   ███",
			"███      ",
			"████████ ",
			"███   ███",
			"███   ███",
			" ███████ ",
		},
		'7': {
			"█████████",
			"      ███",
			"     ███ ",
			"    ███  ",
			"   ███   ",
			"  ███    ",
			" ███     ",
		},
		'8': {
			" ███████ ",
			"███   ███",
			"███   ███",
			" ███████ ",
			"███   ███",
			"███   ███",
			" ███████ ",
		},
		'9': {
			" ███████ ",
			"███   ███",
			"███   ███",
			" ████████",
			"      ███",
			"███   ███",
			" ███████ ",
		},
		'.': {
			"      ",
			"      ",
			"      ",
			"      ",
			"      ",
			" ███  ",
			" ███  ",
		},
		' ': {
			"         ",
			"         ",
			"         ",
			"         ",
			"         ",
			"         ",
			"         ",
		},
	}
}

// getLargeLetterPatterns returns large ASCII art patterns for large screens
func (r *RainbowTUIPlugin) getLargeLetterPatterns() map[rune][]string {
	return map[rune][]string{
		'$': {
			"     ████     ",
			"  ███████████ ",
			" ████ ███     ",
			"████  ████    ",
			" ███████████  ",
			"  ███████████ ",
			"     ████ ████",
			"████████  ████",
			" ███████████  ",
			"     ████     ",
		},
		'0': {
			"  ██████████  ",
			" ████    ████ ",
			"████      ████",
			"████      ████",
			"████      ████",
			"████      ████",
			"████      ████",
			"████      ████",
			" ████    ████ ",
			"  ██████████  ",
		},
		'1': {
			"     ████     ",
			"  ███████     ",
			"     ████     ",
			"     ████     ",
			"     ████     ",
			"     ████     ",
			"     ████     ",
			"     ████     ",
			"     ████     ",
			"██████████████",
		},
		'2': {
			"  ███████████ ",
			" ████     ████",
			"          ████",
			"         ████ ",
			"       ████   ",
			"     ████     ",
			"   ████       ",
			" ████         ",
			"████          ",
			"██████████████",
		},
		'3': {
			"  ███████████ ",
			" ████     ████",
			"          ████",
			"          ████",
			"     █████████",
			"          ████",
			"          ████",
			"          ████",
			" ████     ████",
			"  ███████████ ",
		},
		'4': {
			"████      ████",
			"████      ████",
			"████      ████",
			"████      ████",
			"██████████████",
			"          ████",
			"          ████",
			"          ████",
			"          ████",
			"          ████",
		},
		'5': {
			"██████████████",
			"████          ",
			"████          ",
			"████          ",
			"█████████████ ",
			"          ████",
			"          ████",
			"          ████",
			" ████     ████",
			"  ███████████ ",
		},
		'6': {
			"  ███████████ ",
			" ████     ████",
			"████          ",
			"████          ",
			"█████████████ ",
			"████      ████",
			"████      ████",
			"████      ████",
			" ████     ████",
			"  ███████████ ",
		},
		'7': {
			"██████████████",
			"          ████",
			"         ████ ",
			"        ████  ",
			"       ████   ",
			"      ████    ",
			"     ████     ",
			"    ████      ",
			"   ████       ",
			"  ████        ",
		},
		'8': {
			"  ██████████  ",
			" ████    ████ ",
			"████      ████",
			" ████    ████ ",
			"  ██████████  ",
			" ████    ████ ",
			"████      ████",
			"████      ████",
			" ████    ████ ",
			"  ██████████  ",
		},
		'9': {
			"  ██████████  ",
			" ████    ████ ",
			"████      ████",
			"████      ████",
			" █████████████",
			"          ████",
			"          ████",
			"          ████",
			" ████     ███ ",
			"  ██████████  ",
		},
		'.': {
			"         ",
			"         ",
			"         ",
			"         ",
			"         ",
			"         ",
			"         ",
			" ██████  ",
			" ██████  ",
			" ██████  ",
		},
		' ': {
			"              ",
			"              ",
			"              ",
			"              ",
			"              ",
			"              ",
			"              ",
			"              ",
			"              ",
			"              ",
		},
	}
}

// generateASCIIArtWithSize converts a dollar amount to ASCII art with specified size
func (r *RainbowTUIPlugin) generateASCIIArtWithSize(amount float64, width, height int) string {
	text := fmt.Sprintf("$%.2f", amount)

	// Choose pattern set based on screen size
	var patterns map[rune][]string
	var numRows int

	// Use small patterns for smaller screens
	if width < 100 || height < 25 {
		patterns = r.getSmallLetterPatterns()
		numRows = 7
	} else {
		patterns = r.getLargeLetterPatterns()
		numRows = 10
	}

	// Build ASCII art line by line with spacing between characters
	lines := make([]string, numRows)
	for charIndex, char := range text {
		if pattern, exists := patterns[char]; exists {
			for i, line := range pattern {
				lines[i] += line
				// Add spacing between characters (except for the last character)
				if charIndex < len(text)-1 {
					lines[i] += "  " // 2 spaces between characters
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

// applyRainbowColors applies rainbow colors to text based on animation frame
func (r *RainbowTUIPlugin) applyRainbowColors(text string, animation *domain.AnimationFrame) string {
	if animation == nil || len(animation.Colors) == 0 {
		return text
	}

	var styledText strings.Builder
	lines := strings.Split(text, "\n")

	for lineIndex, line := range lines {
		for i, char := range line {
			colorIndex := (lineIndex*len(line) + i) % len(animation.Colors)
			color := animation.Colors[colorIndex]

			charStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			styledText.WriteString(charStyle.Render(string(char)))
		}
		if lineIndex < len(lines)-1 {
			styledText.WriteString("\n")
		}
	}

	return styledText.String()
}
