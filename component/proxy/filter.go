package proxy

import "github.com/kudoochui/kudos/rpc"

type filter interface {
	Before(route string, msgId int, session *rpc.Session, msgReq interface{})
	After(route string, msgId int, session *rpc.Session, msgResp interface{})
}
