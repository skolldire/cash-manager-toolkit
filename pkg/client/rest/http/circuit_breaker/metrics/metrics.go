package metrics

import (
	"container/ring"
	"time"
)

type counter struct {
	second  int64
	success uint
	failure uint
}

func (c *counter) reset(second int64) {
	c.failure = 0
	c.success = 0
	c.second = second
}

type summary struct {
	total  uint
	errors uint
	rate   float64
}

type metric struct {
	r       *ring.Ring
	seconds uint
	now     func() time.Time
}

func newMetric(window time.Duration, now func() time.Time) *metric {
	seconds := int(window / time.Second)

	if seconds <= 0 {
		panic("metrics must have a window of at least 1 Second")
	}

	r := ring.New(seconds)
	for i := 0; i < seconds; i++ {
		r.Value = &counter{}
		r = r.Next()
	}

	return &metric{r: r, seconds: uint(seconds), now: now}
}

func (m *metric) next() *counter {
	bucket := m.now().Unix()
	c := m.r.Value.(*counter)
	if c.second == bucket {
		return c
	}

	step := bucket - c.second
	// consider the data are invalid when clock jumps back
	if step < 0 || step > int64(m.seconds) {
		step = int64(m.seconds)
	}

	for i := int64(1); i <= step; i++ {
		m.r = m.r.Next()
		c = m.r.Value.(*counter)
		c.reset(bucket - step + i)
	}

	return c
}

func (m *metric) Success() {
	m.next().success++
}

func (m *metric) Failure() {
	m.next().failure++
}

func (m metric) Summary() summary {
	var sum summary

	m.r.Do(func(v interface{}) {
		c := v.(*counter)
		sum.total += c.success + c.failure
		sum.errors += c.failure
	})

	if sum.total > 0 {
		sum.rate = float64(sum.errors) / float64(sum.total)
	}

	return sum
}
