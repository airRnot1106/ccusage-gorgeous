package datasource_test

import (
	"context"
	"testing"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/stretchr/testify/assert"
)

func TestNewBankruptcyDataSourcePlugin(t *testing.T) {
	plugin := datasource.NewBankruptcyDataSourcePlugin()
	assert.NotNil(t, plugin)
	assert.Equal(t, "bankruptcy-datasource", plugin.Name())
	assert.Equal(t, "1.0.0", plugin.Version())
	assert.Equal(t, "Bankruptcy data source plugin that returns fixed $9999.99", plugin.Description())
	assert.False(t, plugin.IsEnabled()) // Should be disabled initially
}

func TestBankruptcyDataSourcePlugin_Initialize(t *testing.T) {
	plugin := datasource.NewBankruptcyDataSourcePlugin()

	// Test with empty config
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())
}

func TestBankruptcyDataSourcePlugin_Shutdown(t *testing.T) {
	plugin := datasource.NewBankruptcyDataSourcePlugin()

	// Initialize first
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())

	// Shutdown
	err = plugin.Shutdown()
	assert.NoError(t, err)
	assert.False(t, plugin.IsEnabled())
}

func TestBankruptcyDataSourcePlugin_FetchCostData_NotEnabled(t *testing.T) {
	plugin := datasource.NewBankruptcyDataSourcePlugin()
	ctx := context.Background()

	// Should fail when plugin is not enabled
	_, err := plugin.FetchCostData(ctx)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrPluginNotEnabled, err)
}

func TestBankruptcyDataSourcePlugin_FetchCostData_Success(t *testing.T) {
	plugin := datasource.NewBankruptcyDataSourcePlugin()
	ctx := context.Background()

	// Initialize plugin
	err := plugin.Initialize(map[string]interface{}{})
	assert.NoError(t, err)

	// Fetch cost data
	costData, err := plugin.FetchCostData(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, costData)

	// Verify bankruptcy cost data
	assert.Equal(t, 9999.99, costData.TotalCost)
	assert.Equal(t, "USD", costData.Currency)
	assert.NotEmpty(t, costData.Timestamp)
	assert.Equal(t, map[string]float64{"bankruptcy-mode": 9999.99}, costData.ModelBreakdown)
}
