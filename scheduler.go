package gon

import (
	"sync"
)

type Scheduler struct {
	sync.RWMutex
	seconds      map[int]*Ticker
	delayseconds map[int]*Ticker
	minutes      map[int]*Ticker
	delayminutes map[int]*Ticker
	hours        map[int]*Ticker
	delayhours   map[int]*Ticker
}

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
	timer, ok := sc.seconds[n]
	if !ok {
		timer = NewTimer(n)
		sc.seconds[n] = timer
	}
	timer.AddFunc(f)
	timer.Start()
}

// Wait for all waitgroups in tickers and timers.
func (sc *Scheduler) Wait() {
	for _, t := range sc.seconds {
		t.Wait()
	}
}
