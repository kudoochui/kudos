package mailbox

import (
	"github.com/kudoochui/kudos/rpcx/mailbox/queue/goring"
	"github.com/kudoochui/kudos/rpcx/mailbox/queue/mpsc"
)

type unboundedMailboxQueue struct {
	userMailbox *goring.Queue
}

func (q *unboundedMailboxQueue) Push(m interface{}) {
	q.userMailbox.Push(m)
}

func (q *unboundedMailboxQueue) Pop() interface{} {
	m, o := q.userMailbox.Pop()
	if o {
		return m
	}
	return nil
}

// Unbounded creates an unbounded mailbox
func Unbounded() Mailbox {
	q := &unboundedMailboxQueue{
		userMailbox: goring.New(10),
	}
	return &defaultMailbox{
		systemMailbox: mpsc.New(),
		userMailbox:   q,
	}
}