package rpc

import (
	"github.com/mitchellh/mapstructure"
)

type Args struct {
	Session Session
	MsgId int
	MsgReq interface{}
}

func (a *Args) GetObject(t interface{}) error {
	return mapstructure.Decode(a.MsgReq, t)
}

type Reply struct {
	Code 	int
	ErrMsg 	string
	MsgResp interface{}
}

type Call struct {
	Session *Session
	MsgId 	int
	ServicePath string
	ServiceName string
	MsgReq interface{}
	MsgResp interface{}
	Done 	interface{}
}

// agent route msg to proxy
type RpcRouter interface {
	Go(call *Call)
}

// Route msg to the specified node
type CustomerRoute func(session *Session, servicePath, serviceName string) (string, error)

// proxy return msg to agent
type RpcResponder interface {
	Cb(session *Session, msgId int, msg interface{})
}

// register msg handler as service
type HandlerRegister interface {
	GetRemoteAddrs() string
	RegisterHandler(rcvr interface{}, metadata string) error
	RegisterName(name string, rcvr interface{}, metadata string) error
}