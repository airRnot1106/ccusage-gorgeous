package core

import (
	"errors"
	"flag"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// FlagConfig represents configuration from command line flags
type FlagConfig struct {
	Animation  FlagAnimationConfig
	Bankruptcy bool
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

	// Define only supported flags
	animationSpeedFlag := flagSet.String("animation-speed", "", "Animation speed (e.g., 100ms)")
	animationPatternFlag := flagSet.String("animation-pattern", "", "Animation pattern (rainbow, gradient, pulse, wave)")
	noAnimationFlag := flagSet.Bool("no-animation", false, "Disable animation")
	bankruptcyFlag := flagSet.Bool("bankruptcy", false, "") // Hidden flag

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	// Process animation flags
	if err := processAnimationFlags(config, *animationSpeedFlag, *animationPatternFlag, *noAnimationFlag); err != nil {
		return nil, err
	}

	// Process bankruptcy flag
	config.Bankruptcy = *bankruptcyFlag

	return config, nil
}

// processAnimationFlags processes animation-related flags
func processAnimationFlags(config *FlagConfig, speed, pattern string, noAnimation bool) error {
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

	// Process no-animation flag - only set if flag was provided
	if noAnimation {
		val := false
		config.Animation.Enabled = &val
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
