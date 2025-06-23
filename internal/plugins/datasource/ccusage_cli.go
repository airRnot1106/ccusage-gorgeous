package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// CcusageCliPlugin implements the DataSourcePlugin interface for ccusage CLI
type CcusageCliPlugin struct {
	name        string
	version     string
	description string
	enabled     bool
	ccusagePath string
	timeout     time.Duration
	cacheTime   time.Duration
	lastUpdate  time.Time
	cachedData  *domain.CostData
}

// CcusageResponse represents the JSON response from ccusage CLI
type CcusageResponse struct {
	TotalCost      float64            `json:"total_cost"`
	Currency       string             `json:"currency"`
	Date           string             `json:"date"`
	ModelBreakdown map[string]float64 `json:"model_breakdown,omitempty"`
}

// NewCcusageCliPlugin creates a new ccusage CLI plugin
func NewCcusageCliPlugin() *CcusageCliPlugin {
	return &CcusageCliPlugin{
		name:        "ccusage-cli",
		version:     "1.0.0",
		description: "ccusage CLI data source plugin",
		enabled:     false,
		ccusagePath: "ccusage",
		timeout:     30 * time.Second,
		cacheTime:   10 * time.Second,
	}
}

// Name returns the plugin name
func (c *CcusageCliPlugin) Name() string {
	return c.name
}

// Version returns the plugin version
func (c *CcusageCliPlugin) Version() string {
	return c.version
}

// Description returns the plugin description
func (c *CcusageCliPlugin) Description() string {
	return c.description
}

// IsEnabled returns whether the plugin is enabled
func (c *CcusageCliPlugin) IsEnabled() bool {
	return c.enabled
}

// Initialize initializes the plugin with configuration
func (c *CcusageCliPlugin) Initialize(config map[string]interface{}) error {
	if ccusagePath, ok := config["ccusage_path"].(string); ok {
		c.ccusagePath = ccusagePath
	}

	if timeout, ok := config["timeout"].(string); ok {
		if duration, err := time.ParseDuration(timeout); err == nil {
			c.timeout = duration
		}
	}

	if cacheTime, ok := config["cache_time"].(string); ok {
		if duration, err := time.ParseDuration(cacheTime); err == nil {
			c.cacheTime = duration
		}
	}

	c.enabled = true
	return nil
}

// Shutdown shuts down the plugin
func (c *CcusageCliPlugin) Shutdown() error {
	c.enabled = false
	c.cachedData = nil
	return nil
}

// FetchCostData fetches cost data from ccusage CLI
func (c *CcusageCliPlugin) FetchCostData(ctx context.Context) (*domain.CostData, error) {
	if !c.enabled {
		return nil, fmt.Errorf("plugin is not enabled")
	}

	// Check cache first
	if c.cachedData != nil && time.Since(c.lastUpdate) < c.cacheTime {
		return c.cachedData, nil
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Execute ccusage command with JSON output via npx
	var cmd *exec.Cmd
	if c.ccusagePath == "ccusage" {
		// Use npx for default ccusage command
		cmd = exec.CommandContext(timeoutCtx, "npx", "ccusage", "daily", "--json")
	} else {
		// Use custom path as specified
		cmd = exec.CommandContext(timeoutCtx, c.ccusagePath, "daily", "--json")
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute ccusage command: %w", err)
	}

	// Parse JSON response
	var response CcusageResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse ccusage JSON output: %w", err)
	}

	// Parse date
	var timestamp time.Time
	if response.Date != "" {
		if parsedTime, err := time.Parse("2006-01-02", response.Date); err == nil {
			timestamp = parsedTime
		} else {
			timestamp = time.Now()
		}
	} else {
		timestamp = time.Now()
	}

	// Convert to domain model
	costData := &domain.CostData{
		TotalCost:      response.TotalCost,
		Currency:       response.Currency,
		Timestamp:      timestamp,
		ModelBreakdown: response.ModelBreakdown,
	}

	// Update cache
	c.cachedData = costData
	c.lastUpdate = time.Now()

	return costData, nil
}

// GetLastUpdated returns the timestamp of the last data update
func (c *CcusageCliPlugin) GetLastUpdated(ctx context.Context) (time.Time, error) {
	if !c.enabled {
		return time.Time{}, fmt.Errorf("plugin is not enabled")
	}

	return c.lastUpdate, nil
}

// SupportsRealtime returns whether the plugin supports real-time data
func (c *CcusageCliPlugin) SupportsRealtime() bool {
	return false // ccusage CLI doesn't support real-time streaming
}
