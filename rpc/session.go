package rpc

import (
	"github.com/kudoochui/kudos/service/idService"
	"sync"
)

type Session struct {
	NodeId	string
	SessionId	int64

	UserId 		int64
	mu 	sync.RWMutex
	Settings	map[string]string
	cachedSettings map[string]string
}

func NewSession(nodeId string) *Session  {
	return &Session{
		NodeId: nodeId,
		SessionId: idService.GenerateID().Int64(),
		Settings:  map[string]string{},
		cachedSettings:  map[string]string{},
	}
}

func NewSessionFromRpc(nodeId string, sessionId int64, userId int64) *Session  {
	return &Session{
		NodeId: nodeId,
		SessionId: sessionId,
		UserId: userId,
		Settings:  map[string]string{},
		cachedSettings:  map[string]string{},
	}
}

func (s *Session) GetNodeId() string {
	return s.NodeId
}

func (s *Session) GetSessionId() int64 {
	return s.SessionId
}

func (s *Session) GetUserId() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.UserId
}

func (s *Session) SetUserId(userId int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UserId = userId
}

func (s *Session) SyncSettings(settings map[string]interface{}) {
	_settings := make(map[string]string)
	for k,v := range settings {
		_settings[k] = v.(string)
	}
	s.Settings = _settings
}

func (s *Session) Get(key string) string {
	return s.Settings[key]
}

func (s *Session) Set(key, value string) {
	if s.Settings == nil {
		s.Settings = make(map[string]string)
	}
	s.Settings[key] = value
}

func (s *Session) Remove(key string) {
	delete(s.Settings, key)
}

func (s *Session) GetCache(key string) string {
	return s.cachedSettings[key]
}

func (s *Session) SetCache(key, value string) {
	if s.cachedSettings == nil {
		s.cachedSettings = make(map[string]string)
	}
	s.cachedSettings[key] = value
}

func (s *Session) RemoveCache(key string) {
	delete(s.cachedSettings, key)
}

func (s *Session) Clone() *Session {
	session := &Session{
		NodeId:   s.NodeId,
		SessionId:  s.SessionId,
		UserId:     s.UserId,
		Settings:   map[string]string{},
		cachedSettings: map[string]string{},
	}

	for k,v := range s.Settings {
		session.Settings[k] = v
	}

	for k,v := range s.cachedSettings {
		session.cachedSettings[k] = v
	}
	return session
}