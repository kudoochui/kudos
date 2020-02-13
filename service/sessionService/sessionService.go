package sessionService

import (
	"github.com/kudoochui/kudos/rpc"
	"sync"
)

var service *SessionService
var once sync.Once

func GetSessionService() *SessionService {
	once.Do(func() {
		service = &SessionService{
		}
	})
	return service
}

type SessionService struct {

}

func (s *SessionService) KickBySid(nodeAddr string, sid int64, reason string) {
	args := &rpc.Args{
		Session: rpc.Session{SessionId:sid},
		MsgReq: reason,
	}
	reply := &rpc.Reply{}
	rpc.RpcInvoke(nodeAddr, "SessionRemote", "KickBySid", args, reply)
}