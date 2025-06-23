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
	var output strings.Builder

	// Title
	title := r.applyRainbowColors("💰 Claude Code Usage Cost", data.Animation)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(data.Config.Size.Width).
		MarginBottom(1)
	output.WriteString(titleStyle.Render(title))
	output.WriteString("\n")

	// Main cost display
	if data.Cost != nil {
		costText := fmt.Sprintf("$%.2f", data.Cost.TotalCost)
		rainbowCostText := r.applyRainbowColors(costText, data.Animation)

		costStyle := lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Width(data.Config.Size.Width).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Padding(2, 4).
			MarginBottom(1)

		output.WriteString(costStyle.Render(rainbowCostText))
		output.WriteString("\n")

		// Currency and timestamp
		if data.Config.ShowTimestamp {
			timestampText := fmt.Sprintf("Currency: %s | Last Updated: %s",
				data.Cost.Currency,
				data.Cost.Timestamp.Format("2006-01-02 15:04:05"))

			timestampStyle := lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(data.Config.Size.Width).
				Foreground(lipgloss.Color("#666666")).
				MarginBottom(1)

			output.WriteString(timestampStyle.Render(timestampText))
			output.WriteString("\n")
		}

		// Model breakdown
		if data.Config.ShowBreakdown && data.Cost.ModelBreakdown != nil && len(data.Cost.ModelBreakdown) > 0 {
			output.WriteString(r.renderModelBreakdown(data.Cost.ModelBreakdown, data))
		}
	}

	return output.String(), nil
}

// renderMedium renders the medium format display
func (r *RainbowTUIPlugin) renderMedium(data *domain.DisplayData) (string, error) {
	var output strings.Builder

	if data.Cost != nil {
		costText := fmt.Sprintf("💰 $%.2f", data.Cost.TotalCost)
		rainbowCostText := r.applyRainbowColors(costText, data.Animation)

		costStyle := lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Width(data.Config.Size.Width).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(1, 2).
			MarginBottom(1)

		output.WriteString(costStyle.Render(rainbowCostText))
		output.WriteString("\n")

		if data.Config.ShowTimestamp {
			timestampText := fmt.Sprintf("%s | %s",
				data.Cost.Currency,
				data.Cost.Timestamp.Format("15:04:05"))

			timestampStyle := lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(data.Config.Size.Width).
				Foreground(lipgloss.Color("#888888"))

			output.WriteString(timestampStyle.Render(timestampText))
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}

// renderSmall renders the small format display
func (r *RainbowTUIPlugin) renderSmall(data *domain.DisplayData) (string, error) {
	if data.Cost == nil {
		return "", nil
	}

	costText := fmt.Sprintf("$%.2f", data.Cost.TotalCost)
	rainbowCostText := r.applyRainbowColors(costText, data.Animation)

	costStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(data.Config.Size.Width)

	return costStyle.Render(rainbowCostText), nil
}

// renderMinimal renders the minimal format display
func (r *RainbowTUIPlugin) renderMinimal(data *domain.DisplayData) (string, error) {
	if data.Cost == nil {
		return "", nil
	}

	costText := fmt.Sprintf("$%.2f", data.Cost.TotalCost)
	return r.applyRainbowColors(costText, data.Animation), nil
}

// renderModelBreakdown renders the model cost breakdown
func (r *RainbowTUIPlugin) renderModelBreakdown(breakdown map[string]float64, data *domain.DisplayData) string {
	var output strings.Builder

	breakdownTitle := r.applyRainbowColors("📊 Model Breakdown", data.Animation)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(data.Config.Size.Width).
		MarginBottom(1)

	output.WriteString(titleStyle.Render(breakdownTitle))
	output.WriteString("\n")

	for model, cost := range breakdown {
		modelText := fmt.Sprintf("%-20s $%.2f", model, cost)
		rainbowModelText := r.applyRainbowColors(modelText, data.Animation)

		modelStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(data.Config.Size.Width).
			PaddingLeft(4)

		output.WriteString(modelStyle.Render(rainbowModelText))
		output.WriteString("\n")
	}

	return output.String()
}

// applyRainbowColors applies rainbow colors to text based on animation frame
func (r *RainbowTUIPlugin) applyRainbowColors(text string, animation *domain.AnimationFrame) string {
	if animation == nil || len(animation.Colors) == 0 {
		return text
	}

	var styledText strings.Builder

	for i, char := range text {
		colorIndex := i % len(animation.Colors)
		color := animation.Colors[colorIndex]

		charStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		styledText.WriteString(charStyle.Render(string(char)))
	}

	return styledText.String()
}
