package connector

import (
	"context"
	"github.com/kudoochui/kudos/protocol/message"
	"github.com/kudoochui/kudos/protocol/pkg"
	"github.com/kudoochui/kudos/rpc"
	msgService "github.com/kudoochui/kudos/service/msgService"
)

type ChannelRemote struct {
	connector	*Connector
}

func NewChannelRemote(conn *Connector) *ChannelRemote {
	return &ChannelRemote{connector:conn}
}

// Push message to client by uids.
func (c *ChannelRemote) PushMessageByGroup(ctx context.Context, args *rpc.ArgsGroup, reply *rpc.ReplyGroup) error {
	for _, sid := range args.Sids {
		if a, err := c.connector.sessions.GetAgent(sid); err == nil {
			routeId := msgService.GetMsgService().GetRouteId(args.Route)
			buffer := message.Encode(0, message.TYPE_PUSH, routeId, args.Payload)
			a.Write(pkg.Encode(pkg.TYPE_DATA, buffer)...)
		}
	}
	return nil
}

// Broadcast to all the client connectd with current frontend server.
func (c *ChannelRemote) Broadcast(ctx context.Context, args *rpc.ArgsGroup, reply *rpc.ReplyGroup) error {
	c.connector.sessions.Range(func(key, value interface{}) bool {
		routeId := msgService.GetMsgService().GetRouteId(args.Route)
		buffer := message.Encode(0, message.TYPE_PUSH, routeId, args.Payload)
		value.(*agent).Write(pkg.Encode(pkg.TYPE_DATA, buffer)...)
		return true
	})
	return nil
}