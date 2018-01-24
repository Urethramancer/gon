package gon

import (
	"sync"
)

// Scheduler holds pointers to all the tickers and timers.
type Scheduler struct {
	sync.RWMutex
	seconds map[int]*Ticker
	minutes map[int]*Ticker
	hours   map[int]*Ticker
	alarms  map[int]*Alarm
}

// NewScheduler returns a Scheduler populated with maps.
func NewScheduler() *Scheduler {
	return &Scheduler{
		seconds: make(map[int]*Ticker),
		minutes: make(map[int]*Ticker),
		hours:   make(map[int]*Ticker),
	}
}

func (sc *Scheduler) addTicker(m *map[int]*Ticker, n int, f TickerFunc) {
	sc.Lock()
	defer sc.Unlock()
	t, ok := (*m)[n]
	if !ok {
		t = NewTicker(n)
		(*m)[n] = t
	}
	t.AddFunc(f)
	go t.Start()
}

// RepeatSeconds adds a repeating task based on a seconds interval.
func (sc *Scheduler) RepeatSeconds(n int, f TickerFunc) {
	sc.addTicker(&sc.seconds, n, f)
}

// RepeatMinutes adds a repeating task on a minute-based interval.
func (sc *Scheduler) RepeatMinutes(n int, f TickerFunc) {
	sc.addTicker(&sc.minutes, n, f)
}

// RepeatHours adds a repeating task on an hour-based interval.
func (sc *Scheduler) RepeatHours(n int, f TickerFunc) {
	sc.addTicker(&sc.hours, n, f)
}

// Wait for all waitgroups in tickers and timers.
func (sc *Scheduler) Wait() {
	for _, t := range sc.seconds {
		t.Wait()
	}
}
