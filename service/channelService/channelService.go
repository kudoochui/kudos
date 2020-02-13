package channelService

import (
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
	"sync"
)

var _channelService *ChannelService
var once sync.Once

type ChannelService struct {
	channels sync.Map
}

func GetChannelService() *ChannelService {
	once.Do(func() {
		_channelService = &ChannelService{

		}
	})

	return _channelService
}

func (c *ChannelService) CreateChannel(name string) *Channel {
	channel := NewChannel(name)
	c.channels.Store(name, channel)
	return channel
}

func (c *ChannelService) DestroyChannel(name string) {
	c.channels.Delete(name)
}

func (c *ChannelService) GetChannel(name string) *Channel {
	channel, ok := c.channels.Load(name)
	if ok {
		return channel.(*Channel)
	}
	return nil
}

func (c *ChannelService) PushMessageBySid(nodeAddr string, route string, msg interface{}, sids []int64) {
	data, err := codecService.GetCodecService().Marshal(msg)
	if err != nil {
		log.Error("marshal error: %v", err)
	}
	args := &rpc.ArgsGroup{
		Sids:    sids,
		Route:	 route,
		Payload:  data,
	}
	reply := &rpc.ReplyGroup{}
	rpc.RpcInvoke(nodeAddr, "ChannelRemote", "PushMessageByGroup", args, reply)
}

func (c *ChannelService) Broadcast(nodeAddr string, route string, msg interface{}) {
	data, err := codecService.GetCodecService().Marshal(msg)
	if err != nil {
		log.Error("marshal error: %v", err)
	}
	args := &rpc.ArgsGroup{
		Sids:    []int64{},
		Route:	 route,
		Payload:  data,
	}
	reply := &rpc.ReplyGroup{}
	rpc.RpcInvoke(nodeAddr, "ChannelRemote", "Broadcast", args, reply)
}