package domain

import (
	"time"
)

// CostData represents the cost information from ccusage
type CostData struct {
	TotalCost      float64            `json:"total_cost"`
	Currency       string             `json:"currency"`
	Timestamp      time.Time          `json:"timestamp"`
	ModelBreakdown map[string]float64 `json:"model_breakdown,omitempty"`
}

// CostDataRepository defines the interface for fetching cost data
type CostDataRepository interface {
	FetchCostData() (*CostData, error)
	GetLastUpdated() (time.Time, error)
}

// CostDataService defines the business logic for cost data operations
type CostDataService interface {
	GetCurrentCost() (*CostData, error)
	GetCostHistory(days int) ([]*CostData, error)
	RefreshCostData() error
}
