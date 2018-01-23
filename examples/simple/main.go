package main

import (
	"fmt"
	"time"

	"github.com/Urethramancer/gon"
)

func main() {
	sc := gon.NewScheduler()
	sc.AddEveryNSecond(Tick, 3)
	sc.AddEveryNSecond(Tick, 5)
	sc.AddEveryNSecond(Tick, 10)
	sc.AddEveryNSecond(Tick, 15)
	fmt.Printf("Waiting…\n")
	time.Sleep(time.Second * 22)
	sc.Wait()
}

func Tick(id int64) {
	fmt.Printf("Ticker #%d running\n", id)
}
