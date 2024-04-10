package bucket

import (
	"errors"
	breaker2 "github.com/skolldire/cash-manager-toolkit/pkg/client/rest/http/circuit_breaker/breaker"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestBucketBreaker(t *testing.T) {
	cb := NewBucketBreaker(0.1, breaker2.WithMinObservations(5))

	allowed, success, _ := cb.Allow("bucket_1")
	if !allowed {
		t.Fatal("expected circuit_breaker to allow operation on bucket_1")
	}
	success()

	for i := 0; i < 5; i++ {
		allowed, _, failure := cb.Allow("bucket_1")
		if !allowed {
			t.Fatal("expected circuit_breaker to allow operation on bucket_1")
		}
		failure()
	}

	{
		allowed, _, _ := cb.Allow("bucket_1")
		if allowed {
			t.Fatal("expected circuit_breaker to disallow operation after fail threshold was reach")
		}
	}

	{
		allowedB2, _, _ := cb.Allow("bucket_2")
		if !allowedB2 {
			t.Fatal("expected circuit_breaker to allow operation on bucket_2")
		}
	}

	{
		cb.Remove("bucket_1")
		allowed, _, _ = cb.Allow("bucket_1")
		if !allowed {
			t.Fatal("expected circuit_breaker to allow operation after remove was called")
		}
	}

	{
		// Multiple calls to remove should not panic
		cb.Remove("bucket_1")
		cb.Remove("bucket_1")
	}
}

func TestBucketBreakerDo(t *testing.T) {
	cb := NewBucketBreaker(0.1)
	defer cb.Remove("bucket_1")

	mockErr := errors.New("error")
	failFunc := func() error { return mockErr }
	okFunc := func() error { return nil }
	unexpectedFunc := func() error { return errors.New("unexpected call") }

	if cb.Do("bucket_1", okFunc) != nil {
		t.Fatal("expected Breaker.Do to not return en error")
	}

	for i := 0; i < 10; i++ {
		if err := cb.Do("bucket_1", failFunc); err != mockErr {
			t.Fatal("expected Breaker.Do to return failFunc error")
		}
	}

	if err := cb.Do("bucket_1", unexpectedFunc); err != breaker2.ErrCircuitOpen {
		t.Fatalf("expected error to be %q, got: %q", breaker2.ErrCircuitOpen, err)
	}
}

func TestBucketBreaker_ConcurrentUse(t *testing.T) {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	buckets := []string{"bucket_1", "bucket_2", "bucket_3", "bucket_4"}
	randBucket := func() string { return buckets[source.Intn(len(buckets))] }

	const concurrency = 20
	const iterations = 10000

	cb := NewBucketBreaker(0.1)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		bucket := randBucket()
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				allowed, success, _ := cb.Allow(bucket)
				if !allowed {
					t.Errorf("expected circuit_breaker to allow operation on bucket %q", bucket)
					return
				}
				success()
			}
		}()
	}

	wg.Wait()
	for _, bucket := range buckets {
		cb.Remove(bucket)
	}
}
