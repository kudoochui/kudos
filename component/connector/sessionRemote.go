package connector

import (
	"context"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/protocol/pkg"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
	"runtime"
)

type SessionRemote struct {
	connector	*Connector
}

func NewSessionRemote(c *Connector) *SessionRemote {
	return &SessionRemote{
		connector: c,
	}
}

func (s *SessionRemote) Bind(ctx context.Context, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := args.Session.GetSessionId()
	agent,err := s.connector.sessions.GetAgent(sessioinId)
	if err != nil {
		log.Error("Bind can't find session:%s", sessioinId)
	}
	agent.GetSession().SetUserId(args.MsgReq.(int64))
	log.Debug("Bind success: %d", agent.GetSession().GetUserId())
	return nil
}

func (s *SessionRemote) UnBind(ctx context.Context, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := args.Session.GetSessionId()
	agent,err := s.connector.sessions.GetAgent(sessioinId)
	if err != nil {
		log.Error("UnBind can't find session:%s", sessioinId)
	}
	agent.GetSession().SetUserId(0)
	log.Debug("UnBind success: %d", agent.GetSession().GetUserId())
	return nil
}

func (s *SessionRemote) Push(ctx context.Context, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := args.Session.GetSessionId()
	agent,err := s.connector.sessions.GetAgent(sessioinId)
	if err != nil {
		log.Error("Push can't find session:%s", sessioinId)
	}
	settings := args.MsgReq.(map[string]interface{})
	agent.GetSession().SyncSettings(settings)
	log.Debug("Push success: %v", settings)
	return nil
}

func (s *SessionRemote) KickBySid(ctx context.Context, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := args.Session.GetSessionId()
	agent,err := s.connector.sessions.GetAgent(sessioinId)
	if err != nil {
		log.Error("KickBySid can't find session:%s", sessioinId)
		return err
	}
	reason := args.MsgReq.(string)
	ret := map[string]string{
		"reason": reason,
	}
	buffer,_ := codecService.GetCodecService().Marshal(ret)
	agent.Write(pkg.Encode(pkg.TYPE_KICK, buffer)...)

	runtime.Gosched()
	agent.Close()
	return nil
}