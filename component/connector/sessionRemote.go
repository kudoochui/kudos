package connector

import (
	"context"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"runtime"
)

type SessionRemote struct {
	connector	Connector
}

func NewSessionRemote(c Connector) *SessionRemote {
	return &SessionRemote{
		connector: c,
	}
}

func (s *SessionRemote) Bind(ctx context.Context, session *rpc.Session, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := session.GetSessionId()
	agent,err := s.connector.GetSessionMap().GetAgent(sessioinId)
	if err != nil {
		log.Error("Bind can't find session:%s", sessioinId)
		return nil
	}
	agent.GetSession().SetUserId(args.MsgReq.(int64))
	log.Debug("Bind success: %d", agent.GetSession().GetUserId())
	return nil
}

func (s *SessionRemote) UnBind(ctx context.Context, session *rpc.Session, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := session.GetSessionId()
	agent,err := s.connector.GetSessionMap().GetAgent(sessioinId)
	if err != nil {
		log.Error("UnBind can't find session:%s", sessioinId)
		return nil
	}
	agent.GetSession().SetUserId(0)
	log.Debug("UnBind success: %d", agent.GetSession().GetUserId())
	return nil
}

func (s *SessionRemote) Push(ctx context.Context, session *rpc.Session, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := session.GetSessionId()
	agent,err := s.connector.GetSessionMap().GetAgent(sessioinId)
	if err != nil {
		log.Error("Push can't find session:%s", sessioinId)
		return nil
	}
	settings := args.MsgReq.(map[string]interface{})
	agent.GetSession().SyncSettings(settings)
	log.Debug("Push success: %v", settings)
	return nil
}

func (s *SessionRemote) KickBySid(ctx context.Context, session *rpc.Session, args *rpc.Args, reply *rpc.Reply) error {
	sessioinId := session.GetSessionId()
	agent,err := s.connector.GetSessionMap().GetAgent(sessioinId)
	if err != nil {
		log.Error("KickBySid can't find session:%s", sessioinId)
		return err
	}
	reason := args.MsgReq.(string)
	agent.KickMessage(reason)

	runtime.Gosched()
	agent.Close()
	return nil
}

func (s *SessionRemote) GetSessionCount(ctx context.Context, session *rpc.Session, args *rpc.Args, reply *rpc.Reply) error {
	count := s.connector.GetSessionMap().GetSessionCount()
	reply.MsgResp = count
	return nil
}