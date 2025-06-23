package datasource

import (
	"context"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// BankruptcyDataSourcePlugin implements a data source that always returns bankruptcy cost
type BankruptcyDataSourcePlugin struct {
	name        string
	version     string
	description string
	enabled     bool
}

// NewBankruptcyDataSourcePlugin creates a new bankruptcy data source plugin
func NewBankruptcyDataSourcePlugin() *BankruptcyDataSourcePlugin {
	return &BankruptcyDataSourcePlugin{
		name:        "bankruptcy-datasource",
		version:     "1.0.0",
		description: "Bankruptcy data source plugin that returns fixed $9999.99",
		enabled:     false,
	}
}

// Name returns the plugin name
func (b *BankruptcyDataSourcePlugin) Name() string {
	return b.name
}

// Version returns the plugin version
func (b *BankruptcyDataSourcePlugin) Version() string {
	return b.version
}

// Description returns the plugin description
func (b *BankruptcyDataSourcePlugin) Description() string {
	return b.description
}

// IsEnabled returns whether the plugin is enabled
func (b *BankruptcyDataSourcePlugin) IsEnabled() bool {
	return b.enabled
}

// Initialize initializes the plugin with configuration
func (b *BankruptcyDataSourcePlugin) Initialize(config map[string]interface{}) error {
	b.enabled = true
	return nil
}

// Shutdown shuts down the plugin
func (b *BankruptcyDataSourcePlugin) Shutdown() error {
	b.enabled = false
	return nil
}

// FetchCostData returns bankruptcy cost data ($9999.99)
func (b *BankruptcyDataSourcePlugin) FetchCostData(ctx context.Context) (*domain.CostData, error) {
	if !b.enabled {
		return nil, domain.ErrPluginNotEnabled
	}

	return &domain.CostData{
		TotalCost: 9999.99,
		Currency:  "USD",
		Timestamp: time.Now(),
		ModelBreakdown: map[string]float64{
			"bankruptcy-mode": 9999.99,
		},
	}, nil
}

// GetLastUpdated returns the current time (bankruptcy data is always "fresh")
func (b *BankruptcyDataSourcePlugin) GetLastUpdated(ctx context.Context) (time.Time, error) {
	if !b.enabled {
		return time.Time{}, domain.ErrPluginNotEnabled
	}

	return time.Now(), nil
}

// SupportsRealtime returns false as bankruptcy data is static
func (b *BankruptcyDataSourcePlugin) SupportsRealtime() bool {
	return false
}
