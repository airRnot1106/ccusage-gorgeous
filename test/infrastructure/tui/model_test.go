package tui_test

import (
	"context"
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/airRnot1106/ccusage-gorgeous/internal/infrastructure/tui"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/animation"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/display"
	"github.com/stretchr/testify/assert"
)

func setupTestModel(t *testing.T) (*tui.Model, *core.PluginRegistry) {
	ctx := context.Background()

	// Initialize configuration manager
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Initialize plugin registry
	registry := core.NewPluginRegistry(configManager)

	// Register plugins
	ccusagePlugin := datasource.NewCcusageCliPlugin()
	rainbowAnimationPlugin := animation.NewRainbowAnimationPlugin()
	rainbowDisplayPlugin := display.NewRainbowTUIPlugin()

	err = registry.RegisterDataSource(ccusagePlugin)
	assert.NoError(t, err)

	err = registry.RegisterAnimation(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.RegisterDisplay(rainbowDisplayPlugin)
	assert.NoError(t, err)

	// Initialize plugins
	err = registry.InitializePlugin(ccusagePlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowDisplayPlugin)
	assert.NoError(t, err)

	// Create model
	model := tui.NewModel(ctx, registry, configManager)

	return model, registry
}

func setupBankruptcyTestModel(t *testing.T) (*tui.Model, *core.PluginRegistry) {
	ctx := context.Background()

	// Initialize configuration manager
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	// Update configuration to use bankruptcy data source
	err = configManager.UpdateConfig(map[string]interface{}{
		"plugins.datasource": "bankruptcy-datasource",
	})
	assert.NoError(t, err)

	// Initialize plugin registry
	registry := core.NewPluginRegistry(configManager)

	// Register bankruptcy plugin
	bankruptcyPlugin := datasource.NewBankruptcyDataSourcePlugin()
	rainbowAnimationPlugin := animation.NewRainbowAnimationPlugin()
	rainbowDisplayPlugin := display.NewRainbowTUIPlugin()

	err = registry.RegisterDataSource(bankruptcyPlugin)
	assert.NoError(t, err)

	err = registry.RegisterAnimation(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.RegisterDisplay(rainbowDisplayPlugin)
	assert.NoError(t, err)

	// Initialize plugins
	err = registry.InitializePlugin(bankruptcyPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowAnimationPlugin)
	assert.NoError(t, err)

	err = registry.InitializePlugin(rainbowDisplayPlugin)
	assert.NoError(t, err)

	// Create model
	model := tui.NewModel(ctx, registry, configManager)

	return model, registry
}

func TestModel_NewModel(t *testing.T) {
	ctx := context.Background()
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)

	model := tui.NewModel(ctx, registry, configManager)

	assert.NotNil(t, model)
}

func TestModel_BasicFunctionality(t *testing.T) {
	model, _ := setupTestModel(t)

	// Test that model is created properly
	assert.NotNil(t, model)

	// Create mock cost data
	mockCostData := &domain.CostData{
		TotalCost: 42.50,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	// Note: This is a basic test. In a real scenario, we would need to
	// inject the cost data and test the rendering output
	_ = mockCostData
}

func TestModel_WithBankruptcyDataSource(t *testing.T) {
	model, registry := setupBankruptcyTestModel(t)

	// Test that model is created properly with bankruptcy data source
	assert.NotNil(t, model)

	// Verify that bankruptcy data source is active
	activeDataSource, err := registry.GetActiveDataSource()
	assert.NoError(t, err)
	assert.Equal(t, "bankruptcy-datasource", activeDataSource.Name())
}
