package rpcClientService

import (
	"context"
	"github.com/kudoochui/kudos/filter"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/rpcx/client"
	"sync"
)

var _rpcClientService *RpcClientService
var once sync.Once

type RpcClientService struct {
	opts 			*Options

	rpcClient 		*client.OneClient
	lock 			sync.RWMutex
	rpcFilter     	filter.Filter
}

func GetRpcClientService() *RpcClientService {
	once.Do(func() {
		_rpcClientService = &RpcClientService{

		}
	})

	return _rpcClientService
}


func (r *RpcClientService) Initialize(opts ...Option) {
	options := newOptions(opts...)
	r.opts = options
	var d client.ServiceDiscovery
	switch r.opts.RegistryType {
	case "consul":
		d = client.NewConsulDiscovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	case "etcd":
		d = client.NewEtcdDiscovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	case "etcdv3":
		d = client.NewEtcdV3Discovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
	case "zookeeper":
		d = client.NewZookeeperDiscovery(r.opts.BasePath, "", []string{r.opts.RegistryAddr}, nil)
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

func (r *RpcClientService) Call(servicePath string, serviceMethod string, args interface{}, reply interface{}) error {
	if r.rpcFilter != nil {
		r.rpcFilter.Before(servicePath + "." + serviceMethod, args)
	}
	r.lock.RLock()
	err := r.rpcClient.Call(context.TODO(), servicePath, serviceMethod, args, reply)
	if err != nil {
		log.Error("rpc call error: %v", err)
	}
	r.lock.RUnlock()
	if r.rpcFilter != nil {
		r.rpcFilter.After(servicePath + "." + serviceMethod, reply)
	}
	return err
}

func (r *RpcClientService) Go(servicePath string, serviceMethod string, args interface{}, reply interface{}, chanRet chan *client.Call) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if _,err := r.rpcClient.Go(context.TODO(),servicePath, serviceMethod, args, reply, chanRet); err != nil {
		log.Error("rpc go error: %v", err)
	}
}

// Set a filter for rpc
func (r *RpcClientService) SetRpcFilter(f filter.Filter) {
	r.rpcFilter = f
}