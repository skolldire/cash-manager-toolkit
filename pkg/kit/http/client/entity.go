package client

import (
	"net/http"
	"time"
)

const (
	defaultTimeout           = 5
	defaultCBRatio           = 0.75
	defaultCBWindow          = 5
	defaultCBCoolDown        = 5
	defaultCBMinObservations = 20
	defaultMaxRetries        = 3
	defaultBackoffMin        = 250
	defaultBackoffMax        = 5000
)

type Configuration struct {
	DialTimeOut        time.Duration `json:"dial_timeout"`
	DisableTimeout     bool          `json:"disable_timeout"`
	Timeout            time.Duration `json:"timeout"`
	WithCircuitBreaker bool          `json:"with_circuit_breaker"`
	CBFailureRatio     float64       `json:"cb_failure_ratio"`
	CBWindow           time.Duration `json:"cb_window"`
	CBCoolDown         time.Duration `json:"cb_cool_down"`
	CBMinObservations  uint          `json:"cb_min_observations"`
	WithRetry          bool          `json:"with_retry"`
	MaxRetries         int           `json:"max_retries"`
	BackoffMin         time.Duration `json:"backoff_min"`
	BackoffMax         time.Duration `json:"backoff_max"`
}

type Service interface {
	Init() http.Client
}
