package connector

import (
	"errors"
	"sync"
	"sync/atomic"
)

// goroutine safe
type sessionMap struct {
	sessions sync.Map
	counter int32
}

func (s *sessionMap) AddSession(a *agent)  {
	s.sessions.Store(a.session.GetSessionId(), a)
	atomic.AddInt32(&s.counter, 1)
}

func (s *sessionMap) DelSession(a *agent) {
	s.sessions.Delete(a.session.GetSessionId())
	atomic.AddInt32(&s.counter, -1)
}

func (s *sessionMap) GetAgent(sessionId int64) (*agent, error) {
	a, ok := s.sessions.Load(sessionId)
	if !ok || a == nil {
		return nil, errors.New("No Sesssion found")
	}
	return a.(*agent), nil
}

func (s *sessionMap) Range(f func(k,v interface{})bool) {
	s.sessions.Range(f)
}

func (s *sessionMap) GetSessionCount() int32 {
	return atomic.LoadInt32(&s.counter)
}

