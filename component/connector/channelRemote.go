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
func (c *ChannelRemote) PushMessageByGroup(ctx context.Context, args *rpc.ArgsGroup, reply *rpc.ReplyGroup) error {
	//log.Debug(">>>> %s push: %s, %+v, %s", c.connector.opts.WSAddr, args.Route, args.Sids, string(args.Payload))
	for _, sid := range args.Sids {
		if a, err := c.connector.GetSessionMap().GetAgent(sid); err == nil {
			routeId := msgService.GetMsgService().GetRouteId(args.Route)
			a.PushMessage(routeId, args.Payload)
		}
	}
	return nil
}

// Broadcast to all the client connectd with current frontend server.
func (c *ChannelRemote) Broadcast(ctx context.Context, args *rpc.ArgsGroup, reply *rpc.ReplyGroup) error {
	//log.Debug(">>>> %s broadcast: %s, %s", c.connector.opts.WSAddr, args.Route, string(args.Payload))
	c.connector.GetSessionMap().Range(func(key, value interface{}) bool {
		routeId := msgService.GetMsgService().GetRouteId(args.Route)
		value.(Agent).PushMessage(routeId, args.Payload)
		return true
	})
	return nil
}