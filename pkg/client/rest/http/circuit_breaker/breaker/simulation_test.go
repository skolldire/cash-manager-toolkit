package breaker

import (
	"sync"
	"testing"
	"time"
)

func failrate(b *Breaker, count int, pct float64) {
	chance := int(1 / pct)
	if chance <= 0 {
		chance = 1
	}

	for i := 0; i < count; i++ {
		if b.Allow() {
			if (i%count)%chance == 0 {
				b.Failure()
			} else {
				b.Success()
			}
		}
	}
}

func TestSimulateConcurrentBreakerHandlerWithPartialFailures(t *testing.T) {
	const requestsPerSecond = 100
	const seconds = 5

	var mu sync.RWMutex

	now := time.Now()
	after := make(chan time.Time)

	b := newBreaker(breakerConfig{
		Window:          seconds * time.Second,
		MinObservations: requestsPerSecond / seconds,
		FailureRatio:    0.05,
		Now: func() time.Time {
			mu.RLock()
			defer mu.RUnlock()
			return now
		},
		After: func(time.Duration) <-chan time.Time { return after },
	})
	defer b.Stop()

	for i := 0; i < seconds; i++ {
		failrate(b, requestsPerSecond, 0.20)
		mu.Lock()
		now = now.Add(time.Second)
		mu.Unlock()
	}

	if got, want := b.Allow(), false; got != want {
		t.Fatalf("expected to trip at a high failure rate")
	}

	mu.RLock()
	aftNow := now
	mu.RUnlock()

	after <- aftNow

	if got, want := b.Allow(), true; got != want {
		t.Fatalf("expected to allow in half-open state after cooldown")
	}

	b.Success()

	if got, want := b.Allow(), true; got != want {
		t.Fatalf("expected to close after success from half-open")
	}

	for i := 0; i < seconds; i++ {
		failrate(b, requestsPerSecond, 0.02)
		mu.Lock()
		now = now.Add(time.Second)
		mu.Unlock()
	}

	if got, want := b.Allow(), true; got != want {
		t.Fatalf("expected to stay closed after lower error rate")
	}

	for i := 0; i < seconds; i++ {
		failrate(b, requestsPerSecond, 0.06)
		mu.Lock()
		now = now.Add(time.Second)
		mu.Unlock()
	}

	if got, want := b.Allow(), false; got != want {
		t.Fatalf("expected to open after high error rate again")
	}
}
