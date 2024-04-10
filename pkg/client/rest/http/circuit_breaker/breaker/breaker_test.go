package breaker

import (
	"errors"
	"testing"
	"time"
)

func TestNewBreakerAllows(t *testing.T) {
	c := NewBreaker(0)
	defer c.Stop()

	if !c.Allow() {
		t.Fatal("expected new Breaker to be closed")
	}
}

func TestBreakerSuccessClosesOpenBreaker(t *testing.T) {
	b := NewBreaker(0)
	defer b.Stop()

	b.forceState(tripped)

	if b.Allow() {
		t.Fatal("expected new Breaker to be open after being tripped")
	}

	b.Success()

	if !b.Allow() {
		t.Fatal("expected new Breaker to be closed after a success")
	}
}

func TestBreakerFailTripsBreakerWithASingleFailureAt0PercentThreshold(t *testing.T) {
	b := NewBreaker(0)
	defer b.Stop()

	for i := 0; i < 100; i++ {
		b.Success()
	}

	b.Failure()

	if b.Allow() {
		t.Fatalf("expected failure to not trip circuit at 0%% threshold")
	}
}

func TestBreakerFailDoesNotTripBreakerAt1PercentThreshold(t *testing.T) {
	const threshold = 0.01

	c := NewBreaker(threshold)
	defer c.Stop()

	for i := 0; i < 100-100*threshold; i++ {
		c.Success()
	}

	for i := 0; i < 100*threshold; i++ {
		c.Failure()
	}

	if !c.Allow() {
		t.Fatalf("expected failure to not trip circuit at 1%% threshold")
	}

	c.Failure()

	if c.Allow() {
		t.Fatal("expected failure to trip over the threshold")
	}
}

func TestBreakerAllowsASingleRequestAfterNapTime(t *testing.T) {
	after := make(chan time.Time)

	c := newBreaker(breaker.breakerConfig{
		Window: 5 * time.Second,
		After:  func(time.Duration) <-chan time.Time { return after },
	})
	defer c.Stop()

	c.forceState(tripped)

	after <- time.Now()

	if !c.Allow() {
		t.Fatal("expected to allow once after nap time")
	}

	if c.Allow() {
		t.Fatal("expected to only allow once after nap time")
	}
}

func TestBreakerClosesAfterSuccessAfterNapTime(t *testing.T) {
	after := make(chan time.Time)

	b := newBreaker(breaker.breakerConfig{
		Window: 5 * time.Second,
		After:  func(time.Duration) <-chan time.Time { return after },
	})
	defer b.Stop()

	b.forceState(tripped)

	after <- time.Now()

	if !b.Allow() {
		t.Fatal("expected to allow once after nap time")
	}

	b.Success()

	if !b.Allow() {
		t.Fatal("expected to close after first success")
	}

	if !b.Allow() {
		t.Fatal("expected to stay closed after first success")
	}
}

func TestBreakerReschedulesOnFailureInHalfOpen(t *testing.T) {
	afters := make(chan chan time.Time)

	b := newBreaker(breaker.breakerConfig{
		Window: 5 * time.Second,
		After: func(time.Duration) <-chan time.Time {
			after := make(chan time.Time)
			afters <- after
			return after
		},
	})
	defer b.Stop()

	b.forceState(tripped)

	(<-afters) <- time.Now()

	b.Failure()

	select {
	case after := <-afters:
		after <- time.Now()
	case <-time.After(time.Millisecond):
		t.Fatal("expected to reschedule after failure in half-open, did not")
	}

	b.Success()

	if !b.Allow() {
		t.Fatal("expected to close after failure in half-open")
	}
}

func TestNewBreakerOptions(t *testing.T) {
	c := NewBreaker(0,
		WithCooldown(2*time.Second),
		WithMinObservations(30),
		WithWindow(30*time.Second),
	)
	defer c.Stop()

	if c.config.Cooldown != 2*time.Second || c.config.MinObservations != 30 || c.config.Window != 30*time.Second {
		t.Fatal("unexpected configuration values")
	}
}

func TestBreakerDo(t *testing.T) {
	b := NewBreaker(0.1)
	defer b.Stop()

	mockErr := errors.New("error")
	failFunc := func() error { return mockErr }
	okFunc := func() error { return nil }
	unexpectedFunc := func() error { return errors.New("unexpected call") }

	if b.Do(okFunc) != nil {
		t.Fatal("expected Breaker.Do to not return en error")
	}

	for i := 0; i < 10; i++ {
		if err := b.Do(failFunc); err != mockErr {
			t.Fatal("expected Breaker.Do to return failFunc error")
		}
	}

	if err := b.Do(unexpectedFunc); err != ErrCircuitOpen {
		t.Fatalf("expected error to be %q, got: %q", ErrCircuitOpen, err)
	}
}

func TestBreakerStop(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	b := newBreaker(breaker.breakerConfig{FailureRatio: 0.1})

	if b.Allow() != true {
		t.Fatalf("expected allow to return true")
	}

	b.Success()
	b.Failure()

	b.Stop()

	if b.Allow() != false {
		t.Fatalf("expected allow to return false after circuit_breaker was stopped")
	}

	b.Success()
	b.Failure()
}
