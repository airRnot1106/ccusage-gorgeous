package core_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDisplayFormatFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected domain.DisplayFormat
		wantErr  bool
	}{
		{
			name:     "Large format flag",
			args:     []string{"--format", "large"},
			expected: domain.FormatLarge,
			wantErr:  false,
		},
		{
			name:     "Medium format flag",
			args:     []string{"--format", "medium"},
			expected: domain.FormatMedium,
			wantErr:  false,
		},
		{
			name:     "Small format flag",
			args:     []string{"--format", "small"},
			expected: domain.FormatSmall,
			wantErr:  false,
		},
		{
			name:     "Minimal format flag",
			args:     []string{"--format", "minimal"},
			expected: domain.FormatMinimal,
			wantErr:  false,
		},
		{
			name:     "Invalid format flag",
			args:     []string{"--format", "invalid"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags directly from test args
			flagConfig, err := core.ParseFlagsFromArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, flagConfig.Display.Format)
			}
		})
	}
}

func TestAnimationSpeedFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected time.Duration
		wantErr  bool
	}{
		{
			name:     "Valid speed 50ms",
			args:     []string{"--animation-speed", "50ms"},
			expected: 50 * time.Millisecond,
			wantErr:  false,
		},
		{
			name:     "Valid speed 1s",
			args:     []string{"--animation-speed", "1s"},
			expected: 1 * time.Second,
			wantErr:  false,
		},
		{
			name:     "Invalid speed format",
			args:     []string{"--animation-speed", "invalid"},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags directly from test args
			flagConfig, err := core.ParseFlagsFromArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, flagConfig.Animation.Speed)
			}
		})
	}
}

func TestAnimationPatternFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected domain.AnimationPattern
		wantErr  bool
	}{
		{
			name:     "Rainbow pattern",
			args:     []string{"--animation-pattern", "rainbow"},
			expected: domain.PatternRainbow,
			wantErr:  false,
		},
		{
			name:     "Gradient pattern",
			args:     []string{"--animation-pattern", "gradient"},
			expected: domain.PatternGradient,
			wantErr:  false,
		},
		{
			name:     "Pulse pattern",
			args:     []string{"--animation-pattern", "pulse"},
			expected: domain.PatternPulse,
			wantErr:  false,
		},
		{
			name:     "Wave pattern",
			args:     []string{"--animation-pattern", "wave"},
			expected: domain.PatternWave,
			wantErr:  false,
		},
		{
			name:     "Invalid pattern",
			args:     []string{"--animation-pattern", "invalid"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags directly from test args
			flagConfig, err := core.ParseFlagsFromArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, flagConfig.Animation.Pattern)
			}
		})
	}
}

func TestNoAnimationFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "Animation enabled by default",
			args:     []string{},
			expected: true,
		},
		{
			name:     "Animation disabled with --no-animation",
			args:     []string{"--no-animation"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags directly from test args
			flagConfig, err := core.ParseFlagsFromArgs(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *flagConfig.Animation.Enabled)
		})
	}
}

func TestConfigOverrideOrder(t *testing.T) {
	// Create a test config manager
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Parse flags that should override config
	flagConfig, err := core.ParseFlagsFromArgs([]string{"--format", "small", "--animation-speed", "50ms"})
	assert.NoError(t, err)

	err = configManager.ApplyFlagsToConfig(flagConfig)
	assert.NoError(t, err)

	// Command line flags should override config file values
	displayConfig := configManager.GetDisplayConfig()
	assert.Equal(t, domain.FormatSmall, displayConfig.Format)

	animationConfig := configManager.GetAnimationConfig()
	assert.Equal(t, 50*time.Millisecond, animationConfig.Speed)
}

func TestInvalidFlagValues(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Invalid format value",
			args: []string{"--format", "invalid"},
		},
		{
			name: "Invalid animation speed",
			args: []string{"--animation-speed", "invalid"},
		},
		{
			name: "Invalid animation pattern",
			args: []string{"--animation-pattern", "invalid"},
		},
		{
			name: "Invalid width value (must be positive)",
			args: []string{"--width", "-10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse flags directly from test args
			_, err := core.ParseFlagsFromArgs(tt.args)
			assert.Error(t, err, "Expected error for invalid flag value")
		})
	}
}
