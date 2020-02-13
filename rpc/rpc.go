package rpc

import "github.com/mitchellh/mapstructure"

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
}

type Call struct {
	Session *Session
	MsgId 	int
	ServicePath string
	ServiceName string
	MsgReq interface{}
	MsgResp interface{}
}

// agent route msg to proxy
type RpcRouter interface {
	Go(call *Call)
}

// proxy return msg to agent
type RpcResponder interface {
	Cb(session *Session, msgId int, msg interface{})
}

// register msg handler as service
type HandlerRegister interface {
	GetRemoteAddrs() string
	RegisterHandler(rcvr interface{}, metadata string) error
}