package core

import (
	"errors"
	"flag"
	"strconv"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// FlagConfig represents configuration from command line flags
type FlagConfig struct {
	Display   FlagDisplayConfig
	Animation FlagAnimationConfig
}

// FlagDisplayConfig represents display configuration from flags
type FlagDisplayConfig struct {
	Format        domain.DisplayFormat
	Width         int
	Height        int
	ShowTimestamp *bool // pointer to distinguish between set and unset
	ShowBreakdown *bool // pointer to distinguish between set and unset
}

// FlagAnimationConfig represents animation configuration from flags
type FlagAnimationConfig struct {
	Speed   time.Duration
	Pattern domain.AnimationPattern
	Enabled *bool // pointer to distinguish between set and unset
}

// ParseFlags parses command line flags and returns FlagConfig
func ParseFlags() (*FlagConfig, error) {
	return ParseFlagsFromArgs(flag.Args())
}

// ParseFlagsFromArgs parses flags from given arguments
func ParseFlagsFromArgs(args []string) (*FlagConfig, error) {
	flagSet := flag.NewFlagSet("ccugorg", flag.ContinueOnError)
	config := &FlagConfig{}

	// Define flags
	formatFlag := flagSet.String("format", "", "Display format (large, medium, small, minimal)")
	widthFlag := flagSet.Int("width", 0, "Display width")
	heightFlag := flagSet.Int("height", 0, "Display height")
	flagSet.Bool("show-timestamp", false, "Show timestamp")
	flagSet.Bool("no-timestamp", false, "Hide timestamp")
	flagSet.Bool("show-breakdown", false, "Show breakdown")
	flagSet.Bool("no-breakdown", false, "Hide breakdown")
	animationSpeedFlag := flagSet.String("animation-speed", "", "Animation speed (e.g., 100ms)")
	animationPatternFlag := flagSet.String("animation-pattern", "", "Animation pattern (rainbow, gradient, pulse, wave)")
	flagSet.Bool("no-animation", false, "Disable animation")
	flagSet.Bool("bankruptcy", false, "") // Hidden flag

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	// Process flags using helper functions
	if err := processDisplayFlags(flagSet, config, *formatFlag, *widthFlag, *heightFlag); err != nil {
		return nil, err
	}

	if err := processAnimationFlags(flagSet, config, *animationSpeedFlag, *animationPatternFlag); err != nil {
		return nil, err
	}

	return config, nil
}

// processDisplayFlags processes display-related flags
func processDisplayFlags(flagSet *flag.FlagSet, config *FlagConfig, format string, width, height int) error {
	if format != "" {
		if err := parseDisplayFormat(format, &config.Display.Format); err != nil {
			return err
		}
	}

	if width < 0 || height < 0 {
		return errors.New("width and height must be positive")
	}
	config.Display.Width = width
	config.Display.Height = height

	processBooleanFlags(flagSet, config)
	return nil
}

// processAnimationFlags processes animation-related flags
func processAnimationFlags(flagSet *flag.FlagSet, config *FlagConfig, speed, pattern string) error {
	if speed != "" {
		duration, err := time.ParseDuration(speed)
		if err != nil {
			return errors.New("invalid animation speed format")
		}
		config.Animation.Speed = duration
	}

	if pattern != "" {
		if err := parseAnimationPattern(pattern, &config.Animation.Pattern); err != nil {
			return err
		}
	}

	// Process no-animation flag
	val := flagSet.Lookup("no-animation").Value.String() != "true"
	config.Animation.Enabled = &val
	return nil
}

// processBooleanFlags processes boolean display flags
func processBooleanFlags(flagSet *flag.FlagSet, config *FlagConfig) {
	if flagSet.Lookup("show-timestamp").Value.String() == "true" {
		val := true
		config.Display.ShowTimestamp = &val
	}
	if flagSet.Lookup("no-timestamp").Value.String() == "true" {
		val := false
		config.Display.ShowTimestamp = &val
	}
	if flagSet.Lookup("show-breakdown").Value.String() == "true" {
		val := true
		config.Display.ShowBreakdown = &val
	}
	if flagSet.Lookup("no-breakdown").Value.String() == "true" {
		val := false
		config.Display.ShowBreakdown = &val
	}
}

// parseDisplayFormat parses display format string
func parseDisplayFormat(format string, target *domain.DisplayFormat) error {
	switch format {
	case "large":
		*target = domain.FormatLarge
	case "medium":
		*target = domain.FormatMedium
	case "small":
		*target = domain.FormatSmall
	case "minimal":
		*target = domain.FormatMinimal
	default:
		return errors.New("invalid format: must be large, medium, small, or minimal")
	}
	return nil
}

// parseAnimationPattern parses animation pattern string
func parseAnimationPattern(pattern string, target *domain.AnimationPattern) error {
	switch pattern {
	case "rainbow":
		*target = domain.PatternRainbow
	case "gradient":
		*target = domain.PatternGradient
	case "pulse":
		*target = domain.PatternPulse
	case "wave":
		*target = domain.PatternWave
	default:
		return errors.New("invalid animation pattern: must be rainbow, gradient, pulse, or wave")
	}
	return nil
}

// ValidateFlagValue validates individual flag values
func ValidateFlagValue(flagName, value string) error {
	switch flagName {
	case "format":
		var dummy domain.DisplayFormat
		return parseDisplayFormat(value, &dummy)
	case "width", "height":
		if val, err := strconv.Atoi(value); err != nil || val < 0 {
			return errors.New(flagName + " must be a positive number")
		}
	case "animation-speed":
		if _, err := time.ParseDuration(value); err != nil {
			return errors.New("invalid animation speed format")
		}
	case "animation-pattern":
		var dummy domain.AnimationPattern
		return parseAnimationPattern(value, &dummy)
	}
	return nil
}
