package gon

import (
	"sync"
)

// Scheduler holds pointers to all the tickers and timers.
type Scheduler struct {
	sync.RWMutex
	seconds      map[int]*Ticker
	delayseconds map[int]*Ticker
	minutes      map[int]*Ticker
	delayminutes map[int]*Ticker
	hours        map[int]*Ticker
	delayhours   map[int]*Ticker
}

// NewScheduler returns a Scheduler populated with maps.
func NewScheduler() *Scheduler {
	return &Scheduler{
		seconds:      make(map[int]*Ticker),
		delayseconds: make(map[int]*Ticker),
		minutes:      make(map[int]*Ticker),
		delayminutes: make(map[int]*Ticker),
		hours:        make(map[int]*Ticker),
		delayhours:   make(map[int]*Ticker),
	}
}

// AddEveryNSecond adds a repeating task based on a seconds interval.
func (sc *Scheduler) AddEveryNSecond(f TickerFunc, n int) {
	sc.Lock()
	defer sc.Unlock()
	t, ok := sc.seconds[n]
	if !ok {
		t = NewTicker(n)
		sc.seconds[n] = t
	}
	t.AddFunc(f)
	t.Start()
}

// Wait for all waitgroups in tickers and timers.
func (sc *Scheduler) Wait() {
	for _, t := range sc.seconds {
		t.Wait()
	}
}
