package gon

import (
	"sync"
	"time"
)

var tickerid int64

type Ticker struct {
	sync.RWMutex
	wg       sync.WaitGroup
	interval int
	ticker   *time.Ticker
	quit     chan bool
	funcs    map[int64]TickerFunc
}

type TickerFunc func(int64)

func NewTimer(i int) *Ticker {
	t := &Ticker{
		interval: i,
		funcs:    make(map[int64]TickerFunc)}
	t.quit = make(chan bool, 0)
	return t
}

func (t *Ticker) AddFunc(f TickerFunc) {
	t.Lock()
	defer t.Unlock()
	tickerid++
	t.funcs[tickerid] = f
}

func (t *Ticker) Start() {
	go t.Tick()
}

func (t *Ticker) Stop() {
	t.quit <- true
}

func (t *Ticker) Tick() {
	t.ticker = time.NewTicker(time.Second * time.Duration(t.interval))
	for {
		select {
		case <-t.ticker.C:
			for k, f := range t.funcs {
				go func(id int64, tf TickerFunc) {
					t.wg.Add(1)
					tf(id)
					t.wg.Done()
				}(k, f)
			}
		case <-t.quit:
			t.ticker.Stop()
			return
		}
	}
}

func (t *Ticker) Wait() {
	t.wg.Wait()
}
