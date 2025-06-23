package domain

import "errors"

// Common errors used across the domain
var (
	ErrPluginNotEnabled = errors.New("plugin is not enabled")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrDataNotFound     = errors.New("data not found")
)
