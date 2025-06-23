package core

import (
	"fmt"
	"sync"

	"github.com/airRnot1106/ccusage-gorgeous/internal/application/interfaces"
)

// PluginRegistry implements the plugin registry interface
type PluginRegistry struct {
	mu            sync.RWMutex
	dataSources   map[string]interfaces.DataSourcePlugin
	displays      map[string]interfaces.DisplayPlugin
	animations    map[string]interfaces.AnimationPlugin
	configManager *ConfigManager
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry(configManager *ConfigManager) *PluginRegistry {
	return &PluginRegistry{
		dataSources:   make(map[string]interfaces.DataSourcePlugin),
		displays:      make(map[string]interfaces.DisplayPlugin),
		animations:    make(map[string]interfaces.AnimationPlugin),
		configManager: configManager,
	}
}

// RegisterDataSource registers a data source plugin
func (pr *PluginRegistry) RegisterDataSource(plugin interfaces.DataSourcePlugin) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	name := plugin.Name()
	if _, exists := pr.dataSources[name]; exists {
		return fmt.Errorf("data source plugin '%s' already registered", name)
	}

	pr.dataSources[name] = plugin
	return nil
}

// RegisterDisplay registers a display plugin
func (pr *PluginRegistry) RegisterDisplay(plugin interfaces.DisplayPlugin) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	name := plugin.Name()
	if _, exists := pr.displays[name]; exists {
		return fmt.Errorf("display plugin '%s' already registered", name)
	}

	pr.displays[name] = plugin
	return nil
}

// RegisterAnimation registers an animation plugin
func (pr *PluginRegistry) RegisterAnimation(plugin interfaces.AnimationPlugin) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	name := plugin.Name()
	if _, exists := pr.animations[name]; exists {
		return fmt.Errorf("animation plugin '%s' already registered", name)
	}

	pr.animations[name] = plugin
	return nil
}

// GetDataSource retrieves a data source plugin by name
func (pr *PluginRegistry) GetDataSource(name string) (interfaces.DataSourcePlugin, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	plugin, exists := pr.dataSources[name]
	if !exists {
		return nil, fmt.Errorf("data source plugin '%s' not found", name)
	}

	return plugin, nil
}

// GetDisplay retrieves a display plugin by name
func (pr *PluginRegistry) GetDisplay(name string) (interfaces.DisplayPlugin, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	plugin, exists := pr.displays[name]
	if !exists {
		return nil, fmt.Errorf("display plugin '%s' not found", name)
	}

	return plugin, nil
}

// GetAnimation retrieves an animation plugin by name
func (pr *PluginRegistry) GetAnimation(name string) (interfaces.AnimationPlugin, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	plugin, exists := pr.animations[name]
	if !exists {
		return nil, fmt.Errorf("animation plugin '%s' not found", name)
	}

	return plugin, nil
}

// ListPlugins returns all registered plugins
func (pr *PluginRegistry) ListPlugins() []interfaces.Plugin {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	var plugins []interfaces.Plugin

	for _, plugin := range pr.dataSources {
		plugins = append(plugins, plugin)
	}

	for _, plugin := range pr.displays {
		plugins = append(plugins, plugin)
	}

	for _, plugin := range pr.animations {
		plugins = append(plugins, plugin)
	}

	return plugins
}

// ShutdownAll shuts down all registered plugins
func (pr *PluginRegistry) ShutdownAll() error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	var errors []error

	// Shutdown data source plugins
	for name, plugin := range pr.dataSources {
		if err := plugin.Shutdown(); err != nil {
			errors = append(errors, fmt.Errorf("failed to shutdown data source plugin '%s': %w", name, err))
		}
	}

	// Shutdown display plugins
	for name, plugin := range pr.displays {
		if err := plugin.Shutdown(); err != nil {
			errors = append(errors, fmt.Errorf("failed to shutdown display plugin '%s': %w", name, err))
		}
	}

	// Shutdown animation plugins
	for name, plugin := range pr.animations {
		if err := plugin.Shutdown(); err != nil {
			errors = append(errors, fmt.Errorf("failed to shutdown animation plugin '%s': %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("plugin shutdown errors: %v", errors)
	}

	return nil
}

// GetActiveDataSource returns the active data source plugin based on config
func (pr *PluginRegistry) GetActiveDataSource() (interfaces.DataSourcePlugin, error) {
	config := pr.configManager.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("no configuration available")
	}

	return pr.GetDataSource(config.Plugins.DataSource)
}

// GetActiveDisplay returns the active display plugin based on config
func (pr *PluginRegistry) GetActiveDisplay() (interfaces.DisplayPlugin, error) {
	config := pr.configManager.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("no configuration available")
	}

	return pr.GetDisplay(config.Plugins.Display)
}

// GetActiveAnimation returns the active animation plugin based on config
func (pr *PluginRegistry) GetActiveAnimation() (interfaces.AnimationPlugin, error) {
	config := pr.configManager.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("no configuration available")
	}

	return pr.GetAnimation(config.Plugins.Animation)
}

// InitializePlugin initializes a plugin with its configuration
func (pr *PluginRegistry) InitializePlugin(plugin interfaces.Plugin) error {
	config := pr.configManager.GetConfig()
	if config == nil {
		return fmt.Errorf("no configuration available")
	}

	pluginConfig := config.Plugins.Config
	if pluginConfig == nil {
		pluginConfig = make(map[string]interface{})
	}

	return plugin.Initialize(pluginConfig)
}

// GetPluginCount returns the number of registered plugins by type
func (pr *PluginRegistry) GetPluginCount() (dataSources, displays, animations int) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	return len(pr.dataSources), len(pr.displays), len(pr.animations)
}
