package datasource_test

import (
	"context"
	"testing"

	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/stretchr/testify/assert"
)

func TestNewCcusageCliPlugin(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()
	assert.NotNil(t, plugin)
	assert.Equal(t, "ccusage-cli", plugin.Name())
	assert.Equal(t, "1.0.0", plugin.Version())
	assert.Equal(t, "ccusage CLI data source plugin", plugin.Description())
	assert.False(t, plugin.IsEnabled()) // Should be disabled initially
}

func TestCcusageCliPlugin_Initialize(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()

	// Test with empty config
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())

	// Test with custom config
	config := map[string]interface{}{
		"ccusage_path": "/custom/path/ccusage",
		"timeout":      "60s",
		"cache_time":   "30s",
	}

	err = plugin.Initialize(config)
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())
}

func TestCcusageCliPlugin_Shutdown(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()

	// Initialize first
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())

	// Shutdown
	err = plugin.Shutdown()
	assert.NoError(t, err)
	assert.False(t, plugin.IsEnabled())
}

func TestCcusageCliPlugin_FetchCostData_NotEnabled(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()
	ctx := context.Background()

	// Should fail when plugin is not enabled
	_, err := plugin.FetchCostData(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin is not enabled")
}

func TestCcusageCliPlugin_GetLastUpdated_NotEnabled(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()
	ctx := context.Background()

	// Should fail when plugin is not enabled
	_, err := plugin.GetLastUpdated(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin is not enabled")
}

func TestCcusageCliPlugin_GetLastUpdated_Enabled(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	// Should return zero time initially
	lastUpdated, err := plugin.GetLastUpdated(ctx)
	assert.NoError(t, err)
	assert.True(t, lastUpdated.IsZero())
}

func TestCcusageCliPlugin_SupportsRealtime(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()

	// ccusage CLI doesn't support real-time data
	assert.False(t, plugin.SupportsRealtime())
}

func TestCcusageCliPlugin_FetchCostData_CommandNotFound(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()
	ctx := context.Background()

	// Initialize with non-existent ccusage path
	config := map[string]interface{}{
		"ccusage_path": "/non/existent/ccusage",
		"timeout":      "1s", // Short timeout for faster test
	}

	err := plugin.Initialize(config)
	assert.NoError(t, err)

	// Should fail to execute command
	_, err = plugin.FetchCostData(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute ccusage command")
}

func TestCcusageCliPlugin_Initialize_InvalidTimeout(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()

	// Test with invalid timeout format
	config := map[string]interface{}{
		"timeout": "invalid-duration",
	}

	// Should still succeed, invalid values are ignored
	err := plugin.Initialize(config)
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())
}

func TestCcusageCliPlugin_Initialize_TypeConversion(t *testing.T) {
	plugin := datasource.NewCcusageCliPlugin()

	// Test with wrong types (should be ignored gracefully)
	config := map[string]interface{}{
		"ccusage_path": 12345, // should be string
		"timeout":      123,   // should be string
		"cache_time":   true,  // should be string
	}

	err := plugin.Initialize(config)
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())
}
