package timers

import (
	"github.com/kudoochui/kudos/utils/timer"
	"time"
)

type Timers struct {
	opts		*Options
	dispatcher  *timer.Dispatcher
}

func NewTimer(opts ...Option) *Timers {
	options := newOptions(opts...)

	return &Timers{
		opts: options,
	}
}

func (t *Timers) OnInit() {

}

func (t *Timers) OnDestroy() {

}

func (t *Timers) Run(closeSig chan bool) {
	t.dispatcher = timer.NewDispatcher(t.opts.TimerDispatcherLen)

	for {
		select {
		case <-closeSig:
			return
		case tt := <-t.dispatcher.ChanTimer:
			tt.Cb()
		}
	}
}

func (t *Timers) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	if t.opts.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return t.dispatcher.AfterFunc(d, cb)
}

func (t *Timers) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	if t.opts.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return t.dispatcher.CronFunc(cronExpr, cb)
}