package main

import (
	"testing"
	"time"
)

func alertRun(dur time.Duration, threshold int, emitter func() int) *AlertMonitor {
	al := NewAlertMonitor(emitter, dur, threshold)
	time.Sleep(1 * time.Second)
	al.Stop()
	return al
}

func TestAlertCross(t *testing.T) {
	emitter := func() int {
		return 10
	}
	al := alertRun(10*time.Millisecond, 100, emitter)
	al.Do(func(alert Alert) {
		t.Fatal("Should not alert unless threshold has been crossed")
	})

	al = alertRun(10*time.Millisecond, 10, emitter)
	al.Do(func(alert Alert) {
		t.Fatal("Should not alert unless threshold has been crossed")
	})
}

func TestAlertFlap(t *testing.T) {
	emitter := func() int {
		return 10
	}
	al := alertRun(10*time.Millisecond, 9, emitter)
	alertCount := 0
	al.Do(func(alert Alert) {
		alertCount += 1
	})

	if alertCount != 1 {
		t.Fatal("Should have just 1 alert")
	}

	flap := false
	emitter = func() int {
		flap = !flap
		if flap {
			return 100
		} else {
			return 1
		}
	}

	// underlying ringbuffer has an even number of slots
	al = alertRun(10*time.Millisecond, 10, emitter)
	alOn := 0
	alOff := 0
	al.Do(func(alert Alert) {
		if alert.Active {
			alOn += 1
		} else {
			alOff += 1
		}
	})

	if alOn != alOff {
		t.Fatal("Alert buffer should have same number of rise/clear alerts")
	}
}
