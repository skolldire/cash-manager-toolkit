package breaker

import (
	"time"
)

const (
	DefaultWindow          = 5 * time.Second
	DefaultCooldown        = 1 * time.Second
	DefaultMinObservations = 15
)

type breakerConfig struct {
	FailureRatio float64

	Window          time.Duration
	Cooldown        time.Duration
	MinObservations uint

	Now   func() time.Time
	After func(time.Duration) <-chan time.Time
}

type Option func(*breakerConfig)

func WithMinObservations(min uint) Option {
	return func(config *breakerConfig) {
		config.MinObservations = min
	}
}

func WithWindow(w time.Duration) Option {
	return func(config *breakerConfig) {
		config.Window = w
	}
}

func WithCooldown(c time.Duration) Option {
	return func(config *breakerConfig) {
		config.Cooldown = c
	}
}
