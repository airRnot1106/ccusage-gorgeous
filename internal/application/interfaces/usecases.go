package interfaces

import (
	"context"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
)

// CostFetcher defines the use case for fetching cost data
type CostFetcher interface {
	GetCurrentCost(ctx context.Context) (*domain.CostData, error)
	GetCostHistory(ctx context.Context, days int) ([]*domain.CostData, error)
	RefreshCostData(ctx context.Context) error
}

// Animator defines the use case for animation control
type Animator interface {
	GenerateAnimationFrame(ctx context.Context, text string, frameNumber int) (*domain.AnimationFrame, error)
	GetAnimationConfig(ctx context.Context) (*domain.AnimationConfig, error)
	UpdateAnimationConfig(ctx context.Context, config *domain.AnimationConfig) error
	StartAnimation(ctx context.Context) error
	StopAnimation(ctx context.Context) error
}

// Displayer defines the use case for display control
type Displayer interface {
	RenderDisplay(ctx context.Context, costData *domain.CostData, animationFrame *domain.AnimationFrame) (string, error)
	GetDisplayConfig(ctx context.Context) (*domain.DisplayConfig, error)
	UpdateDisplayConfig(ctx context.Context, config *domain.DisplayConfig) error
	GetDisplayCapabilities(ctx context.Context) (*DisplayCapabilities, error)
}

// AppController defines the main application controller
type AppController interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Refresh(ctx context.Context) error
	GetStatus(ctx context.Context) (*AppStatus, error)
}

// AppStatus represents the current status of the application
type AppStatus struct {
	IsRunning     bool             `json:"is_running"`
	LastUpdate    time.Time        `json:"last_update"`
	CurrentCost   *domain.CostData `json:"current_cost,omitempty"`
	ActivePlugins []string         `json:"active_plugins"`
	ErrorCount    int              `json:"error_count"`
	LastError     string           `json:"last_error,omitempty"`
}
