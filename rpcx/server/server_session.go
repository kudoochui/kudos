package server

import (
	"context"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/rpcx/protocol"
	"github.com/kudoochui/kudos/rpcx/share"
	"net"
	"sync/atomic"
	"time"
)

type ServerSession struct {
	NodeId	string
	SessionId	int64
	UserId 		int64

	Settings	map[string]string
	server 		*Server
	conn 		net.Conn
}

func NewSessionFromRpc(nodeId string, sessionId int64, userId int64) *ServerSession  {
	return &ServerSession{
		NodeId: nodeId,
		SessionId: sessionId,
		UserId: userId,
		Settings:  map[string]string{},
	}
}

func (s *ServerSession) GetNodeId() string {
	return s.NodeId
}

func (s *ServerSession) GetSessionId() int64 {
	return s.SessionId
}

func (s *ServerSession) GetUserId() int64 {
	return s.UserId
}

func (s *ServerSession) SetUserId(userId int64) {
	s.UserId = userId
}

func (s *ServerSession) SyncSettings(settings map[string]interface{}) {
	_settings := make(map[string]string)
	for k,v := range settings {
		_settings[k] = v.(string)
	}
	s.Settings = _settings
}

func (s *ServerSession) Bind(userId int64) {
	s.UserId = userId

	args := &rpc.Args{
		MsgReq:  userId,
	}

	s.sendMessage("SessionRemote","Bind", nil, args)
}

func (s *ServerSession) UnBind() {
	s.UserId = 0

	args := &rpc.Args{
	}

	s.sendMessage("SessionRemote","UnBind", nil, args)
}

func (s *ServerSession) Get(key string) string {
	return s.Settings[key]
}

func (s *ServerSession) Set(key, value string) {
	if s.Settings == nil {
		s.Settings = make(map[string]string)
	}
	s.Settings[key] = value
}

func (s *ServerSession) Remove(key string) {
	delete(s.Settings, key)
}

func (s *ServerSession) Clone() *ServerSession {
	session := &ServerSession{
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
func (s *ServerSession) Push(){
	args := &rpc.Args{
		MsgReq: s.Settings,
	}
	s.sendMessage("SessionRemote","Push", nil, args)
}

func (s *ServerSession) Close(reason string) {
	args := &rpc.Args{
		MsgReq: reason,
	}

	s.sendMessage("SessionRemote","KickBySid", nil, args)
}

// Server push message
func (s *ServerSession) PushMessage(servicePath string, serviceMethod string, args interface{}) {
	s.sendMessage(servicePath, serviceMethod, nil, args)
}

func (s *ServerSession) sendMessage(servicePath, serviceMethod string, metadata map[string]string, args interface{}) error {
	ctx := share.WithValue(context.Background(), StartSendRequestContextKey, time.Now().UnixNano())
	s.server.Plugins.DoPreWriteRequest(ctx)

	req := protocol.GetPooledMsg()
	req.SetMessageType(protocol.Request)

	seq := atomic.AddUint64(&s.server.seq, 1)
	req.SetSeq(seq)
	req.SetOneway(true)
	req.SetSerializeType(protocol.MsgPack)
	req.ServicePath = servicePath
	req.ServiceMethod = serviceMethod
	req.NodeId = s.NodeId
	req.SessionId = s.SessionId
	req.UserId = s.UserId
	req.Metadata = metadata

	codec := share.Codecs[protocol.MsgPack]
	data, err := codec.Encode(args)
	if err != nil {
		return err
	}
	req.Payload = data

	b := req.EncodeSlicePointer()
	_, err = s.conn.Write(*b)
	protocol.PutData(b)

	s.server.Plugins.DoPostWriteRequest(ctx, req, err)
	protocol.FreeMsg(req)
	return err
}