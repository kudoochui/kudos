package timer

import (
	"fmt"
	"time"
)

func ExampleTimer() {
	d := NewDispatcher(10)

	// timer 1
	d.AfterFunc(1, func() {
		fmt.Println("My name is Leaf")
	})

	// timer 2
	t := d.AfterFunc(1, func() {
		fmt.Println("will not print")
	})
	t.Stop()

	// dispatch
	(<-d.ChanTimer).Cb()

	// Output:
	// My name is Leaf
}

func ExampleCronExpr() {
	cronExpr, err := NewCronExpr("0 * * * *")
	if err != nil {
		return
	}

	fmt.Println(cronExpr.Next(time.Date(
		2000, 1, 1,
		20, 10, 5,
		0, time.UTC,
	)))

	// Output:
	// 2000-01-01 21:00:00 +0000 UTC
}

func ExampleCron() {
	d := NewDispatcher(10)

	// cron expr
	cronExpr, err := NewCronExpr("* * * * * *")
	if err != nil {
		return
	}

	// cron
	var c *Cron
	c = d.CronFunc(cronExpr, func() {
		fmt.Println("My name is Leaf")
		c.Stop()
	})

	// dispatch
	(<-d.ChanTimer).Cb()

	// Output:
	// My name is Leaf
}
