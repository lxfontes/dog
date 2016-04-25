package main

import (
	"container/ring"
	"sync"
	"time"
)

type DataPoints map[string][]int

func (dp DataPoints) Push(key string, val int) {
	_, ok := dp[key]
	if !ok {
		dp[key] = []int{}
	}
	dp[key] = append(dp[key], val)
}

type TimeData struct {
	Time time.Time
	Data DataPoints
}

func (td *TimeData) Count() int {
	return len(td.Data)
}

func (td *TimeData) CountEntries(key string) int {
	return len(td.Data[key])
}

func (td *TimeData) SumEntries(key string) int {
	t := 0
	for _, k := range td.Data[key] {
		t += k
	}

	return t
}

func (td *TimeData) AverageEntries(key string) int {
	c := td.CountEntries(key)
	if c > 0 {
		return td.SumEntries(key) / c
	}

	return 0
}

type SlidingCounter struct {
	Window time.Duration
	Slots  int

	curData DataPoints
	data    *ring.Ring

	sync.RWMutex
}

func NewSlidingCounter(window time.Duration, slots int) *SlidingCounter {
	ret := &SlidingCounter{
		Window:  window,
		Slots:   slots,
		data:    ring.New(slots),
		curData: make(DataPoints),
	}

	go ret.start()
	return ret
}

func (tc *SlidingCounter) start() {
	for {
		select {
		case <-time.After(tc.Window):
			tc.cut()
		}
	}
}

func (tc *SlidingCounter) cut() {
	dt := TimeData{
		Time: time.Now(),
		Data: tc.curData,
	}

	tc.Lock()
	tc.data.Value = dt
	tc.data = tc.data.Next()
	tc.Unlock()

	tc.curData = make(DataPoints)
}

func (tc *SlidingCounter) Push(key string, value int) {
	tc.Lock()
	tc.curData.Push(key, value)
	tc.Unlock()
}

func (tc *SlidingCounter) AverageValuesBack(key string, amount time.Duration) int {
	tc.RLock()
	defer tc.RUnlock()

	delta := time.Now().Add(amount)

	r := tc.data.Prev()
	t := 0
	n := 0
	for r != nil {
		td, ok := r.Value.(TimeData)
		if !ok {
			break
		}

		if td.Time.After(delta) {
			t += td.SumEntries(key)
			n += td.CountEntries(key)
		} else {
			break
		}
		r = r.Prev()
	}

	if n == 0 {
		return 0
	}

	return t / n
}

func (tc *SlidingCounter) CountEntriesBack(key string, amount time.Duration) int {
	tc.RLock()
	defer tc.RUnlock()

	delta := time.Now().Add(amount)

	r := tc.data.Prev()
	t := 0
	for r != nil {
		td, ok := r.Value.(TimeData)
		if !ok {
			break
		}

		if td.Time.After(delta) {
			t += td.CountEntries(key)
		} else {
			break
		}
		r = r.Prev()
	}

	return t
}

func (tc *SlidingCounter) SumEntriesBack(key string, amount time.Duration) int {
	tc.RLock()
	defer tc.RUnlock()

	delta := time.Now().Add(amount)

	r := tc.data.Prev()
	t := 0
	for r != nil {
		td, ok := r.Value.(TimeData)
		if !ok {
			break
		}

		if td.Time.After(delta) {
			t += td.SumEntries(key)
		} else {
			break
		}
		r = r.Prev()
	}

	return t
}

func (tc *SlidingCounter) AverageEntriesBack(key string, amount time.Duration) int {
	c := tc.CountEntriesBack(key, amount)
	if c == 0 {
		return 0
	}

	return tc.SumEntriesBack(key, amount) / c
}
