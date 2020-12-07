package connector

import (
	"context"
	"github.com/kudoochui/kudos/protocol/pomelo/message"
	"github.com/kudoochui/kudos/protocol/pomelo/pkg"
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
			buffer := message.Encode(0, message.TYPE_PUSH, routeId, args.Payload)
			a.Write(pkg.Encode(pkg.TYPE_DATA, buffer))
		}
	}
	return nil
}

// Broadcast to all the client connectd with current frontend server.
func (c *ChannelRemote) Broadcast(ctx context.Context, args *rpc.ArgsGroup, reply *rpc.ReplyGroup) error {
	//log.Debug(">>>> %s broadcast: %s, %s", c.connector.opts.WSAddr, args.Route, string(args.Payload))
	c.connector.GetSessionMap().Range(func(key, value interface{}) bool {
		routeId := msgService.GetMsgService().GetRouteId(args.Route)
		buffer := message.Encode(0, message.TYPE_PUSH, routeId, args.Payload)
		value.(Agent).Write(pkg.Encode(pkg.TYPE_DATA, buffer))
		return true
	})
	return nil
}