package config

import "errors"

var (
	ErrDatabaseURL = errors.New("database url is required")
	ErrDuration    = errors.New("durations must be positive")
)
