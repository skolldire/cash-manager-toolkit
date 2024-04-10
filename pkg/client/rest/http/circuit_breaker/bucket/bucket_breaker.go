package bucket

import (
	breaker2 "github.com/skolldire/cash-manager-toolkit/pkg/client/rest/http/circuit_breaker/breaker"
	"sync"
)

type BucketBreaker struct {
	breakerFactory func() *breaker2.Breaker

	mu       sync.RWMutex
	breakers map[string]*breaker2.Breaker
}

func NewBucketBreaker(failureRatio float64, opts ...breaker2.Option) *BucketBreaker {
	config := breaker.breakerConfig{
		MinObservations: breaker2.DefaultMinObservations,
		FailureRatio:    failureRatio,
	}

	for _, opt := range opts {
		opt(&config)
	}

	b := BucketBreaker{
		breakerFactory: func() *breaker2.Breaker {
			return breaker.newBreaker(config)
		},
		breakers: make(map[string]*breaker2.Breaker),
	}

	return &b
}

func (b *BucketBreaker) Allow(bucket string) (allowed bool, success, failure func()) {
	cb := b.breakerForBucket(bucket)
	return cb.Allow(), cb.Success, cb.Failure
}

func (b *BucketBreaker) Do(bucket string, f func() error) error {
	cb := b.breakerForBucket(bucket)
	if !cb.Allow() {
		return breaker2.ErrCircuitOpen
	}

	if err := f(); err != nil {
		cb.Failure()
		return err
	}

	cb.Success()
	return nil
}

func (b *BucketBreaker) Remove(bucket string) {
	b.mu.RLock()
	cb, ok := b.breakers[bucket]
	b.mu.RUnlock()
	if !ok {
		return
	}

	b.mu.Lock()
	delete(b.breakers, bucket)
	b.mu.Unlock()

	cb.Stop()
}

func (b *BucketBreaker) breakerForBucket(bucket string) *breaker2.Breaker {
	b.mu.RLock()
	cb, ok := b.breakers[bucket]
	b.mu.RUnlock()

	if ok {
		return cb
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	cb, ok = b.breakers[bucket]
	if ok {
		return cb
	}

	cb = b.breakerFactory()
	b.breakers[bucket] = cb
	return cb
}
