package breaker

import (
	"errors"
	"github.com/skolldire/cash-manager-toolkit/pkg/client/rest/http/circuit_breaker/metrics"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit_breaker: circuit circuit_breaker open")

type states int

const (
	reset states = iota
	tripped
	closed
	open
	halfOpen
)

type Breaker struct {
	force   chan states
	allow   chan bool
	success chan struct{}
	failure chan struct{}

	config breaker.breakerConfig

	done chan struct{}
	wg   sync.WaitGroup
}

func newBreaker(c breaker.breakerConfig) *Breaker {
	if c.FailureRatio < 0.0 {
		c.FailureRatio = 0.0
	}

	if c.FailureRatio > 1.0 {
		c.FailureRatio = 1.0
	}

	if c.Window == 0 {
		c.Window = DefaultWindow
	}

	if c.Cooldown == 0 {
		c.Cooldown = DefaultCooldown
	}

	if c.Now == nil {
		c.Now = time.Now
	}

	if c.After == nil {
		c.After = time.After
	}

	b := Breaker{
		force:   make(chan states),
		allow:   make(chan bool),
		success: make(chan struct{}),
		failure: make(chan struct{}),
		done:    make(chan struct{}),
		config:  c,
	}

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.run()
	}()

	return &b
}

func NewBreaker(failureRatio float64, opts ...Option) *Breaker {
	config := breaker.breakerConfig{
		MinObservations: DefaultMinObservations,
		FailureRatio:    failureRatio,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return newBreaker(config)
}

func (b *Breaker) shouldOpen(m *metrics.metric) bool {
	s := m.Summary()
	return s.total > b.config.MinObservations && s.rate > b.config.FailureRatio
}

func (b *Breaker) run() {
	var (
		state   states
		timeout <-chan time.Time
		metrics *metrics.metric
	)

	for {
		switch state {
		case reset:
			metrics = metrics.newMetric(b.config.Window, b.config.Now)
			timeout = nil
			state = closed

		case closed:
			select {
			case b.allow <- true:
			case <-b.success:
				metrics.Success()
			case <-b.failure:
				metrics.Failure()
				if b.shouldOpen(metrics) {
					state = tripped
				}
			case state = <-b.force:
			case <-b.done:
				return
			}

		case tripped:
			timeout = b.config.After(b.config.Cooldown)
			state = open

		case open:
			select {
			case b.allow <- false:
			case <-b.success:
				state = reset
			case <-b.failure:
			case <-timeout:
				state = halfOpen
			case state = <-b.force:
			case <-b.done:
				return
			}

		case halfOpen:
			select {
			case b.allow <- true:
				state = tripped
			case <-b.success:
				state = reset
			case <-b.failure:
				state = tripped
			case state = <-b.force:
			case <-b.done:
				return
			}
		}
	}
}

func (b *Breaker) Success() {
	defer func() { _ = recover() }()
	b.success <- struct{}{}
}

func (b *Breaker) Failure() {
	defer func() { _ = recover() }()
	b.failure <- struct{}{}
}

func (b *Breaker) Allow() bool {
	return <-b.allow
}

func (b *Breaker) Do(f func() error) error {
	if !b.Allow() {
		return ErrCircuitOpen
	}

	if err := f(); err != nil {
		b.Failure()
		return err
	}

	b.Success()
	return nil
}

func (b *Breaker) Stop() {
	close(b.done)
	b.wg.Wait()

	close(b.allow)
	close(b.success)
	close(b.failure)
}

func (b *Breaker) forceState(s states) { b.force <- s } // nolint:unused
