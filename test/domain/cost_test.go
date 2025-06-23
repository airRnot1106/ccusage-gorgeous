package domain_test

import (
	"testing"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCostData_Creation(t *testing.T) {
	now := time.Now()
	costData := &domain.CostData{
		TotalCost: 15.50,
		Currency:  "USD",
		Timestamp: now,
		ModelBreakdown: map[string]float64{
			"claude-3-opus":   10.00,
			"claude-3-sonnet": 5.50,
		},
	}

	assert.Equal(t, 15.50, costData.TotalCost)
	assert.Equal(t, "USD", costData.Currency)
	assert.Equal(t, now, costData.Timestamp)
	assert.Len(t, costData.ModelBreakdown, 2)
	assert.Equal(t, 10.00, costData.ModelBreakdown["claude-3-opus"])
	assert.Equal(t, 5.50, costData.ModelBreakdown["claude-3-sonnet"])
}

func TestCostData_EmptyModelBreakdown(t *testing.T) {
	costData := &domain.CostData{
		TotalCost: 15.50,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	assert.Equal(t, 15.50, costData.TotalCost)
	assert.Nil(t, costData.ModelBreakdown)
}

// Mock implementations for testing interfaces
type MockCostDataRepository struct {
	mockCostData    *domain.CostData
	mockLastUpdated time.Time
	shouldError     bool
}

func (m *MockCostDataRepository) FetchCostData() (*domain.CostData, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return m.mockCostData, nil
}

func (m *MockCostDataRepository) GetLastUpdated() (time.Time, error) {
	if m.shouldError {
		return time.Time{}, assert.AnError
	}
	return m.mockLastUpdated, nil
}

func TestCostDataRepository_Interface(t *testing.T) {
	now := time.Now()
	mockCostData := &domain.CostData{
		TotalCost: 25.75,
		Currency:  "USD",
		Timestamp: now,
	}

	repo := &MockCostDataRepository{
		mockCostData:    mockCostData,
		mockLastUpdated: now,
		shouldError:     false,
	}

	// Test successful fetch
	costData, err := repo.FetchCostData()
	assert.NoError(t, err)
	assert.Equal(t, 25.75, costData.TotalCost)
	assert.Equal(t, "USD", costData.Currency)

	// Test successful GetLastUpdated
	lastUpdated, err := repo.GetLastUpdated()
	assert.NoError(t, err)
	assert.Equal(t, now, lastUpdated)

	// Test error cases
	repo.shouldError = true

	_, err = repo.FetchCostData()
	assert.Error(t, err)

	_, err = repo.GetLastUpdated()
	assert.Error(t, err)
}

type MockCostDataService struct {
	mockCostData *domain.CostData
	mockHistory  []*domain.CostData
	shouldError  bool
}

func (m *MockCostDataService) GetCurrentCost() (*domain.CostData, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return m.mockCostData, nil
}

func (m *MockCostDataService) GetCostHistory(days int) ([]*domain.CostData, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return m.mockHistory, nil
}

func (m *MockCostDataService) RefreshCostData() error {
	if m.shouldError {
		return assert.AnError
	}
	return nil
}

func TestCostDataService_Interface(t *testing.T) {
	mockCostData := &domain.CostData{
		TotalCost: 30.25,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	mockHistory := []*domain.CostData{
		{TotalCost: 10.00, Currency: "USD", Timestamp: time.Now().AddDate(0, 0, -2)},
		{TotalCost: 20.00, Currency: "USD", Timestamp: time.Now().AddDate(0, 0, -1)},
		mockCostData,
	}

	service := &MockCostDataService{
		mockCostData: mockCostData,
		mockHistory:  mockHistory,
		shouldError:  false,
	}

	// Test GetCurrentCost
	currentCost, err := service.GetCurrentCost()
	assert.NoError(t, err)
	assert.Equal(t, 30.25, currentCost.TotalCost)

	// Test GetCostHistory
	history, err := service.GetCostHistory(7)
	assert.NoError(t, err)
	assert.Len(t, history, 3)
	assert.Equal(t, 30.25, history[2].TotalCost)

	// Test RefreshCostData
	err = service.RefreshCostData()
	assert.NoError(t, err)

	// Test error cases
	service.shouldError = true

	_, err = service.GetCurrentCost()
	assert.Error(t, err)

	_, err = service.GetCostHistory(7)
	assert.Error(t, err)

	err = service.RefreshCostData()
	assert.Error(t, err)
}
