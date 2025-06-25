package core

import (
	"fmt"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/spf13/cobra"
)

// FlagConfig represents the command line flag configuration
type FlagConfig struct {
	Animation struct {
		Speed   time.Duration
		Pattern domain.AnimationPattern
		Enabled *bool
	}
	Bankruptcy bool
}

// NewRootCommand creates the root cobra command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ccugorg",
		Short: "TUI application for displaying Claude API costs with rainbow animations",
		Long: `ccugorg is a terminal user interface application that displays
Claude API usage costs with beautiful rainbow animations and ASCII art.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add flags
	cmd.Flags().String("animation-speed", "", "Animation speed (e.g., 100ms)")
	cmd.Flags().String("animation-pattern", "", "Animation pattern (rainbow, gradient, pulse, wave)")
	cmd.Flags().Bool("no-animation", false, "Disable animation")

	// Hidden bankruptcy flag
	cmd.Flags().Bool("bankruptcy", false, "")
	_ = cmd.Flags().MarkHidden("bankruptcy") // Hide bankruptcy flag from help

	return cmd
}

// ParseCobraFlagsFromArgs parses cobra command flags from args directly
func ParseCobraFlagsFromArgs(args []string) (*FlagConfig, error) {
	cmd := NewRootCommand()
	cmd.SetArgs(args)

	// Parse flags
	if err := cmd.ParseFlags(args); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	flagConfig := &FlagConfig{}

	// Parse animation speed
	speedStr, _ := cmd.Flags().GetString("animation-speed")
	if speedStr != "" {
		speed, err := time.ParseDuration(speedStr)
		if err != nil {
			return nil, fmt.Errorf("invalid animation speed format '%s': %w", speedStr, err)
		}
		flagConfig.Animation.Speed = speed
	}

	// Parse animation pattern
	patternStr, _ := cmd.Flags().GetString("animation-pattern")
	if patternStr != "" {
		pattern := domain.AnimationPattern(patternStr)
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
			return nil, fmt.Errorf("invalid animation pattern '%s'. Valid patterns: rainbow, gradient, pulse, wave", patternStr)
		}
		flagConfig.Animation.Pattern = pattern
	}

	// Parse no-animation flag
	noAnimation, _ := cmd.Flags().GetBool("no-animation")
	if noAnimation {
		enabled := false
		flagConfig.Animation.Enabled = &enabled
	}

	// Parse bankruptcy flag
	bankruptcy, _ := cmd.Flags().GetBool("bankruptcy")
	flagConfig.Bankruptcy = bankruptcy

	return flagConfig, nil
}

// ParseCobraFlags parses cobra command flags and returns flag configuration (for backwards compatibility)
func ParseCobraFlags(cmd *cobra.Command) (*FlagConfig, error) {
	return ParseCobraFlagsFromArgs(cmd.Flags().Args())
}
