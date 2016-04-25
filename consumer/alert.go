package main

import (
	"container/ring"
	"fmt"
	"time"
)

type Alert struct {
	Active bool
	Time   time.Time
	Value  int
}

type AlertMonitor struct {
	valueSource  func() int
	interval     time.Duration
	threshold    int
	InAlert      bool
	alertHistory *ring.Ring
	quit         chan bool
}

func NewAlertMonitor(source func() int, interval time.Duration, threshold int) *AlertMonitor {
	r := &AlertMonitor{
		valueSource:  source,
		interval:     interval,
		threshold:    threshold,
		alertHistory: ring.New(6),
		quit:         make(chan bool),
	}

	go r.monitor()
	return r
}

func (am *AlertMonitor) monitor() {
	running := true
	for running {
		select {
		case <-time.After(am.interval):
			am.check()
		case <-am.quit:
			running = false
		}
	}
}

func (am *AlertMonitor) Stop() {
	close(am.quit)
}

func (am *AlertMonitor) check() {
	monitoredValue := am.valueSource()
	if monitoredValue > am.threshold {
		if !am.InAlert {
			am.InAlert = true
			am.alertHistory.Value = Alert{Active: true, Value: monitoredValue, Time: time.Now()}
			am.alertHistory = am.alertHistory.Next()
		}
	} else {
		if am.InAlert {
			am.InAlert = false
			am.alertHistory.Value = Alert{Active: false, Value: monitoredValue, Time: time.Now()}
			am.alertHistory = am.alertHistory.Next()
		}
	}
}

func (am *AlertMonitor) Print() {
	am.Do(func(alert Alert) {
		if alert.Active {
			fmt.Printf("[rising] High traffic generated an alert - hits = %d, triggered at %s\n", alert.Value, alert.Time)
		} else {
			fmt.Printf("[clearing] Low traffic generated an alert - hits = %d, triggered at %s\n", alert.Value, alert.Time)
		}
	})
}

func (am *AlertMonitor) Do(eachFunc func(alert Alert)) {
	am.alertHistory.Do(func(rAlert interface{}) {
		alert, ok := rAlert.(Alert)
		if !ok {
			return
		}
		eachFunc(alert)
	})
}
