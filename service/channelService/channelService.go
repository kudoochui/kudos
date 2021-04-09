package channelService

import (
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