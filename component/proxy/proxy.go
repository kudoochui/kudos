package proxy

import (
	"context"
	"github.com/kudoochui/kudos/component"
	"github.com/kudoochui/kudos/filter"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/rpcx/client"
	"sync"
)

// It is deprecated. Use RpcClientService instead
type Proxy struct {
	opts 			*Options

	rpcClient 		*client.OneClient
	lock 			sync.RWMutex
	rpcFilter     	filter.Filter
}

func NewProxy(opts ...Option) *Proxy {
	options := newOptions(opts...)

	return &Proxy{
		opts:      options,
	}
}

func (r *Proxy) OnInit(s component.ServerImpl) {

}

func (r *Proxy) OnDestroy() {

}

func (r *Proxy) OnRun(closeSig chan bool) {
	var d client.ServiceDiscovery
	switch r.opts.RegistryType {
	case "consul":
		d,_ = client.NewConsulDiscovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	case "etcd":
		d,_ = client.NewEtcdDiscovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	case "etcdv3":
		d,_ = client.NewEtcdV3Discovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	case "zookeeper":
		d,_ = client.NewZookeeperDiscovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	}

	var s client.SelectMode
	switch r.opts.SelectMode {
	case "RoundRobin":
		s = client.RoundRobin
	case "WeightedRoundRobin":
		s = client.WeightedRoundRobin
	case "WeightedICMP":
		s = client.WeightedICMP
	case "ConsistentHash":
		s = client.ConsistentHash
	case "Closest":
		s = client.Closest
	default:
		s = client.RandomSelect
	}

	r.lock.Lock()
	r.rpcClient = client.NewOneClient(client.Failtry, s, d, client.DefaultOption)
	r.lock.Unlock()
}

func (r *Proxy) Call(nodeName string, servicePath string, serviceMethod string, args *rpc.Args, reply interface{}) error {
	if r.rpcFilter != nil {
		r.rpcFilter.Before(servicePath + "." + serviceMethod, args)
	}
	r.lock.RLock()
	err := r.rpcClient.Call(context.TODO(), nodeName, servicePath, serviceMethod, args, reply)
	r.lock.RUnlock()
	if r.rpcFilter != nil {
		r.rpcFilter.After(servicePath + "." + serviceMethod,reply)
	}
	return err
}

func (r *Proxy) Go(nodeName string, servicePath string, serviceMethod string, args *rpc.Args, reply interface{}, chanRet chan *client.Call) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if _,err := r.rpcClient.Go(context.TODO(),nodeName, servicePath, serviceMethod, args, reply, chanRet); err != nil {
		log.Error("rpc call error: %v", err)
	}
}

// Set a filter for rpc
func (r *Proxy) SetRpcFilter(f filter.Filter) {
	r.rpcFilter = f
}