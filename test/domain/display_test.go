package domain_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDisplayConfig_Creation(t *testing.T) {
	config := &domain.DisplayConfig{
		RefreshRate: 1 * time.Second,
		Size: domain.DisplaySize{
			Width:  80,
			Height: 24,
		},
	}

	assert.Equal(t, 1*time.Second, config.RefreshRate)
	assert.Equal(t, 80, config.Size.Width)
	assert.Equal(t, 24, config.Size.Height)
}

// TestDisplayFormat_Constants removed - DisplayFormat enum no longer exists

func TestDisplaySize_Creation(t *testing.T) {
	size := domain.DisplaySize{
		Width:  120,
		Height: 30,
	}

	assert.Equal(t, 120, size.Width)
	assert.Equal(t, 30, size.Height)
}

func TestDisplayData_Creation(t *testing.T) {
	now := time.Now()

	costData := &domain.CostData{
		TotalCost: 15.50,
		Currency:  "USD",
		Timestamp: now,
	}

	animationFrame := &domain.AnimationFrame{
		Colors:    []string{"#FF0000", "#00FF00"},
		Text:      "$15.50",
		Timestamp: now,
	}

	config := &domain.DisplayConfig{
		RefreshRate: 500 * time.Millisecond,
		Size:        domain.DisplaySize{Width: 80, Height: 24},
	}

	displayData := &domain.DisplayData{
		Cost:        costData,
		Animation:   animationFrame,
		Config:      config,
		LastUpdated: now,
	}

	assert.Equal(t, costData, displayData.Cost)
	assert.Equal(t, animationFrame, displayData.Animation)
	assert.Equal(t, config, displayData.Config)
	assert.Equal(t, now, displayData.LastUpdated)
}

// Mock DisplayService for testing
type MockDisplayService struct {
	mockRender      string
	mockConfig      *domain.DisplayConfig
	shouldError     bool
	validationError bool
}

func (m *MockDisplayService) Render(data *domain.DisplayData) (string, error) {
	if m.shouldError {
		return "", assert.AnError
	}
	return m.mockRender, nil
}

func (m *MockDisplayService) GetDefaultConfig() *domain.DisplayConfig {
	if m.mockConfig == nil {
		return &domain.DisplayConfig{
			RefreshRate: 1 * time.Second,
			Size: domain.DisplaySize{
				Width:  80,
				Height: 24,
			},
		}
	}
	return m.mockConfig
}

func (m *MockDisplayService) ValidateConfig(config *domain.DisplayConfig) error {
	if m.validationError {
		return assert.AnError
	}
	return nil
}

func TestDisplayService_Interface(t *testing.T) {
	service := &MockDisplayService{
		mockRender:  "Rendered Output",
		shouldError: false,
	}

	displayData := &domain.DisplayData{
		Cost: &domain.CostData{
			TotalCost: 25.75,
			Currency:  "USD",
			Timestamp: time.Now(),
		},
		Animation: &domain.AnimationFrame{
			Colors:    []string{"#FF0000", "#00FF00"},
			Text:      "$25.75",
			Timestamp: time.Now(),
		},
		LastUpdated: time.Now(),
	}

	// Test Render
	rendered, err := service.Render(displayData)
	assert.NoError(t, err)
	assert.Equal(t, "Rendered Output", rendered)

	// Test GetDefaultConfig
	config := service.GetDefaultConfig()
	assert.NotNil(t, config)
	assert.Equal(t, 1*time.Second, config.RefreshRate)
	assert.Equal(t, 80, config.Size.Width)
	assert.Equal(t, 24, config.Size.Height)

	// Test ValidateConfig
	err = service.ValidateConfig(config)
	assert.NoError(t, err)

	// Test error cases
	service.shouldError = true
	_, err = service.Render(displayData)
	assert.Error(t, err)

	service.shouldError = false
	service.validationError = true
	err = service.ValidateConfig(config)
	assert.Error(t, err)
}

func TestDisplayData_WithNilValues(t *testing.T) {
	displayData := &domain.DisplayData{
		LastUpdated: time.Now(),
	}

	assert.Nil(t, displayData.Cost)
	assert.Nil(t, displayData.Animation)
	assert.Nil(t, displayData.Config)
	assert.False(t, displayData.LastUpdated.IsZero())
}

// TestDisplayConfig_WithDifferentFormats removed - DisplayFormat enum no longer exists

func TestDisplaySize_EdgeCases(t *testing.T) {
	// Test zero size
	zeroSize := domain.DisplaySize{Width: 0, Height: 0}
	assert.Equal(t, 0, zeroSize.Width)
	assert.Equal(t, 0, zeroSize.Height)

	// Test large size
	largeSize := domain.DisplaySize{Width: 1920, Height: 1080}
	assert.Equal(t, 1920, largeSize.Width)
	assert.Equal(t, 1080, largeSize.Height)
}
