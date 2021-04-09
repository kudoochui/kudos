package connector

import (
	"context"
	"github.com/kudoochui/kudos/rpc"
	msgService "github.com/kudoochui/kudos/service/msgService"
)

type ChannelRemote struct {
	connector	Connector
}

func NewChannelRemote(conn Connector) *ChannelRemote {
	return &ChannelRemote{connector:conn}
}

// Push message to client by uids.
func (c *ChannelRemote) PushMessage(ctx context.Context, session *rpc.Session, args *rpc.ArgsGroup, reply *rpc.ReplyGroup) error {
	//log.Debug(">>>> %s push: %s, %+v, %s", c.connector.opts.WSAddr, args.Route, args.Sids, string(args.Payload))
	sessioinId := session.GetSessionId()
	if a, err := c.connector.GetSessionMap().GetAgent(sessioinId); err == nil {
		routeId := msgService.GetMsgService().GetRouteId(args.Route)
		a.PushMessage(routeId, args.Payload)
	}
	return nil
}