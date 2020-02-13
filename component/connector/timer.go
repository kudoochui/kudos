package connector

import (
	"github.com/kudoochui/kudos/utils/timer"
	"time"
)


type Timers struct {
	dispatcher  *timer.Dispatcher
	timerDispatcherLen int
	chanJob		chan *TimeJob
	chanStop	chan *timer.Timer
}

func NewTimer() *Timers {
	return &Timers{
		timerDispatcherLen: 20,
		dispatcher:         timer.NewDispatcher(20),
		chanJob:            make(chan *TimeJob, 5),
		chanStop:           make(chan *timer.Timer, 20),
	}
}

type TimeJob struct {
	timeout 	time.Duration
	cronExpr *timer.CronExpr
	f 		func()
}

func (t *Timers) RunTimer(closeSig chan bool) {
	//t.dispatcher = timer.NewDispatcher(t.timerDispatcherLen)

	for {
		select {
		case <-closeSig:
			return
		case tt := <-t.dispatcher.ChanTimer:
			tt.Cb()
		case job := <-t.chanJob:
			t.work(job)
		case handler := <-t.chanStop:
			handler.Stop()
		}
	}
}
func (t *Timers) work(job *TimeJob) {
	if job.cronExpr != nil {

	} else {

	}
}

func (t *Timers) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	return t.dispatcher.AfterFunc(d, cb)
}

func (t *Timers) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	return t.dispatcher.CronFunc(cronExpr, cb)
}

func (t *Timers) ClearTimeout(handler *timer.Timer){
	t.chanStop <- handler
}