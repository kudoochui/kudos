package rpc

import (
	"context"
	"github.com/kudoochui/rpcx/client"
	"github.com/kudoochui/kudos/log"
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

func RpcInvoke(addrs string, servicePath, serviceMethod string, args interface{}, reply interface{}) error {
	d := client.NewPeer2PeerDiscovery("tcp@"+addrs, "")
	xclient := client.NewXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	err := xclient.Call(context.TODO(), serviceMethod, args, reply)
	if err != nil {
		log.Error("rpcInvoke error: %v", err)
	}
	return err
}
