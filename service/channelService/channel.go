package channelService

import (
	"errors"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
	"github.com/kudoochui/kudos/service/rpcClientService"
	"github.com/kudoochui/kudos/utils/array"
	"sync"
)

type Channel struct {
	name 		string
	group 		map[int64]*rpc.Session			//uid => session
	nodeMap 	map[string][]int64				//address => [sessionId]
	lock 		sync.RWMutex
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:  name,
		group: map[int64]*rpc.Session{},
		nodeMap: map[string][]int64{},
	}
}

// Add user to channel.
func (c *Channel) Add(s *rpc.Session) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _,ok := c.group[s.GetUserId()]; ok {
		return errors.New("already in channel")
	}
	c.group[s.GetUserId()] = s.Clone()

	a := c.nodeMap[s.NodeId]
	if a != nil {
		c.nodeMap[s.NodeId] = append(a, s.GetSessionId())
	} else {
		a = make([]int64,0)
		c.nodeMap[s.NodeId] = append(a, s.GetSessionId())
	}
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
	if a, ok := c.nodeMap[s.NodeId]; ok {
		c.nodeMap[s.NodeId] = array.PullInt64(a, s.GetSessionId())
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

func (c *Channel) GetSessions() map[int64]*rpc.Session {
	c.lock.RLock()
	defer c.lock.RUnlock()

	m := make(map[int64]*rpc.Session, 0)
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

	c.lock.RLock()
	nodeMap := make(map[string][]int64, 0)
	for k,v := range c.nodeMap {
		w := make([]int64, len(v))
		copy(w, v)
		nodeMap[k] = w
	}

	if len(excludeUid) > 0 {
		for _,uid := range excludeUid {
			if s,ok := c.group[uid]; ok {
				if a, ok := nodeMap[s.NodeId]; ok {
					nodeMap[s.NodeId] = array.PullInt64(a, s.GetSessionId())
				}
			}
		}
	}
	c.lock.RUnlock()

	for nodeId, sids := range nodeMap {
		args := &rpc.ArgsGroup{
			Sids:    sids,
			Route: 	 route,
			Payload:  data,
		}
		reply := &rpc.ReplyGroup{}
		rpcClientService.GetRpcClientService().Go(nodeId+"@ChannelRemote","PushMessageByGroup", args, reply, nil)
	}
}
