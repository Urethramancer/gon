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
	f         EventFunc
	quit      chan bool
	repeat    bool
}

// NewAlarm creates the Alarm structure and quit channel.
func NewAlarm(d time.Duration, aid int64, af EventFunc) *Alarm {
	a := &Alarm{
		delay: d,
		id:    aid,
		f:     af,
	}
	a.quit = make(chan bool)
	return a
}

// Start creates a one-shot timer and runs it after its duration has passed.
// If the alarm is a repeating alarm, a 24-hour Ticker is created.
func (a *Alarm) Start() {
	a.timer = time.NewTimer(a.delay)
	for {
		select {
		case <-a.timer.C:
			a.RLock()
			a.wg.Add(1)
			go func(id int64, af EventFunc) {
				af(id)
				a.wg.Done()
			}(a.id, a.f)
			a.RUnlock()
			a.wg.Wait()
			if a.repeat {
				// The alarm transitions to repeating ticker here
				a.scheduler.RepeatHours(24, a.f)
			}
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
