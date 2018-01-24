package gon

import (
	"sync"
	"time"
)

var alarmid int64

type Alarm struct {
	sync.RWMutex
	wg    sync.WaitGroup
	delay time.Duration
	timer *time.Timer
	quit  chan bool
	funcs map[int64]AlarmFunc
}

// AlarmFunc is the signature of the callbacks run when the timer triggers.
type AlarmFunc func(int64)

// NewAlarm
func NewAlarm(d time.Duration) *Alarm {
	a := &Alarm{
		delay: d,
		funcs: make(map[int64]AlarmFunc)}
	a.quit = make(chan bool, 0)
	return a
}

// AddFunc adds anothe callback to the slice with a new ID.
func (a *Alarm) AddFunc(f AlarmFunc) {
	a.Lock()
	defer a.Unlock()
	alarmid++
	a.funcs[alarmid] = f
}

func (a *Alarm) Start() {
	a.timer = time.NewTimer(a.delay)
	for {
		select {
		case <-a.timer.C:
			a.RLock()
			for k, f := range a.funcs {
				go func(id int64, af AlarmFunc) {
					a.wg.Add(1)
					af(id)
					a.wg.Done()
				}(k, f)
			}
			a.RUnlock()
			a.wg.Wait()
			return
		case <-a.quit:
			a.timer.Stop()
			return
		}
	}
}

func (a *Alarm) Stop() {
	a.quit <- true
}
