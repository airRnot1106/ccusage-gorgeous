package domain

import (
	"time"
)

// DisplayConfig represents the display configuration
type DisplayConfig struct {
	RefreshRate   time.Duration `json:"refresh_rate"`
	ShowTimestamp bool          `json:"show_timestamp"`
	ShowBreakdown bool          `json:"show_breakdown"`
	Format        DisplayFormat `json:"format"`
	Size          DisplaySize   `json:"size"`
}

// DisplayFormat defines the display format options
type DisplayFormat string

const (
	FormatLarge   DisplayFormat = "large"
	FormatMedium  DisplayFormat = "medium"
	FormatSmall   DisplayFormat = "small"
	FormatMinimal DisplayFormat = "minimal"
)

// DisplaySize defines the display size configuration
type DisplaySize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DisplayData represents the data to be displayed
type DisplayData struct {
	Cost        *CostData       `json:"cost"`
	Animation   *AnimationFrame `json:"animation"`
	Config      *DisplayConfig  `json:"config"`
	LastUpdated time.Time       `json:"last_updated"`
}

// DisplayService defines the interface for display operations
type DisplayService interface {
	Render(data *DisplayData) (string, error)
	GetDefaultConfig() *DisplayConfig
	ValidateConfig(config *DisplayConfig) error
}
