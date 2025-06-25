package integration_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestCLIFlags_AnimationOptionsIntegration(t *testing.T) {
	// Test animation flags already working
	args := []string{
		"--animation-speed", "50ms",
		"--animation-pattern", "gradient",
	}

	// Parse flags
	flagConfig, err := core.ParseCobraFlagsFromArgs(args)
	assert.NoError(t, err)

	// Create config manager and apply flags
	configManager := core.NewConfigManager()
	err = configManager.LoadConfig("")
	assert.NoError(t, err)

	err = configManager.ApplyFlagsToConfig(flagConfig)
	assert.NoError(t, err)

	// Check animation settings were applied
	animationConfig := configManager.GetAnimationConfig()
	assert.Equal(t, 50*time.Millisecond, animationConfig.Speed)
	assert.Equal(t, "gradient", string(animationConfig.Pattern))
	assert.True(t, animationConfig.Enabled) // should remain enabled
}

func TestCLIFlags_NoAnimationIntegration(t *testing.T) {
	// Test no-animation flag
	args := []string{"--no-animation"}

	// Parse flags
	flagConfig, err := core.ParseCobraFlagsFromArgs(args)
	assert.NoError(t, err)

	// Create config manager and apply flags
	configManager := core.NewConfigManager()
	err = configManager.LoadConfig("")
	assert.NoError(t, err)

	err = configManager.ApplyFlagsToConfig(flagConfig)
	assert.NoError(t, err)

	// Check animation was disabled
	animationConfig := configManager.GetAnimationConfig()
	assert.False(t, animationConfig.Enabled)
}

func TestCLIFlags_BankruptcyIntegration(t *testing.T) {
	// Test bankruptcy flag
	args := []string{"--bankruptcy"}

	// Parse flags
	flagConfig, err := core.ParseCobraFlagsFromArgs(args)
	assert.NoError(t, err)

	// Create config manager and apply flags
	configManager := core.NewConfigManager()
	err = configManager.LoadConfig("")
	assert.NoError(t, err)

	err = configManager.ApplyFlagsToConfig(flagConfig)
	assert.NoError(t, err)

	// Check bankruptcy flag was parsed correctly
	assert.True(t, flagConfig.Bankruptcy)
}

func TestCLIFlags_EndToEndValidation(t *testing.T) {
	// Test that flags can be parsed and applied without errors
	// This simulates the main.go flow more closely
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Animation speed only",
			args: []string{"--animation-speed", "100ms"},
		},
		{
			name: "Animation pattern only",
			args: []string{"--animation-pattern", "rainbow"},
		},
		{
			name: "No animation",
			args: []string{"--no-animation"},
		},
		{
			name: "Bankruptcy mode",
			args: []string{"--bankruptcy"},
		},
		{
			name: "Mixed animation options",
			args: []string{"--animation-speed", "200ms", "--animation-pattern", "wave"},
		},
		{
			name: "All supported options",
			args: []string{"--animation-speed", "75ms", "--animation-pattern", "pulse", "--bankruptcy"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This should not panic or error - simulating main.go flow
			flagConfig, err := core.ParseCobraFlagsFromArgs(tc.args)
			assert.NoError(t, err, "Flag parsing should succeed")

			configManager := core.NewConfigManager()
			err = configManager.LoadConfig("")
			assert.NoError(t, err, "Config loading should succeed")

			err = configManager.ApplyFlagsToConfig(flagConfig)
			assert.NoError(t, err, "Flag application should succeed")

			err = configManager.ValidateConfig()
			assert.NoError(t, err, "Config validation should succeed")

			// Verify configs can be retrieved without errors
			animationConfig := configManager.GetAnimationConfig()
			assert.NotNil(t, animationConfig, "Animation config should be available")
		})
	}
}

func TestCLIFlags_UnsupportedOptions(t *testing.T) {
	// Test that unsupported flags are rejected
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
			// These flags should cause an error
			_, err := core.ParseCobraFlagsFromArgs(tt.args)
			assert.Error(t, err, "Unsupported flag should cause an error")
		})
	}
}
