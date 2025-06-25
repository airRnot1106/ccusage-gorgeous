package core_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigManager(t *testing.T) {
	cm := core.NewConfigManager()
	assert.NotNil(t, cm)
}

func TestConfigManager_LoadConfig_Defaults(t *testing.T) {
	cm := core.NewConfigManager()

	// LoadConfig is now a no-op but should still work for compatibility
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	config := cm.GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "info", config.App.LogLevel)
	assert.Equal(t, 1*time.Second, config.App.RefreshRate)
	assert.Equal(t, 80, config.Display.Width)
	assert.Equal(t, 24, config.Display.Height)
	assert.True(t, config.Animation.Enabled)
	assert.Equal(t, 100*time.Millisecond, config.Animation.Speed)
	assert.Equal(t, domain.PatternRainbow, config.Animation.Pattern)
	assert.Len(t, config.Animation.Colors, 12)
}

// TestConfigManager_LoadConfig_FromFile removed since file loading is no longer supported

// TestConfigManager_LoadConfig_InvalidYAML removed since file loading is no longer supported

func TestConfigManager_GetDisplayConfig(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	displayConfig := cm.GetDisplayConfig()
	assert.NotNil(t, displayConfig)
	assert.Equal(t, 1*time.Second, displayConfig.RefreshRate)
	assert.Equal(t, 80, displayConfig.Size.Width)
	assert.Equal(t, 24, displayConfig.Size.Height)
}

func TestConfigManager_GetAnimationConfig(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	animationConfig := cm.GetAnimationConfig()
	assert.NotNil(t, animationConfig)
	assert.Equal(t, 100*time.Millisecond, animationConfig.Speed)
	assert.Len(t, animationConfig.Colors, 12)
	assert.True(t, animationConfig.Enabled)
	assert.Equal(t, domain.PatternRainbow, animationConfig.Pattern)
}

func TestConfigManager_GetConfig_Nil(t *testing.T) {
	cm := core.NewConfigManager()

	// Config is now loaded with defaults in NewConfigManager
	config := cm.GetConfig()
	assert.NotNil(t, config)

	displayConfig := cm.GetDisplayConfig()
	assert.NotNil(t, displayConfig)

	animationConfig := cm.GetAnimationConfig()
	assert.NotNil(t, animationConfig)
}

func TestConfigManager_UpdateConfig(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// UpdateConfig now only supports plugin datasource updates
	updates := map[string]interface{}{
		"plugins.datasource": "bankruptcy-datasource",
	}

	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	config := cm.GetConfig()
	assert.Equal(t, "bankruptcy-datasource", config.Plugins.DataSource)
}

func TestConfigManager_ValidateConfig(t *testing.T) {
	cm := core.NewConfigManager()

	// Config is always loaded with defaults now, so validation should pass
	err := cm.ValidateConfig()
	assert.NoError(t, err)

	// Load config (no-op but for compatibility)
	err = cm.LoadConfig("")
	assert.NoError(t, err)

	// Should still pass validation
	err = cm.ValidateConfig()
	assert.NoError(t, err)
}

// Format validation test removed since display formats are no longer supported

func TestConfigManager_ValidateConfig_InvalidPattern(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Manually set invalid pattern for testing
	config := cm.GetConfig()
	config.Animation.Pattern = "invalid-pattern"

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid animation pattern")
}

func TestConfigManager_ValidateConfig_InvalidDimensions(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Manually set invalid dimensions for testing
	config := cm.GetConfig()
	config.Display.Width = -1
	config.Display.Height = 0

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "display dimensions must be positive")
}

func TestConfigManager_ValidateConfig_InvalidRefreshRate(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Manually set invalid refresh rate for testing
	config := cm.GetConfig()
	config.App.RefreshRate = -1 * time.Second

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh rate must be positive")
}

func TestConfigManager_ValidateConfig_InvalidAnimationSpeed(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Manually set invalid animation speed for testing
	config := cm.GetConfig()
	config.Animation.Speed = -100 * time.Millisecond

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animation speed must be positive")
}
