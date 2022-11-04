# gron

> This repository is forked from [roylee0704/gron](https://github.com/roylee0704/gron)
> Thanks for the great work!

Gron provides a clear syntax for writing and deploying cron jobs.

## Goals

- Minimalist APIs for scheduling jobs.
- Thread safety.
- Customizable Job Type.
- Customizable Schedule.

## Installation

```sh
$ go get github.com/wifecooky/gron
```

## Usage

Create `schedule.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/wifecooky/gron"
)

func main() {
	c := gron.New()
	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("runs every hour.")
	})
	c.Start()
}
```

#### Schedule Parameters

All scheduling is done in the machine's local time zone (as provided by the Go [time package](http://www.golang.org/pkg/time)).


Setup basic periodic schedule with `gron.Every()`.

```go
gron.Every(1*time.Second)
gron.Every(1*time.Minute)
gron.Every(1*time.Hour)
```

Also provide some pre-defined schedules.

```go
gron.DAILY
gron.WEEKLY
gron.MONTHLY
gron.YEARLY
```

Schedule to run at specific time with `.At(hh:mm)`
```go
gron.Every(30 * xtime.Day).At("00:00")
gron.Every(1 * xtime.Week).At("23:59")
```

#### Custom Job Type
You may define custom job types by implementing `gron.Job` interface: `Run()`.

For example:

```go
type Reminder struct {
	Msg string
}

func (r Reminder) Run() {
  fmt.Println(r.Msg)
}
```

After job has defined, instantiate it and schedule to run in Gron.
```go
c := gron.New()
r := Reminder{ "Feed the baby!" }
c.Add(gron.Every(8*time.Hour), r)
c.Start()
```

#### Custom Job Func
You may register `Funcs` to be executed on a given schedule. Gron will run them in their own goroutines, asynchronously.

```go
c := gron.New()
c.AddFunc(gron.Every(1*time.Second), func() {
	fmt.Println("runs every second")
})
c.Start()
```


#### Custom Schedule
Schedule is the interface that wraps the basic `Next` method: `Next(p time.Duration) time.Time`

In `gron`, the interface value `Schedule` has the following concrete types:

- **periodicSchedule**. adds time instant t to underlying period p.
- **atSchedule**. reoccurs every period p, at time components(hh:mm).

For more info, checkout `schedule.go`.

### Full Example

```go
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
```
