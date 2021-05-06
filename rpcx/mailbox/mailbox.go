package mailbox

import (
	"github.com/kudoochui/kudos/rpcx/log"
	"github.com/kudoochui/kudos/rpcx/mailbox/queue/mpsc"
	"runtime"
	"sync/atomic"
	"time"
)

// MessageInvoker is the interface used by a mailbox to forward messages for processing
type MessageInvoker interface {
	InvokeSystemMessage(interface{})
	InvokeUserMessage(interface{})
	OnTimeTick(interface{})
}

// Mailbox interface is used to enqueue messages to the mailbox
type Mailbox interface {
	PostUserMessage(message interface{})
	PostSystemMessage(message interface{})
	PostTimeMessage(message interface{})
	RegisterHandlers(invoker MessageInvoker, dispatcher Dispatcher)
}

const (
	idle int32 = iota
	running
)

type defaultMailbox struct {
	userMailbox     queue
	systemMailbox   *mpsc.Queue
	timeMailbox 	interface{}				//only one message
	schedulerStatus int32
	userMessages    int32
	sysMessages     int32
	suspended       int32
	invoker         MessageInvoker
	dispatcher      Dispatcher
}

func (m *defaultMailbox) PostUserMessage(message interface{}) {
	m.userMailbox.Push(message)
	atomic.AddInt32(&m.userMessages, 1)
	m.schedule()
}

func (m *defaultMailbox) PostSystemMessage(message interface{}) {
	m.systemMailbox.Push(message)
	atomic.AddInt32(&m.sysMessages, 1)
	m.schedule()
}

func (m *defaultMailbox) PostTimeMessage(message interface{}) {
	m.timeMailbox = message
}

func (m *defaultMailbox) RegisterHandlers(invoker MessageInvoker, dispatcher Dispatcher) {
	m.invoker = invoker
	m.dispatcher = dispatcher
}

func (m *defaultMailbox) schedule() {
	if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
		m.dispatcher.Schedule(m.processMessages)
	}
}

func (m *defaultMailbox) processMessages() {
process:
	m.run()

	// set mailbox to idle
	atomic.StoreInt32(&m.schedulerStatus, idle)
	sys := atomic.LoadInt32(&m.sysMessages)
	user := atomic.LoadInt32(&m.userMessages)
	// check if there are still messages to process (sent after the message loop ended)
	if sys > 0 || (atomic.LoadInt32(&m.suspended) == 0 && user > 0) {
		// try setting the mailbox back to running
		if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
			//	fmt.Printf("looping %v %v %v\n", sys, user, m.suspended)
			goto process
		}
	}
}

func (m *defaultMailbox) run() {
	var msg interface{}

	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			l := runtime.Stack(buf, false)
			log.Errorf("Mailbox Recovering %v: %s", r, buf[:l])
		}
	}()

	tick := time.Tick(time.Millisecond * 100)
	i, t := 0, m.dispatcher.Throughput()
	for {
		select {
		case <-tick:
			m.invoker.OnTimeTick(m.timeMailbox)
		default:
		}

		if i > t {
			i = 0
			runtime.Gosched()
		}

		i++

		// keep processing system messages until queue is empty
		if msg = m.systemMailbox.Pop(); msg != nil {
			atomic.AddInt32(&m.sysMessages, -1)
			switch msg.(type) {
			case *SuspendMailbox:
				atomic.StoreInt32(&m.suspended, 1)
			case *ResumeMailbox:
				atomic.StoreInt32(&m.suspended, 0)
			default:
				m.invoker.InvokeSystemMessage(msg)
			}
			continue
		}

		// didn't process a system message, so break until we are resumed
		if atomic.LoadInt32(&m.suspended) == 1 {
			return
		}

		if msg = m.userMailbox.Pop(); msg != nil {
			atomic.AddInt32(&m.userMessages, -1)
			m.invoker.InvokeUserMessage(msg)
		} else {
			return
		}
	}

}