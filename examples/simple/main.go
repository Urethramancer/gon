package main

import (
	"fmt"
	"time"

	"github.com/Urethramancer/gon"
)

func main() {
	sc := gon.NewScheduler()
	fmt.Printf("Adding alarm in 10 seconds.\n")
	when := time.Now().Add(time.Second * 10)
	sc.AddAlarmAt(when, event)
	fmt.Printf("Adding tickers at intervals 3, 5, 10 and 15 seconds.\n")
	sc.RepeatSeconds(3, tick)
	sc.RepeatSeconds(5, tick)
	sc.RepeatSeconds(10, tick)
	sc.RepeatSeconds(15, tick)
	fmt.Printf("Waiting…\n")
	time.Sleep(time.Second * 22)
	sc.Wait()
}

func tick(id int64) {
	fmt.Printf("Ticker #%d running\n", id)
}

func event(id int64) {
	fmt.Printf("Alarm #%d triggered\n", id)
}
