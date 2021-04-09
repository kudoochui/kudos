package channelService

import (
	"errors"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/rpcx/server"
	"github.com/kudoochui/kudos/service/codecService"
	"sync"
)

type Channel struct {
	name 		string
	group 		map[int64]*server.ServerSession			//uid => session
	lock 		sync.RWMutex
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:  name,
		group: map[int64]*server.ServerSession{},
	}
}

// Add user to channel.
func (c *Channel) Add(s *server.ServerSession) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _,ok := c.group[s.GetUserId()]; ok {
		return errors.New("already in channel")
	}
	c.group[s.GetUserId()] = s
	return nil
}

// Remove user from channel.
func (c *Channel) Leave(uid int64)  {
	c.lock.Lock()
	defer c.lock.Unlock()

	s := c.group[uid]
	if s == nil {
		return
	}

	delete(c.group, uid)
}

// Get userId array
func (c *Channel) GetMembers() []int64  {
	c.lock.RLock()
	defer c.lock.RUnlock()

	array := make([]int64, 0)
	for k,_ := range c.group {
		array = append(array, k)
	}
	return array
}

func (c *Channel) GetSessions() map[int64]*server.ServerSession {
	c.lock.RLock()
	defer c.lock.RUnlock()

	m := make(map[int64]*server.ServerSession, 0)
	for k,v := range c.group {
		m[k] = v
	}
	return m
}

// Push message to all the members in the channel, exclude uid in the excludeUid.
func (c *Channel) PushMessage(route string, msg interface{}, excludeUid []int64) {
	data, err := codecService.GetCodecService().Marshal(msg)
	if err != nil {
		log.Error("marshal error: %v", err)
	}

	excludeMap := make(map[int64]bool, len(excludeUid))
	for _, uid := range excludeUid {
		excludeMap[uid] = true
	}
	args := &rpc.ArgsGroup{
		Route: 	 route,
		Payload:  data,
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	for uid, session := range c.group {
		if !excludeMap[uid] {
			session.PushMessage("ChannelRemote", "PushMessage", args)
		}
	}
}
