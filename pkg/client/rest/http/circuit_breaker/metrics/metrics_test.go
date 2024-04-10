package metrics

import (
	"testing"
	"time"
)

func TestErrorRateUnderThreshold(t *testing.T) {
	c := newMetric(5*time.Second, time.Now)

	c.Success()
	c.Success()
	c.Success()
	c.Failure()
	c.Success()
	c.Failure()
	c.Success()
	c.Success()

	if r := c.Summary().rate; r == 0.0 {
		t.Errorf("expected error rate to be greater than zero,  got: %f in %+v", r, c.Summary())
	}
}

func TestErrorRateOverThreshold(t *testing.T) {
	c := newMetric(5*time.Second, time.Now)

	c.Failure()
	c.Failure()
	c.Failure()
	c.Failure()
	c.Failure()
	c.Failure()
	c.Failure()
	c.Success()
	c.Success()

	if ex, s := 0.70, c.Summary(); s.rate < ex {
		t.Errorf("expected error rate to be over %d%%, got: %f in %+v", int(ex*100), s.rate, s)
	}
}

func TestErrorRateCalculatedFromLast5Seconds(t *testing.T) {
	fakenow := time.Now()
	c := newMetric(5*time.Second, func() time.Time { return fakenow })

	// 77% error for 5 seconds
	for i := 0; i < 5; i++ {
		fakenow = fakenow.Add(time.Second)

		c.Failure()
		c.Failure()
		c.Success()
	}

	// 33.333% error for 5 seconds
	for i := 0; i < 5; i++ {
		fakenow = fakenow.Add(time.Second)

		c.Failure()
		c.Success()
		c.Success()
	}

	if ex, s := 0.34, c.Summary(); s.rate > ex {
		t.Errorf("expected error rate to be under %d%%, got: %f in %+v", int(ex*100), s.rate, s)
	}

}

func TestErrorRateCalculationWithTimeGap(t *testing.T) {
	fakenow := time.Now()
	c := newMetric(3*time.Second, func() time.Time { return fakenow })

	for i := 0; i < 3; i++ {
		fakenow = fakenow.Add(time.Second)
		c.Failure()
	}

	// the gap >= Metric.seconds
	fakenow = fakenow.Add(4 * time.Second)
	c.Success()

	if ex, s := 0.0, c.Summary(); s.rate != ex {
		t.Errorf("expected error rate to be %d%%, got: %f in %+v", int(ex*100), s.rate, s)
	}

	c = newMetric(3*time.Second, func() time.Time { return fakenow })
	for i := 0; i < 2; i++ {
		fakenow = fakenow.Add(time.Second)
		c.Failure()
	}

	// the gap < Metric.seconds
	fakenow = fakenow.Add(2 * time.Second)
	c.Success()

	if ex, s := 0.5, c.Summary(); s.rate != ex {
		t.Errorf("expected error rate to be %d%%, got: %f in %+v", int(ex*100), s.rate, s)
	}

	c = newMetric(3*time.Second, func() time.Time { return fakenow })
	for i := 0; i < 3; i++ {
		fakenow = fakenow.Add(time.Second)
		c.Success()
	}

	// clock jumps back
	fakenow = fakenow.Add(-time.Second)
	c.Failure()

	if ex, s := 1.0, c.Summary(); s.rate != ex {
		t.Errorf("expected error rate to be %d%%, got: %f in %+v", int(ex*100), s.rate, s)
	}
}
