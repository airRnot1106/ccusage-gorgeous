package domain_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAnimationConfig_Creation(t *testing.T) {
	config := &domain.AnimationConfig{
		Speed:   100 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	assert.Equal(t, 100*time.Millisecond, config.Speed)
	assert.Len(t, config.Colors, 3)
	assert.True(t, config.Enabled)
	assert.Equal(t, domain.PatternRainbow, config.Pattern)
}

func TestAnimationPattern_Constants(t *testing.T) {
	assert.Equal(t, domain.AnimationPattern("rainbow"), domain.PatternRainbow)
	assert.Equal(t, domain.AnimationPattern("gradient"), domain.PatternGradient)
	assert.Equal(t, domain.AnimationPattern("pulse"), domain.PatternPulse)
	assert.Equal(t, domain.AnimationPattern("wave"), domain.PatternWave)
}

func TestAnimationFrame_Creation(t *testing.T) {
	now := time.Now()
	frame := &domain.AnimationFrame{
		Colors:    []string{"#FF0000", "#FF8000"},
		Text:      "$15.50",
		Timestamp: now,
	}

	assert.Len(t, frame.Colors, 2)
	assert.Equal(t, "$15.50", frame.Text)
	assert.Equal(t, now, frame.Timestamp)
}

// Mock AnimationService for testing
type MockAnimationService struct {
	mockFrame       *domain.AnimationFrame
	mockConfig      *domain.AnimationConfig
	shouldError     bool
	validationError bool
}

func (m *MockAnimationService) GenerateFrame(text string, frameNumber int) (*domain.AnimationFrame, error) {
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

func (m *MockAnimationService) GetDefaultConfig() *domain.AnimationConfig {
	if m.mockConfig == nil {
		return &domain.AnimationConfig{
			Speed:   100 * time.Millisecond,
			Colors:  []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00", "#FF00FF", "#00FFFF"},
			Enabled: true,
			Pattern: domain.PatternRainbow,
		}
	}
	return m.mockConfig
}

func (m *MockAnimationService) ValidateConfig(config *domain.AnimationConfig) error {
	if m.validationError {
		return assert.AnError
	}
	return nil
}

func TestAnimationService_Interface(t *testing.T) {
	mockFrame := &domain.AnimationFrame{
		Colors:    []string{"#FF0000", "#00FF00"},
		Text:      "",
		Timestamp: time.Now(),
	}

	service := &MockAnimationService{
		mockFrame:   mockFrame,
		shouldError: false,
	}

	// Test GenerateFrame
	frame, err := service.GenerateFrame("$25.75", 5)
	assert.NoError(t, err)
	assert.Equal(t, "$25.75", frame.Text)
	assert.Len(t, frame.Colors, 2)

	// Test GetDefaultConfig
	config := service.GetDefaultConfig()
	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, domain.PatternRainbow, config.Pattern)
	assert.Len(t, config.Colors, 6)

	// Test ValidateConfig
	err = service.ValidateConfig(config)
	assert.NoError(t, err)

	// Test error cases
	service.shouldError = true
	_, err = service.GenerateFrame("$25.75", 5)
	assert.Error(t, err)

	service.shouldError = false
	service.validationError = true
	err = service.ValidateConfig(config)
	assert.Error(t, err)
}

func TestAnimationConfig_Validation(t *testing.T) {
	service := &MockAnimationService{}

	// Valid config
	validConfig := &domain.AnimationConfig{
		Speed:   50 * time.Millisecond,
		Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
		Enabled: true,
		Pattern: domain.PatternRainbow,
	}

	err := service.ValidateConfig(validConfig)
	assert.NoError(t, err)

	// Test with invalid config (simulated by setting validationError)
	service.validationError = true
	err = service.ValidateConfig(validConfig)
	assert.Error(t, err)
}

func TestAnimationFrame_WithDifferentPatterns(t *testing.T) {
	patterns := []domain.AnimationPattern{
		domain.PatternRainbow,
		domain.PatternGradient,
		domain.PatternPulse,
		domain.PatternWave,
	}

	for _, pattern := range patterns {
		config := &domain.AnimationConfig{
			Speed:   100 * time.Millisecond,
			Colors:  []string{"#FF0000", "#00FF00", "#0000FF"},
			Enabled: true,
			Pattern: pattern,
		}

		assert.Equal(t, pattern, config.Pattern)
	}
}
