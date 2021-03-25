package sessionService

import (
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/rpcClientService"
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

func (s *SessionService) KickBySid(nodeId string, sid int64, reason string) {
	args := &rpc.Args{
		Session: rpc.Session{SessionId:sid},
		MsgReq: reason,
	}
	reply := &rpc.Reply{}
	rpcClientService.GetRpcClientService().Call(nodeId, "SessionRemote","KickBySid", args, reply)
}