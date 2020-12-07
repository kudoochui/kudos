package rpc

import (
	"github.com/kudoochui/kudos/service/idService"
	"github.com/kudoochui/kudos/service/rpcClientService"
	"sync"
)

type Session struct {
	NodeId	string
	SessionId	int64

	UserId 		int64
	userIdLock 	sync.RWMutex

	Settings	map[string]string
}

func NewSession(nodeId string) *Session  {
	return &Session{
		NodeId: nodeId,
		SessionId: idService.GenerateID().Int64(),
		Settings:  map[string]string{},
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
	s.Settings = _settings
}

func (s *Session) Bind(userId int64) {
	s.UserId = userId

	args := &Args{
		Session: *s,
		MsgReq:  userId,
	}
	reply := &Reply{}
	rpcClientService.GetRpcClientService().Call(s.NodeId+"@SessionRemote","Bind", args, reply)
}

func (s *Session) UnBind() {
	s.UserId = 0

	args := &Args{
		Session: *s,
	}
	reply := &Reply{}
	rpcClientService.GetRpcClientService().Call(s.NodeId+"@SessionRemote","UnBind", args, reply)
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

func (s *Session) Clone() *Session {
	session := &Session{
		NodeId:   s.NodeId,
		SessionId:  s.SessionId,
		UserId:     s.UserId,
		Settings:   map[string]string{},
	}

	for k,v := range s.Settings {
		session.Settings[k] = v
	}
	return session
}

// synchronize setting with frontend session
func (s *Session) Push(){
	args := &Args{
		Session: *s,
		MsgReq: s.Settings,
	}
	reply := &Reply{}
	rpcClientService.GetRpcClientService().Call(s.NodeId+"@SessionRemote","Push", args, reply)
}

func (s *Session) Close(reason string) {
	args := &Args{
		Session: *s,
		MsgReq: reason,
	}
	reply := &Reply{}
	rpcClientService.GetRpcClientService().Call(s.NodeId+"@SessionRemote","KickBySid", args, reply)
}