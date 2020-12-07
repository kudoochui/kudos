package connector

import (
	"errors"
	"sync"
	"sync/atomic"
)

// goroutine safe
type SessionMap struct {
	sessions sync.Map
	counter int32
}

func (s *SessionMap) AddSession(a Agent)  {
	s.sessions.Store(a.GetSession().GetSessionId(), a)
	atomic.AddInt32(&s.counter, 1)
}

func (s *SessionMap) DelSession(a Agent) {
	s.sessions.Delete(a.GetSession().GetSessionId())
	atomic.AddInt32(&s.counter, -1)
}

func (s *SessionMap) GetAgent(sessionId int64) (Agent, error) {
	a, ok := s.sessions.Load(sessionId)
	if !ok || a == nil {
		return nil, errors.New("No Sesssion found")
	}
	return a.(Agent), nil
}

func (s *SessionMap) Range(f func(k,v interface{})bool) {
	s.sessions.Range(f)
}

func (s *SessionMap) GetSessionCount() int32 {
	return atomic.LoadInt32(&s.counter)
}

