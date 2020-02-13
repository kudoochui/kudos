package connector

import (
	"errors"
	"sync"
)

// goroutine safe
type sessionMap struct {
	sessions sync.Map
}

func (s *sessionMap) AddSession(a *agent)  {
	s.sessions.Store(a.session.GetSessionId(), a)
}

func (s *sessionMap) DelSession(a *agent) {
	s.sessions.Delete(a.session.GetSessionId())
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

