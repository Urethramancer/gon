package gon

import (
	"sync"
	"time"
)

// Ticker runs one or more functions, repeating at an interval.
type Ticker struct {
	sync.RWMutex
	wg       sync.WaitGroup
	interval int
	duration time.Duration
	ticker   *time.Ticker
	funcs    map[int64]TickerFunc
	quit     chan bool
}

// TickerFunc is the signature of the callbacks run by the ticker.
type TickerFunc func(int64)

// NewTicker
func NewTicker(i int, d time.Duration) *Ticker {
	t := &Ticker{
		interval: i,
		duration: d,
		funcs:    make(map[int64]TickerFunc)}
	t.quit = make(chan bool, 0)
	return t
}

// AddFunc adds another callback to the funcs map with a new ID.
func (t *Ticker) AddFunc(f TickerFunc, id int64) {
	t.Lock()
	defer t.Unlock()
	t.funcs[id] = f
}

// Start creates the time.Ticker and handles the calls at intervals.
func (t *Ticker) Start() {
	t.ticker = time.NewTicker(t.duration)
	for {
		select {
		case <-t.ticker.C:
			t.RLock()
			for k, f := range t.funcs {
				go func(id int64, tf TickerFunc) {
					t.wg.Add(1)
					tf(id)
					t.wg.Done()
				}(k, f)
			}
			t.RUnlock()
		case <-t.quit:
			t.ticker.Stop()
			return
		}
	}
}

// Stop sends a signal to the quit channel.
func (t *Ticker) Stop() {
	t.quit <- true
}

// Wait for the current scheduled task in the ticker to finish.
func (t *Ticker) Wait() {
	t.wg.Wait()
}
