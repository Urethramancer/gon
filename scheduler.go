package gon

import (
	"sync"
	"time"
)

// Scheduler holds pointers to all the tickers and timers.
type Scheduler struct {
	sync.RWMutex
	tickerid       int64
	alarmid        int64
	tickers        map[time.Duration]*Ticker
	dormantTickers map[time.Duration]*Ticker
	alarms         map[int64]*Alarm
}

// EventFunc is the signature of the ticker/alarm callbacks.
type EventFunc func(int64)

// NewScheduler returns a Scheduler populated with maps.
func NewScheduler() *Scheduler {
	return &Scheduler{
		tickers: make(map[time.Duration]*Ticker),
		alarms:  make(map[int64]*Alarm),
	}
}

//
// Tickers
// Repeating events
//

func (sc *Scheduler) addTicker(d time.Duration, f EventFunc) {
	sc.Lock()
	defer sc.Unlock()
	t, ok := sc.tickers[d]
	if !ok {
		t, ok = sc.dormantTickers[d]
		if ok {
			delete(sc.dormantTickers, d)
			sc.tickers[d] = t
		} else {
			t = NewTicker(d)
			sc.tickers[d] = t
		}
	}
	sc.tickerid++
	t.AddFunc(f, sc.tickerid)
	go t.Start()
}

// RemoveTicker removes a ticker by duration, stopping it if necessary.
func (sc *Scheduler) RemoveTicker(d time.Duration) {
	sc.Lock()
	defer sc.Unlock()
	t, ok := sc.tickers[d]
	if ok {
		t.Stop()
		delete(sc.tickers, d)
		sc.dormantTickers[d] = t
	}
}

// RepeatSeconds adds a repeating task based on a seconds interval.
func (sc *Scheduler) RepeatSeconds(n int, f EventFunc) {
	sc.addTicker(time.Second*time.Duration(n), f)
}

// RepeatMinutes adds a repeating task on a minute-based interval.
func (sc *Scheduler) RepeatMinutes(n int, f EventFunc) {
	sc.addTicker(time.Minute*time.Duration(n), f)
}

// RepeatHours adds a repeating task on an hour-based interval.
func (sc *Scheduler) RepeatHours(n int, f EventFunc) {
	sc.addTicker(time.Hour*time.Duration(n), f)
}

//
// Alarms
// One-time events
//

func (sc *Scheduler) addAlarm(d time.Duration, f EventFunc, repeat bool) {
	sc.Lock()
	defer sc.Unlock()
	sc.alarmid++
	alarm := NewAlarm(d, sc.alarmid, f)
	alarm.scheduler = sc
	alarm.repeat = repeat
	sc.alarms[sc.alarmid] = alarm
	go alarm.Start()
}

// RemoveAlarm removes an alarm by id, stopping it if necessary.
func (sc *Scheduler) RemoveAlarm(id int64) {
	sc.Lock()
	defer sc.Unlock()
	alarm, ok := sc.alarms[id]
	if ok {
		alarm.Stop()
		delete(sc.alarms, id)
	}
}

// AddAlarmIn triggers functions after a specific duration has passed.
func (sc *Scheduler) AddAlarmIn(d time.Duration, f EventFunc) {
	sc.addAlarm(d, f, false)
}

// AddAlarmAt triggers functions at a specific time of day.
func (sc *Scheduler) AddAlarmAt(t time.Time, f EventFunc, repeat bool) {
	when := time.Until(t)
	sc.addAlarm(when, f, repeat)
}

// Wait for all waitgroups in tickers and alarms.
func (sc *Scheduler) Wait() {
	for _, t := range sc.tickers {
		t.Wait()
	}
	for _, a := range sc.alarms {
		a.Wait()
	}
}
