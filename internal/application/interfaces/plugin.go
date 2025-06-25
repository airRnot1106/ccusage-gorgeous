package interfaces

import (
	"context"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// Plugin defines the base interface for all plugins
type Plugin interface {
	Name() string
	Version() string
	Description() string
	Initialize(config map[string]interface{}) error
	Shutdown() error
	IsEnabled() bool
}

// DataSourcePlugin defines the interface for data source plugins
type DataSourcePlugin interface {
	Plugin
	FetchCostData(ctx context.Context) (*domain.CostData, error)
	GetLastUpdated(ctx context.Context) (time.Time, error)
	SupportsRealtime() bool
}

// DisplayPlugin defines the interface for display plugins
type DisplayPlugin interface {
	Plugin
	Render(ctx context.Context, data *domain.DisplayData) (string, error)
	GetCapabilities() DisplayCapabilities
	ValidateDisplayConfig(config *domain.DisplayConfig) error
}

// AnimationPlugin defines the interface for animation plugins
type AnimationPlugin interface {
	Plugin
	GenerateFrame(ctx context.Context, text string, frameNumber int, config *domain.AnimationConfig) (*domain.AnimationFrame, error)
	GetSupportedPatterns() []domain.AnimationPattern
	ValidateAnimationConfig(config *domain.AnimationConfig) error
}

// DisplayCapabilities represents the capabilities of a display plugin
type DisplayCapabilities struct {
	MaxWidth        int  `json:"max_width"`
	MaxHeight       int  `json:"max_height"`
	SupportsColor   bool `json:"supports_color"`
	SupportsUnicode bool `json:"supports_unicode"`
}

// PluginRegistry defines the interface for plugin management
type PluginRegistry interface {
	RegisterDataSource(plugin DataSourcePlugin) error
	RegisterDisplay(plugin DisplayPlugin) error
	RegisterAnimation(plugin AnimationPlugin) error
	GetDataSource(name string) (DataSourcePlugin, error)
	GetDisplay(name string) (DisplayPlugin, error)
	GetAnimation(name string) (AnimationPlugin, error)
	ListPlugins() []Plugin
	ShutdownAll() error
}
