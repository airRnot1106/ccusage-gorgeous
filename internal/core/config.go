package core

import (
	"fmt"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// Config represents the application configuration
type Config struct {
	App        AppConfig
	Display    DisplayConfig
	Animation  AnimationConfig
	DataSource DataSourceConfig
	Plugins    PluginsConfig
}

// AppConfig represents general application settings
type AppConfig struct {
	LogLevel    string
	RefreshRate time.Duration
}

// DisplayConfig represents display-specific settings
type DisplayConfig struct {
	Width  int
	Height int
}

// AnimationConfig represents animation-specific settings
type AnimationConfig struct {
	Enabled bool
	Speed   time.Duration
	Pattern domain.AnimationPattern
	Colors  []string
}

// DataSourceConfig represents data source settings
type DataSourceConfig struct {
	CcusagePath string
	Timeout     time.Duration
	CacheTime   time.Duration
}

// PluginsConfig represents plugin configuration
type PluginsConfig struct {
	DataSource string
	Display    string
	Animation  string
}

// ConfigManager provides configuration management functionality
type ConfigManager struct {
	config *Config
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config: getDefaultConfig(),
	}
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			LogLevel:    "info",
			RefreshRate: 1 * time.Second,
		},
		Display: DisplayConfig{
			Width:  80,
			Height: 24,
		},
		Animation: AnimationConfig{
			Enabled: true,
			Speed:   100 * time.Millisecond,
			Pattern: domain.PatternRainbow,
			Colors: []string{
				"#FF0000", // Red
				"#FF8000", // Orange
				"#FFFF00", // Yellow
				"#80FF00", // Light Green
				"#00FF00", // Green
				"#00FF80", // Cyan Green
				"#00FFFF", // Cyan
				"#0080FF", // Light Blue
				"#0000FF", // Blue
				"#8000FF", // Purple
				"#FF00FF", // Magenta
				"#FF0080", // Pink
			},
		},
		DataSource: DataSourceConfig{
			CcusagePath: "ccusage",
			Timeout:     30 * time.Second,
			CacheTime:   10 * time.Second,
		},
		Plugins: PluginsConfig{
			DataSource: "ccusage-cli",
			Display:    "rainbow-display",
			Animation:  "rainbow-animation",
		},
	}
}

// LoadConfig loads configuration with defaults only (no file loading)
func (cm *ConfigManager) LoadConfig(configPath string) error {
	// Configuration is already loaded with defaults in NewConfigManager
	// This method is kept for compatibility but doesn't load from files
	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// GetDisplayConfig converts core DisplayConfig to domain DisplayConfig
func (cm *ConfigManager) GetDisplayConfig() *domain.DisplayConfig {
	if cm.config == nil {
		return nil
	}

	return &domain.DisplayConfig{
		RefreshRate: cm.config.App.RefreshRate,
		Size: domain.DisplaySize{
			Width:  cm.config.Display.Width,
			Height: cm.config.Display.Height,
		},
	}
}

// GetAnimationConfig converts core AnimationConfig to domain AnimationConfig
func (cm *ConfigManager) GetAnimationConfig() *domain.AnimationConfig {
	if cm.config == nil {
		return nil
	}

	return &domain.AnimationConfig{
		Speed:   cm.config.Animation.Speed,
		Colors:  cm.config.Animation.Colors,
		Enabled: cm.config.Animation.Enabled,
		Pattern: cm.config.Animation.Pattern,
	}
}

// UpdateConfig updates the configuration
func (cm *ConfigManager) UpdateConfig(updates map[string]interface{}) error {
	// Apply updates to specific fields
	for key, value := range updates {
		switch key {
		case "plugins.datasource":
			if v, ok := value.(string); ok {
				cm.config.Plugins.DataSource = v
			}
		}
	}
	return nil
}

// ValidateConfig validates the current configuration
func (cm *ConfigManager) ValidateConfig() error {
	if cm.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	// Display format validation is no longer needed since we removed formats

	// Validate animation pattern
	validPatterns := []domain.AnimationPattern{
		domain.PatternRainbow, domain.PatternGradient,
		domain.PatternPulse, domain.PatternWave,
	}

	patternValid := false
	for _, pattern := range validPatterns {
		if cm.config.Animation.Pattern == pattern {
			patternValid = true
			break
		}
	}
	if !patternValid {
		return fmt.Errorf("invalid animation pattern: %s", cm.config.Animation.Pattern)
	}

	// Validate display dimensions
	if cm.config.Display.Width <= 0 || cm.config.Display.Height <= 0 {
		return fmt.Errorf("display dimensions must be positive")
	}

	// Validate refresh rate
	if cm.config.App.RefreshRate <= 0 {
		return fmt.Errorf("refresh rate must be positive")
	}

	// Validate animation speed
	if cm.config.Animation.Speed <= 0 {
		return fmt.Errorf("animation speed must be positive")
	}

	return nil
}

// ApplyFlagsToConfig applies command line flag values to configuration
func (cm *ConfigManager) ApplyFlagsToConfig(flagConfig *FlagConfig) error {
	if cm.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	// Apply animation configuration from flags
	if flagConfig.Animation.Speed > 0 {
		cm.config.Animation.Speed = flagConfig.Animation.Speed
	}

	if flagConfig.Animation.Pattern != "" {
		cm.config.Animation.Pattern = flagConfig.Animation.Pattern
	}

	if flagConfig.Animation.Enabled != nil {
		cm.config.Animation.Enabled = *flagConfig.Animation.Enabled
	}

	// Apply bankruptcy mode (note: this affects datasource configuration)
	// Bankruptcy mode is handled by the main application, not by configuration

	return nil
}
