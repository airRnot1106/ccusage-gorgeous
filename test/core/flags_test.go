package core_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestRequiredFlags tests only the flags that should be supported
func TestRequiredFlags_ShouldNotExist(t *testing.T) {
	// Red phase: Test that unsupported flags are rejected
	unsupportedFlags := []struct {
		name string
		args []string
	}{
		{
			name: "Format flag should not exist",
			args: []string{"--format", "large"},
		},
		{
			name: "Width flag should not exist",
			args: []string{"--width", "100"},
		},
		{
			name: "Height flag should not exist",
			args: []string{"--height", "30"},
		},
		{
			name: "Show timestamp flag should not exist",
			args: []string{"--show-timestamp"},
		},
		{
			name: "No timestamp flag should not exist",
			args: []string{"--no-timestamp"},
		},
		{
			name: "Show breakdown flag should not exist",
			args: []string{"--show-breakdown"},
		},
		{
			name: "No breakdown flag should not exist",
			args: []string{"--no-breakdown"},
		},
	}

	for _, tt := range unsupportedFlags {
		t.Run(tt.name, func(t *testing.T) {
			// These flags should cause an error because they are not supported
			_, err := core.ParseFlagsFromArgs(tt.args)
			assert.Error(t, err, "Unsupported flag should cause an error")
		})
	}
}

func TestSupportedFlags_ShouldWork(t *testing.T) {
	// Test that only supported flags work correctly
	supportedFlags := []struct {
		name string
		args []string
	}{
		{
			name: "Animation speed flag should work",
			args: []string{"--animation-speed", "50ms"},
		},
		{
			name: "Animation pattern flag should work",
			args: []string{"--animation-pattern", "rainbow"},
		},
		{
			name: "No animation flag should work",
			args: []string{"--no-animation"},
		},
		{
			name: "Bankruptcy flag should work",
			args: []string{"--bankruptcy"},
		},
	}

	for _, tt := range supportedFlags {
		t.Run(tt.name, func(t *testing.T) {
			// These flags should not cause an error because they are supported
			_, err := core.ParseFlagsFromArgs(tt.args)
			assert.NoError(t, err, "Supported flag should not cause an error")
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
			if flagConfig.Animation.Enabled != nil {
				assert.Equal(t, tt.expected, *flagConfig.Animation.Enabled)
			} else {
				// If Enabled is nil, assume default enabled behavior
				assert.True(t, tt.expected)
			}
		})
	}
}

func TestBankruptcyFlag(t *testing.T) {
	// Test that bankruptcy flag should work
	_, err := core.ParseFlagsFromArgs([]string{"--bankruptcy"})
	assert.NoError(t, err, "Bankruptcy flag should be supported")
}

func TestConfigOverrideOrder(t *testing.T) {
	// Create a test config manager
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Parse flags that should override config (only supported flags)
	flagConfig, err := core.ParseFlagsFromArgs([]string{"--animation-speed", "50ms", "--no-animation"})
	assert.NoError(t, err)

	err = configManager.ApplyFlagsToConfig(flagConfig)
	assert.NoError(t, err)

	// Command line flags should override config file values
	animationConfig := configManager.GetAnimationConfig()
	assert.Equal(t, 50*time.Millisecond, animationConfig.Speed)
	assert.False(t, animationConfig.Enabled)
}

func TestInvalidFlagValues(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Invalid animation speed",
			args: []string{"--animation-speed", "invalid"},
		},
		{
			name: "Invalid animation pattern",
			args: []string{"--animation-pattern", "invalid"},
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
