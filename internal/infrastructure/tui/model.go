package tui

import (
	"context"
	"strconv"
	"time"

	"github.com/airRnot1106/ccusage-gorgeous/internal/core"
	"github.com/airRnot1106/ccusage-gorgeous/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the TUI application model
type Model struct {
	ctx         context.Context
	registry    *core.PluginRegistry
	config      *core.ConfigManager
	width       int
	height      int
	frameCount  int
	currentCost *domain.CostData
	lastUpdate  time.Time
	error       error
	isLoading   bool
	isQuitting  bool
}

// NewModel creates a new TUI model
func NewModel(ctx context.Context, registry *core.PluginRegistry, config *core.ConfigManager) *Model {
	return &Model{
		ctx:        ctx,
		registry:   registry,
		config:     config,
		frameCount: 0,
		isLoading:  true,
	}
}

// Init initializes the TUI model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchCostData(),
		m.tick(),
	)
}

// Update handles TUI updates
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.isQuitting = true
			return m, tea.Quit
		case "r":
			// Refresh data
			return m, m.fetchCostData()
		}

	case costDataMsg:
		m.currentCost = msg.costData
		m.error = msg.err
		m.lastUpdate = time.Now()
		m.isLoading = false
		return m, nil

	case tickMsg:
		m.frameCount++
		if m.isQuitting {
			return m, nil
		}
		return m, m.tick()

	case errorMsg:
		m.error = msg.err
		m.isLoading = false
		return m, nil
	}

	return m, nil
}

// View renders the TUI view
func (m *Model) View() string {
	if m.isQuitting {
		return "Goodbye! 👋\n"
	}

	if m.isLoading {
		return "Loading cost data... ⏳\n"
	}

	if m.error != nil {
		return "Error: " + m.error.Error() + "\n\nPress 'r' to retry or 'q' to quit.\n"
	}

	if m.currentCost == nil {
		return "No cost data available.\n\nPress 'r' to refresh or 'q' to quit.\n"
	}

	// Get active plugins
	animationPlugin, err := m.registry.GetActiveAnimation()
	if err != nil {
		return "Error getting animation plugin: " + err.Error() + "\n"
	}

	displayPlugin, err := m.registry.GetActiveDisplay()
	if err != nil {
		return "Error getting display plugin: " + err.Error() + "\n"
	}

	// Generate animation frame
	animationConfig := m.config.GetAnimationConfig()
	costText := "$" + formatFloat(m.currentCost.TotalCost)

	animationFrame, err := animationPlugin.GenerateFrame(m.ctx, costText, m.frameCount, animationConfig)
	if err != nil {
		return "Error generating animation: " + err.Error() + "\n"
	}

	// Create display data
	displayConfig := m.config.GetDisplayConfig()
	if displayConfig != nil {
		displayConfig.Size.Width = m.width
		displayConfig.Size.Height = m.height
	}

	displayData := &domain.DisplayData{
		Cost:        m.currentCost,
		Animation:   animationFrame,
		Config:      displayConfig,
		LastUpdated: m.lastUpdate,
	}

	// Render display
	output, err := displayPlugin.Render(m.ctx, displayData)
	if err != nil {
		return "Error rendering display: " + err.Error() + "\n"
	}

	// Add controls help
	output += "\n\nControls: 'r' to refresh, 'q' to quit\n"

	return output
}

// Messages for the TUI update loop
type (
	costDataMsg struct {
		costData *domain.CostData
		err      error
	}
	tickMsg  struct{}
	errorMsg struct{ err error }
)

// fetchCostData fetches cost data from the active data source plugin
func (m *Model) fetchCostData() tea.Cmd {
	return func() tea.Msg {
		dataSourcePlugin, err := m.registry.GetActiveDataSource()
		if err != nil {
			return costDataMsg{nil, err}
		}

		costData, err := dataSourcePlugin.FetchCostData(m.ctx)
		return costDataMsg{costData, err}
	}
}

// tick creates a tick command for animation
func (m *Model) tick() tea.Cmd {
	animationConfig := m.config.GetAnimationConfig()
	if animationConfig == nil || !animationConfig.Enabled {
		return tea.Tick(1*time.Second, func(time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return tea.Tick(animationConfig.Speed, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// formatFloat formats a float64 to a string with 2 decimal places
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
