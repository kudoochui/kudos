package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/kudoochui/kudos/component"
	"github.com/kudoochui/kudos/filter"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/rpcx/client"
	"github.com/kudoochui/kudos/rpcx/protocol"
	"github.com/kudoochui/kudos/rpcx/server"
	"github.com/kudoochui/kudos/rpcx/share"
	"reflect"
	"sync"
)

// It is deprecated. Use RpcClientService instead
type Proxy struct {
	opts 			*Options

	rpcClient 		*client.OneClient
	lock 			sync.RWMutex
	rpcFilter     	filter.Filter

	ch				chan *protocol.Message
	serviceMap   	map[string]*server.Service
}

func NewProxy(opts ...Option) *Proxy {
	options := newOptions(opts...)

	return &Proxy{
		opts:      options,
		ch: make(chan *protocol.Message, 100),
		serviceMap: make(map[string]*server.Service),
	}
}

func (r *Proxy) OnInit(s component.ServerImpl) {
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

	var sm client.SelectMode
	switch r.opts.SelectMode {
	case "RoundRobin":
		sm = client.RoundRobin
	case "WeightedRoundRobin":
		sm = client.WeightedRoundRobin
	case "WeightedICMP":
		sm = client.WeightedICMP
	case "ConsistentHash":
		sm = client.ConsistentHash
	case "Closest":
		sm = client.Closest
	default:
		sm = client.RandomSelect
	}

	r.lock.Lock()
	r.rpcClient = client.NewBidirectionalOneClient(client.Failtry, sm, d, client.DefaultOption, r.ch)
	r.lock.Unlock()
}

func (r *Proxy) OnDestroy() {

}

func (r *Proxy) OnRun(closeSig chan bool) {
	go func() {
		for {
			select {
			case <-closeSig:
				return
			case m := <-r.ch:
				if err := r.handleRequest(context.Background(), m); err != nil {
					log.Error("proxy: failed to handle request: %v", err)
				}
			}
		}
	}()
}

func (r *Proxy) Call(nodeName string, servicePath string, serviceMethod string, session protocol.ISession, args *rpc.Args, reply interface{}) error {
	if r.rpcFilter != nil {
		r.rpcFilter.Before(servicePath + "." + serviceMethod, args)
	}
	r.lock.RLock()
	err := r.rpcClient.Call(context.TODO(), nodeName, servicePath, serviceMethod, session, args, reply)
	r.lock.RUnlock()
	if r.rpcFilter != nil {
		r.rpcFilter.After(servicePath + "." + serviceMethod,reply)
	}
	return err
}

func (r *Proxy) Go(nodeName string, servicePath string, serviceMethod string, session protocol.ISession, args *rpc.Args, reply interface{}, chanRet chan *client.Call) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if _,err := r.rpcClient.Go(context.TODO(),nodeName, servicePath, serviceMethod, session, args, reply, chanRet); err != nil {
		log.Error("rpc call error: %v", err)
	}
}

// Set a filter for rpc
func (r *Proxy) SetRpcFilter(f filter.Filter) {
	r.rpcFilter = f
}

func (r *Proxy) Register(rcvr interface{}) {
	service := new(server.Service)
	service.Typ = reflect.TypeOf(rcvr)
	service.Rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(service.Rcvr).Type().Name() // Type
	if sname == "" {
		errorStr := "proxy.Register: no service name for type " + service.Typ.String()
		log.Error(errorStr)
	}
	service.Name = sname

	// Install the methods
	service.Method = server.SuitableMethods(service.Typ, true)

	if len(service.Method) == 0 {
		var errorStr string

		// To help the user, see if a pointer receiver would work.
		method := server.SuitableMethods(reflect.PtrTo(service.Typ), false)
		if len(method) != 0 {
			errorStr = "proxy.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			errorStr = "proxy.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Error(errorStr)
	}
	r.serviceMap[service.Name] = service
}

func (r *Proxy) handleRequest(ctx context.Context, req *protocol.Message) (err error) {
	serviceName := req.ServicePath
	methodName := req.ServiceMethod

	session := rpc.NewSessionFromRpc(req.NodeId, req.SessionId, req.UserId)

	service := r.serviceMap[serviceName]
	if service == nil {
		err = errors.New("rpcx: can't find service " + serviceName)
		return err
	}
	mtype := service.Method[methodName]
	if mtype == nil {
		err = errors.New("rpcx: can't find method " + methodName)
		return err
	}

	var argv interface{}
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem()).Interface()
	} else {
		argv = reflect.New(mtype.ArgType).Interface()
	}

	codec := share.Codecs[req.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		return err
	}

	err = codec.Decode(req.Payload, argv)
	if err != nil {
		return err
	}

	var replyv interface{}
	if mtype.ReplyType.Kind() == reflect.Ptr {
		replyv = reflect.New(mtype.ReplyType.Elem()).Interface()
	} else {
		replyv = reflect.New(mtype.ReplyType).Interface()
	}

	if mtype.ArgType.Kind() != reflect.Ptr {
		err = service.Call(ctx, mtype, reflect.ValueOf(session), reflect.ValueOf(argv).Elem(), reflect.ValueOf(replyv))
	} else {
		err = service.Call(ctx, mtype, reflect.ValueOf(session), reflect.ValueOf(argv), reflect.ValueOf(replyv))
	}

	return err
}