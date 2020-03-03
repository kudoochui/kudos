package rpc

import (
	"context"
	"github.com/kudoochui/rpcx/client"
	"github.com/kudoochui/kudos/log"
	"sync"
)

// Group message request
type ArgsGroup struct {
	Sids 	[]int64
	Route 	string
	Payload []byte
}

// Group message response
type ReplyGroup struct {

}

var RpcMap sync.Map

func RpcInvoke(addrs string, servicePath, serviceMethod string, args interface{}, reply interface{}) error {
	var xclient *client.OneClient
	a, ok := RpcMap.Load(addrs)
	if !ok || a == nil {
		d := client.NewPeer2PeerDiscovery("tcp@"+addrs, "")
		xclient = client.NewOneClient(client.Failtry, client.RandomSelect, d, client.DefaultOption)
		RpcMap.Store(addrs, xclient)
	} else {
		xclient = a.(*client.OneClient)
	}

	err := xclient.Call(context.TODO(), servicePath, serviceMethod, args, reply)
	if err != nil {
		log.Error("rpcInvoke error: %v", err)
	}
	return err
}

func Cleanup() {
	RpcMap.Range(func(key, value interface{}) bool {
		c := value.(*client.OneClient)
		c.Close()
		return true
	})
}
