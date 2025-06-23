package interfaces_test

import (
	"context"
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/application/interfaces"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing

type MockPlugin struct {
	name        string
	version     string
	description string
	enabled     bool
	initialized bool
	shouldError bool
}

func (m *MockPlugin) Name() string        { return m.name }
func (m *MockPlugin) Version() string     { return m.version }
func (m *MockPlugin) Description() string { return m.description }
func (m *MockPlugin) IsEnabled() bool     { return m.enabled }

func (m *MockPlugin) Initialize(config map[string]interface{}) error {
	if m.shouldError {
		return assert.AnError
	}
	m.initialized = true
	return nil
}

func (m *MockPlugin) Shutdown() error {
	if m.shouldError {
		return assert.AnError
	}
	m.enabled = false
	return nil
}

// MockDataSourcePlugin
type MockDataSourcePlugin struct {
	MockPlugin
	mockCostData    *domain.CostData
	mockLastUpdated time.Time
	realtime        bool
}

func (m *MockDataSourcePlugin) FetchCostData(ctx context.Context) (*domain.CostData, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return m.mockCostData, nil
}

func (m *MockDataSourcePlugin) GetLastUpdated(ctx context.Context) (time.Time, error) {
	if m.shouldError {
		return time.Time{}, assert.AnError
	}
	return m.mockLastUpdated, nil
}

func (m *MockDataSourcePlugin) SupportsRealtime() bool {
	return m.realtime
}

// MockDisplayPlugin
type MockDisplayPlugin struct {
	MockPlugin
	mockRender       string
	mockCapabilities interfaces.DisplayCapabilities
}

func (m *MockDisplayPlugin) Render(ctx context.Context, data *domain.DisplayData) (string, error) {
	if m.shouldError {
		return "", assert.AnError
	}
	return m.mockRender, nil
}

func (m *MockDisplayPlugin) GetCapabilities() interfaces.DisplayCapabilities {
	return m.mockCapabilities
}

func (m *MockDisplayPlugin) ValidateDisplayConfig(config *domain.DisplayConfig) error {
	if m.shouldError {
		return assert.AnError
	}
	return nil
}

// MockAnimationPlugin
type MockAnimationPlugin struct {
	MockPlugin
	mockFrame         *domain.AnimationFrame
	supportedPatterns []domain.AnimationPattern
}

func (m *MockAnimationPlugin) GenerateFrame(ctx context.Context, text string, frameNumber int, config *domain.AnimationConfig) (*domain.AnimationFrame, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	frame := &domain.AnimationFrame{
		Colors:    m.mockFrame.Colors,
		Text:      text,
		Timestamp: time.Now(),
	}
	return frame, nil
}

func (m *MockAnimationPlugin) GetSupportedPatterns() []domain.AnimationPattern {
	return m.supportedPatterns
}

func (m *MockAnimationPlugin) ValidateAnimationConfig(config *domain.AnimationConfig) error {
	if m.shouldError {
		return assert.AnError
	}
	return nil
}

// Test Plugin base interface
func TestPlugin_Interface(t *testing.T) {
	plugin := &MockPlugin{
		name:        "test-plugin",
		version:     "1.0.0",
		description: "Test plugin for unit testing",
		enabled:     true,
	}

	assert.Equal(t, "test-plugin", plugin.Name())
	assert.Equal(t, "1.0.0", plugin.Version())
	assert.Equal(t, "Test plugin for unit testing", plugin.Description())
	assert.True(t, plugin.IsEnabled())

	// Test initialization
	config := map[string]interface{}{
		"setting1": "value1",
		"setting2": 42,
	}

	err := plugin.Initialize(config)
	assert.NoError(t, err)
	assert.True(t, plugin.initialized)

	// Test shutdown
	err = plugin.Shutdown()
	assert.NoError(t, err)
	assert.False(t, plugin.enabled)

	// Test error cases
	plugin.shouldError = true
	err = plugin.Initialize(config)
	assert.Error(t, err)

	plugin.enabled = true
	err = plugin.Shutdown()
	assert.Error(t, err)
}

// Test DataSourcePlugin interface
func TestDataSourcePlugin_Interface(t *testing.T) {
	now := time.Now()
	mockCostData := &domain.CostData{
		TotalCost: 25.75,
		Currency:  "USD",
		Timestamp: now,
	}

	plugin := &MockDataSourcePlugin{
		MockPlugin: MockPlugin{
			name:        "ccusage-datasource",
			version:     "1.0.0",
			description: "ccusage CLI data source plugin",
			enabled:     true,
		},
		mockCostData:    mockCostData,
		mockLastUpdated: now,
		realtime:        true,
	}

	ctx := context.Background()

	// Test FetchCostData
	costData, err := plugin.FetchCostData(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 25.75, costData.TotalCost)
	assert.Equal(t, "USD", costData.Currency)

	// Test GetLastUpdated
	lastUpdated, err := plugin.GetLastUpdated(ctx)
	assert.NoError(t, err)
	assert.Equal(t, now, lastUpdated)

	// Test SupportsRealtime
	assert.True(t, plugin.SupportsRealtime())

	// Test error cases
	plugin.shouldError = true
	_, err = plugin.FetchCostData(ctx)
	assert.Error(t, err)

	_, err = plugin.GetLastUpdated(ctx)
	assert.Error(t, err)
}

// Test DisplayPlugin interface
func TestDisplayPlugin_Interface(t *testing.T) {
	capabilities := interfaces.DisplayCapabilities{
		SupportedFormats: []domain.DisplayFormat{domain.FormatLarge, domain.FormatMedium},
		MaxWidth:         120,
		MaxHeight:        40,
		SupportsColor:    true,
		SupportsUnicode:  true,
	}

	plugin := &MockDisplayPlugin{
		MockPlugin: MockPlugin{
			name:        "rainbow-display",
			version:     "1.0.0",
			description: "Rainbow animation display plugin",
			enabled:     true,
		},
		mockRender:       "Rainbow Animated Display",
		mockCapabilities: capabilities,
	}

	ctx := context.Background()

	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 30.25,
			Currency:  "USD",
			Timestamp: time.Now(),
		},
		LastUpdated: time.Now(),
	}

	// Test Render
	rendered, err := plugin.Render(ctx, displayData)
	assert.NoError(t, err)
	assert.Equal(t, "Rainbow Animated Display", rendered)

	// Test GetCapabilities
	caps := plugin.GetCapabilities()
	assert.Len(t, caps.SupportedFormats, 2)
	assert.Equal(t, 120, caps.MaxWidth)
	assert.Equal(t, 40, caps.MaxHeight)
	assert.True(t, caps.SupportsColor)
	assert.True(t, caps.SupportsUnicode)

	// Test ValidateDisplayConfig
	config := &domain.DisplayConfig{
		Format: domain.FormatLarge,
		Size:   domain.DisplaySize{Width: 80, Height: 24},
	}
	err = plugin.ValidateDisplayConfig(config)
	assert.NoError(t, err)

	// Test error cases
	plugin.shouldError = true
	_, err = plugin.Render(ctx, displayData)
	assert.Error(t, err)

	err = plugin.ValidateDisplayConfig(config)
	assert.Error(t, err)
}

// Test AnimationPlugin interface
func TestAnimationPlugin_Interface(t *testing.T) {
	mockFrame := &domain.AnimationFrame{
		Colors: []string{"#FF0000", "#00FF00", "#0000FF"},
	}

	supportedPatterns := []domain.AnimationPattern{
		domain.PatternRainbow,
		domain.PatternGradient,
		domain.PatternWave,
	}

	plugin := &MockAnimationPlugin{
		MockPlugin: MockPlugin{
			name:        "rainbow-animator",
			version:     "1.0.0",
			description: "Rainbow animation plugin",
			enabled:     true,
		},
		mockFrame:         mockFrame,
		supportedPatterns: supportedPatterns,
	}

	ctx := context.Background()

	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Pattern: domain.PatternRainbow,
		Enabled: true,
	}

	// Test GenerateFrame
	frame, err := plugin.GenerateFrame(ctx, "$35.50", 10, config)
	assert.NoError(t, err)
	assert.Equal(t, "$35.50", frame.Text)
	assert.Len(t, frame.Colors, 3)

	// Test GetSupportedPatterns
	patterns := plugin.GetSupportedPatterns()
	assert.Len(t, patterns, 3)
	assert.Contains(t, patterns, domain.PatternRainbow)
	assert.Contains(t, patterns, domain.PatternGradient)
	assert.Contains(t, patterns, domain.PatternWave)

	// Test ValidateAnimationConfig
	err = plugin.ValidateAnimationConfig(config)
	assert.NoError(t, err)

	// Test error cases
	plugin.shouldError = true
	_, err = plugin.GenerateFrame(ctx, "$35.50", 10, config)
	assert.Error(t, err)

	err = plugin.ValidateAnimationConfig(config)
	assert.Error(t, err)
}

// Test DisplayCapabilities
func TestDisplayCapabilities(t *testing.T) {
	capabilities := interfaces.DisplayCapabilities{
		SupportedFormats: []domain.DisplayFormat{
			domain.FormatLarge,
			domain.FormatMedium,
			domain.FormatSmall,
		},
		MaxWidth:        200,
		MaxHeight:       50,
		SupportsColor:   true,
		SupportsUnicode: false,
	}

	assert.Len(t, capabilities.SupportedFormats, 3)
	assert.Equal(t, 200, capabilities.MaxWidth)
	assert.Equal(t, 50, capabilities.MaxHeight)
	assert.True(t, capabilities.SupportsColor)
	assert.False(t, capabilities.SupportsUnicode)
}
