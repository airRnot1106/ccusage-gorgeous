package core_test

import (
	"testing"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/animation"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/datasource"
	"github.com/airRnot1106/ccusage-gorgeous/internal/plugins/display"
	"github.com/stretchr/testify/assert"
)

func TestNewPluginRegistry(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	assert.NotNil(t, registry)

	// Should start empty
	dataSourceCount, displayCount, animationCount := registry.GetPluginCount()
	assert.Equal(t, 0, dataSourceCount)
	assert.Equal(t, 0, displayCount)
	assert.Equal(t, 0, animationCount)
}

func TestPluginRegistry_RegisterDataSource(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := datasource.NewCcusageCliPlugin()

	// Register plugin
	err = registry.RegisterDataSource(plugin)
	assert.NoError(t, err)

	// Check count
	dataSourceCount, _, _ := registry.GetPluginCount()
	assert.Equal(t, 1, dataSourceCount)

	// Try to register the same plugin again (should fail)
	err = registry.RegisterDataSource(plugin)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestPluginRegistry_RegisterDisplay(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := display.NewRainbowTUIPlugin()

	// Register plugin
	err = registry.RegisterDisplay(plugin)
	assert.NoError(t, err)

	// Check count
	_, displayCount, _ := registry.GetPluginCount()
	assert.Equal(t, 1, displayCount)

	// Try to register the same plugin again (should fail)
	err = registry.RegisterDisplay(plugin)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestPluginRegistry_RegisterAnimation(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := animation.NewRainbowAnimationPlugin()

	// Register plugin
	err = registry.RegisterAnimation(plugin)
	assert.NoError(t, err)

	// Check count
	_, _, animationCount := registry.GetPluginCount()
	assert.Equal(t, 1, animationCount)

	// Try to register the same plugin again (should fail)
	err = registry.RegisterAnimation(plugin)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestPluginRegistry_GetDataSource(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := datasource.NewCcusageCliPlugin()

	// Should fail when plugin not registered
	_, err = registry.GetDataSource("ccusage-cli")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Register and try again
	err = registry.RegisterDataSource(plugin)
	assert.NoError(t, err)

	retrievedPlugin, err := registry.GetDataSource("ccusage-cli")
	assert.NoError(t, err)
	assert.Equal(t, plugin, retrievedPlugin)
}

func TestPluginRegistry_GetDisplay(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := display.NewRainbowTUIPlugin()

	// Should fail when plugin not registered
	_, err = registry.GetDisplay("rainbow-display")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Register and try again
	err = registry.RegisterDisplay(plugin)
	assert.NoError(t, err)

	retrievedPlugin, err := registry.GetDisplay("rainbow-display")
	assert.NoError(t, err)
	assert.Equal(t, plugin, retrievedPlugin)
}

func TestPluginRegistry_GetAnimation(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := animation.NewRainbowAnimationPlugin()

	// Should fail when plugin not registered
	_, err = registry.GetAnimation("rainbow-animation")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Register and try again
	err = registry.RegisterAnimation(plugin)
	assert.NoError(t, err)

	retrievedPlugin, err := registry.GetAnimation("rainbow-animation")
	assert.NoError(t, err)
	assert.Equal(t, plugin, retrievedPlugin)
}

func TestPluginRegistry_ListPlugins(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)

	// Should start empty
	plugins := registry.ListPlugins()
	assert.Len(t, plugins, 0)

	// Register plugins
	dsPlugin := datasource.NewCcusageCliPlugin()
	dispPlugin := display.NewRainbowTUIPlugin()
	animPlugin := animation.NewRainbowAnimationPlugin()

	err = registry.RegisterDataSource(dsPlugin)
	assert.NoError(t, err)
	err = registry.RegisterDisplay(dispPlugin)
	assert.NoError(t, err)
	err = registry.RegisterAnimation(animPlugin)
	assert.NoError(t, err)

	// Should have 3 plugins
	plugins = registry.ListPlugins()
	assert.Len(t, plugins, 3)
}

func TestPluginRegistry_ShutdownAll(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)

	// Register and initialize plugins
	dsPlugin := datasource.NewCcusageCliPlugin()
	dispPlugin := display.NewRainbowTUIPlugin()
	animPlugin := animation.NewRainbowAnimationPlugin()

	err = registry.RegisterDataSource(dsPlugin)
	assert.NoError(t, err)
	err = registry.RegisterDisplay(dispPlugin)
	assert.NoError(t, err)
	err = registry.RegisterAnimation(animPlugin)
	assert.NoError(t, err)

	// Initialize plugins
	err = registry.InitializePlugin(dsPlugin)
	assert.NoError(t, err)
	err = registry.InitializePlugin(dispPlugin)
	assert.NoError(t, err)
	err = registry.InitializePlugin(animPlugin)
	assert.NoError(t, err)

	// All should be enabled
	assert.True(t, dsPlugin.IsEnabled())
	assert.True(t, dispPlugin.IsEnabled())
	assert.True(t, animPlugin.IsEnabled())

	// Shutdown all
	err = registry.ShutdownAll()
	assert.NoError(t, err)

	// All should be disabled
	assert.False(t, dsPlugin.IsEnabled())
	assert.False(t, dispPlugin.IsEnabled())
	assert.False(t, animPlugin.IsEnabled())
}

func TestPluginRegistry_GetActivePlugins(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)

	// Should fail when no plugins registered
	_, err = registry.GetActiveDataSource()
	assert.Error(t, err)
	_, err = registry.GetActiveDisplay()
	assert.Error(t, err)
	_, err = registry.GetActiveAnimation()
	assert.Error(t, err)

	// Register plugins with names matching config defaults
	dsPlugin := datasource.NewCcusageCliPlugin()
	dispPlugin := display.NewRainbowTUIPlugin()
	animPlugin := animation.NewRainbowAnimationPlugin()

	err = registry.RegisterDataSource(dsPlugin)
	assert.NoError(t, err)
	err = registry.RegisterDisplay(dispPlugin)
	assert.NoError(t, err)
	err = registry.RegisterAnimation(animPlugin)
	assert.NoError(t, err)

	// Should now work
	activeDS, err := registry.GetActiveDataSource()
	assert.NoError(t, err)
	assert.Equal(t, dsPlugin, activeDS)

	activeDisp, err := registry.GetActiveDisplay()
	assert.NoError(t, err)
	assert.Equal(t, dispPlugin, activeDisp)

	activeAnim, err := registry.GetActiveAnimation()
	assert.NoError(t, err)
	assert.Equal(t, animPlugin, activeAnim)
}

func TestPluginRegistry_InitializePlugin(t *testing.T) {
	configManager := core.NewConfigManager()
	err := configManager.LoadConfig("")
	assert.NoError(t, err)

	registry := core.NewPluginRegistry(configManager)
	plugin := datasource.NewCcusageCliPlugin()

	// Should start disabled
	assert.False(t, plugin.IsEnabled())

	// Initialize
	err = registry.InitializePlugin(plugin)
	assert.NoError(t, err)

	// Should be enabled
	assert.True(t, plugin.IsEnabled())
}
