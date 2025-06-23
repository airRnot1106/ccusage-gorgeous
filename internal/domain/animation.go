package domain

import (
	"time"
)

// AnimationConfig represents configuration for rainbow animation
type AnimationConfig struct {
	Speed   time.Duration    `json:"speed"`
	Colors  []string         `json:"colors"`
	Enabled bool             `json:"enabled"`
	Pattern AnimationPattern `json:"pattern"`
}

// AnimationPattern defines the type of animation pattern
type AnimationPattern string

const (
	PatternRainbow  AnimationPattern = "rainbow"
	PatternGradient AnimationPattern = "gradient"
	PatternPulse    AnimationPattern = "pulse"
	PatternWave     AnimationPattern = "wave"
)

// AnimationFrame represents a single frame of animation
type AnimationFrame struct {
	Colors    []string  `json:"colors"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

// AnimationService defines the interface for animation operations
type AnimationService interface {
	GenerateFrame(text string, frameNumber int) (*AnimationFrame, error)
	GetDefaultConfig() *AnimationConfig
	ValidateConfig(config *AnimationConfig) error
}
