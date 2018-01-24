package gon

import (
	"sync"
	"time"
)

var tickerid int64

// Ticker runs one or more functions, repeating at an interval.
type Ticker struct {
	sync.RWMutex
	wg       sync.WaitGroup
	interval int
	ticker   *time.Ticker
	quit     chan bool
	funcs    map[int64]TickerFunc
}

// TickerFunc is the signature of the callbacks run by the ticker.
type TickerFunc func(int64)

// NewTickerer
func NewTicker(i int) *Ticker {
	t := &Ticker{
		interval: i,
		funcs:    make(map[int64]TickerFunc)}
	t.quit = make(chan bool, 0)
	return t
}

// AddFunc adds anothe callback to the slice with a new ID.
func (t *Ticker) AddFunc(f TickerFunc) {
	t.Lock()
	defer t.Unlock()
	tickerid++
	t.funcs[tickerid] = f
}

// Start creates the time.Ticker and handles the calls at intervals.
func (t *Ticker) Start() {
	t.ticker = time.NewTicker(time.Second * time.Duration(t.interval))
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
