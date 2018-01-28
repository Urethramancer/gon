package gon

import (
	"sync"
	"time"
)

// Ticker runs one or more functions, repeating at an interval.
type Ticker struct {
	sync.RWMutex
	wg       sync.WaitGroup
	duration time.Duration
	ticker   *time.Ticker
	funcs    map[int64]EventFunc
	quit     chan bool
}

// NewTicker creates the Ticker structure and quit channel.
func NewTicker(d time.Duration) *Ticker {
	t := &Ticker{
		duration: d,
		funcs:    make(map[int64]EventFunc)}
	t.quit = make(chan bool)
	return t
}

// AddFunc adds another callback to the funcs map with a new ID.
func (t *Ticker) AddFunc(f EventFunc, id int64) {
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
				t.wg.Add(1)
				go func(id int64, tf EventFunc) {
					tf(id)
					t.wg.Done()
				}(k, f)
			}
			t.RUnlock()
		case <-t.quit:
			t.ticker.Stop()
			for k := range t.funcs {
				delete(t.funcs, k)
			}
			return
		}
	}
}

// Stop the ticker.
func (t *Ticker) Stop() {
	t.Wait()
	t.quit <- true
}

// Wait for the current scheduled task in the ticker to finish.
func (t *Ticker) Wait() {
	t.wg.Wait()
}
