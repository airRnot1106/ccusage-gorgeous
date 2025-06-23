package core_test

import (
	"os"
	"path/filepath"
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

	// Load with non-existent file should work with defaults
	err := cm.LoadConfig("non-existent-file.yaml")
	assert.NoError(t, err)

	config := cm.GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "info", config.App.LogLevel)
	assert.Equal(t, 1*time.Second, config.App.RefreshRate)
	assert.Equal(t, domain.FormatLarge, config.Display.Format)
	assert.Equal(t, 80, config.Display.Width)
	assert.Equal(t, 24, config.Display.Height)
	assert.True(t, config.Display.ShowTimestamp)
	assert.True(t, config.Display.ShowBreakdown)
	assert.True(t, config.Animation.Enabled)
	assert.Equal(t, 100*time.Millisecond, config.Animation.Speed)
	assert.Equal(t, domain.PatternRainbow, config.Animation.Pattern)
	assert.Len(t, config.Animation.Colors, 12)
}

func TestConfigManager_LoadConfig_FromFile(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	configContent := `
app:
  log_level: "debug"
  refresh_rate: "2s"

display:
  format: "medium"
  width: 120
  height: 30
  show_timestamp: false
  show_breakdown: false

animation:
  enabled: false
  speed: "200ms"
  pattern: "gradient"
  colors: ["#FF0000", "#00FF00", "#0000FF"]

datasource:
  ccusage_path: "/usr/local/bin/ccusage"
  timeout: "60s"
  cache_time: "20s"

plugins:
  datasource: "custom-datasource"
  display: "custom-display"
  animation: "custom-animation"
`

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cm := core.NewConfigManager()
	err = cm.LoadConfig(configPath)
	assert.NoError(t, err)

	config := cm.GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "debug", config.App.LogLevel)
	assert.Equal(t, 2*time.Second, config.App.RefreshRate)
	assert.Equal(t, domain.FormatMedium, config.Display.Format)
	assert.Equal(t, 120, config.Display.Width)
	assert.Equal(t, 30, config.Display.Height)
	assert.False(t, config.Display.ShowTimestamp)
	assert.False(t, config.Display.ShowBreakdown)
	assert.False(t, config.Animation.Enabled)
	assert.Equal(t, 200*time.Millisecond, config.Animation.Speed)
	assert.Equal(t, domain.PatternGradient, config.Animation.Pattern)
	assert.Len(t, config.Animation.Colors, 3)
	assert.Equal(t, "/usr/local/bin/ccusage", config.DataSource.CcusagePath)
	assert.Equal(t, 60*time.Second, config.DataSource.Timeout)
	assert.Equal(t, 20*time.Second, config.DataSource.CacheTime)
	assert.Equal(t, "custom-datasource", config.Plugins.DataSource)
	assert.Equal(t, "custom-display", config.Plugins.Display)
	assert.Equal(t, "custom-animation", config.Plugins.Animation)
}

func TestConfigManager_LoadConfig_InvalidYAML(t *testing.T) {
	// Create temporary invalid config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid-config.yaml")

	invalidContent := `
app:
  log_level: "debug
  # Missing closing quote - invalid YAML
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0o644)
	assert.NoError(t, err)

	cm := core.NewConfigManager()
	err = cm.LoadConfig(configPath)
	assert.Error(t, err)
}

func TestConfigManager_GetDisplayConfig(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	displayConfig := cm.GetDisplayConfig()
	assert.NotNil(t, displayConfig)
	assert.Equal(t, 1*time.Second, displayConfig.RefreshRate)
	assert.True(t, displayConfig.ShowTimestamp)
	assert.True(t, displayConfig.ShowBreakdown)
	assert.Equal(t, domain.FormatLarge, displayConfig.Format)
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

	// Before loading config, should return nil
	config := cm.GetConfig()
	assert.Nil(t, config)

	displayConfig := cm.GetDisplayConfig()
	assert.Nil(t, displayConfig)

	animationConfig := cm.GetAnimationConfig()
	assert.Nil(t, animationConfig)
}

func TestConfigManager_UpdateConfig(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Update some values
	updates := map[string]interface{}{
		"app.log_level":     "error",
		"display.width":     100,
		"animation.enabled": false,
		"animation.pattern": "pulse",
	}

	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	config := cm.GetConfig()
	assert.Equal(t, "error", config.App.LogLevel)
	assert.Equal(t, 100, config.Display.Width)
	assert.False(t, config.Animation.Enabled)
	assert.Equal(t, domain.PatternPulse, config.Animation.Pattern)
}

func TestConfigManager_ValidateConfig(t *testing.T) {
	cm := core.NewConfigManager()

	// Should error when no config is loaded
	err := cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no configuration loaded")

	// Load valid config
	err = cm.LoadConfig("")
	assert.NoError(t, err)

	// Should pass validation
	err = cm.ValidateConfig()
	assert.NoError(t, err)
}

func TestConfigManager_ValidateConfig_InvalidFormat(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Update with invalid display format
	updates := map[string]interface{}{
		"display.format": "invalid-format",
	}
	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid display format")
}

func TestConfigManager_ValidateConfig_InvalidPattern(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Update with invalid animation pattern
	updates := map[string]interface{}{
		"animation.pattern": "invalid-pattern",
	}
	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid animation pattern")
}

func TestConfigManager_ValidateConfig_InvalidDimensions(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Update with invalid dimensions
	updates := map[string]interface{}{
		"display.width":  -1,
		"display.height": 0,
	}
	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid display dimensions")
}

func TestConfigManager_ValidateConfig_InvalidRefreshRate(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Update with invalid refresh rate
	updates := map[string]interface{}{
		"app.refresh_rate": "-1s",
	}
	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh rate")
}

func TestConfigManager_ValidateConfig_InvalidAnimationSpeed(t *testing.T) {
	cm := core.NewConfigManager()
	err := cm.LoadConfig("")
	assert.NoError(t, err)

	// Update with invalid animation speed
	updates := map[string]interface{}{
		"animation.speed": "-100ms",
	}
	err = cm.UpdateConfig(updates)
	assert.NoError(t, err)

	// Should fail validation
	err = cm.ValidateConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid animation speed")
}
