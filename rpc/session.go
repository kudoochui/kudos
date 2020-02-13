package rpc

import (
	"github.com/kudoochui/kudos/service/idService"
	"sync"
)

type Session struct {
	NodeAddr	string
	SessionId	int64

	UserId 		int64
	userIdLock 	sync.RWMutex

	settings	map[string]string
}

func NewSession(nodeAddr string) *Session  {
	return &Session{
		NodeAddr:  nodeAddr,
		SessionId: idService.GenerateID().Int64(),
		settings:  map[string]string{},
	}
}

func (s *Session) GetSessionId() int64 {
	return s.SessionId
}

func (s *Session) GetUserId() int64 {
	s.userIdLock.RLock()
	defer s.userIdLock.RUnlock()
	return s.UserId
}

func (s *Session) SetUserId(userId int64) {
	s.userIdLock.Lock()
	defer s.userIdLock.Unlock()
	s.UserId = userId
}

func (s *Session) SyncSettings(settings map[string]interface{}) {
	_settings := make(map[string]string)
	for k,v := range settings {
		_settings[k] = v.(string)
	}
	s.settings = _settings
}

func (s *Session) Bind(userId int64) {
	s.UserId = userId

	args := &Args{
		Session: *s,
		MsgReq:  userId,
	}
	reply := &Reply{}
	RpcInvoke(s.NodeAddr, "SessionRemote","Bind", args, reply)
}

func (s *Session) UnBind() {
	s.UserId = 0

	args := &Args{
		Session: *s,
	}
	reply := &Reply{}
	RpcInvoke(s.NodeAddr, "SessionRemote","UnBind", args, reply)
}

func (s *Session) Get(key string) string {
	return s.settings[key]
}

func (s *Session) Set(key, value string) {
	if s.settings == nil {
		s.settings = make(map[string]string)
	}
	s.settings[key] = value
}

func (s *Session) Remove(key string) {
	delete(s.settings, key)
}

func (s *Session) Clone() *Session {
	session := &Session{
		NodeAddr:   s.NodeAddr,
		SessionId:  s.SessionId,
		UserId:     s.UserId,
		settings:   map[string]string{},
	}

	for k,v := range s.settings {
		session.settings[k] = v
	}
	return session
}

// synchronize setting with frontend session
func (s *Session) Push(){
	args := &Args{
		Session: *s,
		MsgReq: s.settings,
	}
	reply := &Reply{}
	RpcInvoke(s.NodeAddr, "SessionRemote", "Push", args, reply)
}

func (s *Session) Close(reason string) {
	args := &Args{
		Session: *s,
		MsgReq: reason,
	}
	reply := &Reply{}
	RpcInvoke(s.NodeAddr, "SessionRemote", "KickBySid", args, reply)
}