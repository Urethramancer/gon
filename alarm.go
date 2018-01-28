package gon

import (
	"sync"
	"time"
)

type Alarm struct {
	sync.RWMutex
	scheduler *Scheduler
	wg        sync.WaitGroup
	delay     time.Duration
	timer     *time.Timer
	id        int64
	f         AlarmFunc
	quit      chan bool
}

// AlarmFunc is the signature of the callbacks run when the timer triggers.
type AlarmFunc func(int64)

// NewAlarm
func NewAlarm(d time.Duration, aid int64, af AlarmFunc) *Alarm {
	a := &Alarm{
		delay: d,
		id:    aid,
		f:     af,
	}
	a.quit = make(chan bool, 0)
	return a
}

func (a *Alarm) Start() {
	a.timer = time.NewTimer(a.delay)
	for {
		select {
		case <-a.timer.C:
			a.RLock()
			go func(id int64, af AlarmFunc) {
				a.wg.Add(1)
				af(id)
				a.wg.Done()
			}(a.id, a.f)
			a.RUnlock()
			a.wg.Wait()
			a.scheduler.RemoveAlarm(a.id)
			return
		case <-a.quit:
			a.timer.Stop()
			return
		}
	}
}

// Stop and remove the alarm.
func (a *Alarm) Stop() {
	a.quit <- true
	a.scheduler.RemoveAlarm(a.id)
}

// Wait for the event to finish.
func (a *Alarm) Wait() {
	a.wg.Wait()
}
