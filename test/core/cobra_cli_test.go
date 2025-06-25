package core_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestCobraCLI_RootCommand tests the root command functionality
func TestCobraCLI_RootCommand(t *testing.T) {
	// Red phase: Test that cobra root command can be created
	rootCmd := core.NewRootCommand()
	assert.NotNil(t, rootCmd, "Root command should be created")
	assert.Equal(t, "ccugorg", rootCmd.Use, "Root command should have correct name")
	assert.NotEmpty(t, rootCmd.Short, "Root command should have short description")
}

// TestCobraCLI_AnimationSpeedFlag tests animation speed flag with cobra
func TestCobraCLI_AnimationSpeedFlag(t *testing.T) {
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
			// Use ParseCobraFlagsFromArgs directly with args
			flagConfig, err := core.ParseCobraFlagsFromArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, flagConfig.Animation.Speed)
			}
		})
	}
}

// TestCobraCLI_AnimationPatternFlag tests animation pattern flag with cobra
func TestCobraCLI_AnimationPatternFlag(t *testing.T) {
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
			flagConfig, err := core.ParseCobraFlagsFromArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, flagConfig.Animation.Pattern)
			}
		})
	}
}

// TestCobraCLI_NoAnimationFlag tests no-animation flag with cobra
func TestCobraCLI_NoAnimationFlag(t *testing.T) {
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
			flagConfig, err := core.ParseCobraFlagsFromArgs(tt.args)
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

// TestCobraCLI_BankruptcyFlag tests hidden bankruptcy flag with cobra
func TestCobraCLI_BankruptcyFlag(t *testing.T) {
	flagConfig, err := core.ParseCobraFlagsFromArgs([]string{"--bankruptcy"})
	assert.NoError(t, err, "Bankruptcy flag should be supported")
	assert.True(t, flagConfig.Bankruptcy, "Bankruptcy flag should be set")
}

// TestCobraCLI_UnsupportedFlags tests that unsupported flags are rejected
func TestCobraCLI_UnsupportedFlags(t *testing.T) {
	unsupportedFlags := []struct {
		name string
		args []string
	}{
		{
			name: "Format flag should be unsupported",
			args: []string{"--format", "large"},
		},
		{
			name: "Width flag should be unsupported",
			args: []string{"--width", "100"},
		},
		{
			name: "Height flag should be unsupported",
			args: []string{"--height", "30"},
		},
		{
			name: "Show timestamp flag should be unsupported",
			args: []string{"--show-timestamp"},
		},
		{
			name: "Show breakdown flag should be unsupported",
			args: []string{"--show-breakdown"},
		},
	}

	for _, tt := range unsupportedFlags {
		t.Run(tt.name, func(t *testing.T) {
			_, err := core.ParseCobraFlagsFromArgs(tt.args)
			assert.Error(t, err, "Unsupported flag should cause an error")
		})
	}
}

// TestCobraCLI_HelpText tests that help text doesn't show bankruptcy flag
func TestCobraCLI_HelpText(t *testing.T) {
	// Red phase: Test that help text is correctly generated and bankruptcy is hidden
	rootCmd := core.NewRootCommand()

	helpText := rootCmd.UsageString()
	assert.NotEmpty(t, helpText, "Help text should be generated")
	assert.Contains(t, helpText, "animation-speed", "Help should contain animation-speed flag")
	assert.Contains(t, helpText, "animation-pattern", "Help should contain animation-pattern flag")
	assert.Contains(t, helpText, "no-animation", "Help should contain no-animation flag")
	assert.NotContains(t, helpText, "bankruptcy", "Help should NOT contain bankruptcy flag (hidden)")
}

// TestCobraCLI_ConfigIntegration tests integration with config manager
func TestCobraCLI_ConfigIntegration(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Parse cobra flags
	flagConfig, err := core.ParseCobraFlagsFromArgs([]string{"--animation-speed", "75ms", "--no-animation"})
	assert.NoError(t, err)

	// Apply flags to config
	err = configManager.ApplyFlagsToConfig(flagConfig)
	assert.NoError(t, err)

	// Verify flags were applied
	animationConfig := configManager.GetAnimationConfig()
	assert.Equal(t, 75*time.Millisecond, animationConfig.Speed)
	assert.False(t, animationConfig.Enabled)
}
