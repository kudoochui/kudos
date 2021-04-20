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
	nodeId	string
	sessionId	int64
	userId 		int64

	settings	map[string]string
	server 		*Server
	conn 		net.Conn
	agent 		*ConnAgent
}

func NewSessionFromRpc(nodeId string, sessionId int64, userId int64, agent *ConnAgent) *ServerSession  {
	return &ServerSession{
		nodeId: nodeId,
		sessionId: sessionId,
		userId: userId,
		settings:  map[string]string{},
		agent: agent,
	}
}

func (s *ServerSession) GetConnectionAgent() *ConnAgent {
	return s.agent
}

func (s *ServerSession) GetNodeId() string {
	return s.nodeId
}

func (s *ServerSession) GetSessionId() int64 {
	return s.sessionId
}

func (s *ServerSession) GetUserId() int64 {
	return s.userId
}

func (s *ServerSession) SetUserId(userId int64) {
	s.userId = userId
}

func (s *ServerSession) SyncSettings(settings map[string]interface{}) {
	_settings := make(map[string]string)
	for k,v := range settings {
		_settings[k] = v.(string)
	}
	s.settings = _settings
}

func (s *ServerSession) Bind(userId int64) {
	s.userId = userId

	args := &rpc.Args{
		MsgReq:  userId,
	}

	s.sendMessage("SessionRemote","Bind", nil, args)
}

func (s *ServerSession) UnBind() {
	s.userId = 0

	args := &rpc.Args{
	}

	s.sendMessage("SessionRemote","UnBind", nil, args)
}

func (s *ServerSession) Get(key string) string {
	return s.settings[key]
}

func (s *ServerSession) Set(key, value string) {
	if s.settings == nil {
		s.settings = make(map[string]string)
	}
	s.settings[key] = value
}

func (s *ServerSession) Remove(key string) {
	delete(s.settings, key)
}

func (s *ServerSession) Clone() *ServerSession {
	session := &ServerSession{
		nodeId:   s.nodeId,
		sessionId:  s.sessionId,
		userId:     s.userId,
		settings:   map[string]string{},
	}

	for k,v := range s.settings {
		session.settings[k] = v
	}
	return session
}

// synchronize setting with frontend session
func (s *ServerSession) Push(){
	args := &rpc.Args{
		MsgReq: s.settings,
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
	req.NodeId = s.nodeId
	req.SessionId = s.sessionId
	req.UserId = s.userId
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