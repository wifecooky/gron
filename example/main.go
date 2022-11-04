package main

import (
	"fmt"
	"time"

	"github.com/wifecooky/gron"
)

type printJob struct{ Msg string }

func (p printJob) Run() {
	fmt.Println(p.Msg)
}

func main() {

	var (
		purgeTask = func() { fmt.Println("purge unwanted records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	c := gron.New()

	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("Every 1 hour")
	})
	c.Start()

	c.AddFunc(gron.WEEKLY, func() { fmt.Println("Every week") })
	c.Add(gron.DAILY.At("12:30"), printFoo)
	c.Start()

	// Jobs may also be added to a running Cron
	c.Add(gron.MONTHLY, printBar)
	c.AddFunc(gron.YEARLY, purgeTask)

	// Stop the scheduler (does not stop any jobs already running).
	defer c.Stop()
}
