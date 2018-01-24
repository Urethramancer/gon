package main

import (
	"fmt"
	"time"

	"github.com/Urethramancer/gon"
)

func main() {
	sc := gon.NewScheduler()
	fmt.Printf("Adding tickers at intervals 3, 5, 10 and 15 seconds.\n")
	sc.RepeatSeconds(3, Tick)
	sc.RepeatSeconds(5, Tick)
	sc.RepeatSeconds(10, Tick)
	sc.RepeatSeconds(15, Tick)
	fmt.Printf("Waiting…\n")
	time.Sleep(time.Second * 22)
	sc.Wait()
}

func Tick(id int64) {
	fmt.Printf("Ticker #%d running\n", id)
}
