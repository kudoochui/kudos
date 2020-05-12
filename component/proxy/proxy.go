package proxy

import (
	"context"
	"github.com/kudoochui/rpcx/client"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"sync"
)

type Proxy struct {
	opts 			*Options

	pool          *client.OneClientPool
	lock          sync.RWMutex
	chanCall      chan *rpc.Call
	chanRet       chan *client.Call
	responder     rpc.RpcResponder
	handlerFilter filter
	rpcFilter     filter
}

func NewProxy(opts ...Option) *Proxy {
	options := newOptions(opts...)

	return &Proxy{
		opts:      options,
	}
}

func (r *Proxy) OnInit() {
	r.chanCall = make(chan *rpc.Call, r.opts.ChanCallSize)
	r.chanRet = make(chan *client.Call, r.opts.ChanRetSize)
}

func (r *Proxy) OnDestroy() {

}

func (r *Proxy) Run(closeSig chan bool) {
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
	r.pool = client.NewOneClientPool(r.opts.RpcPoolSize, client.Failtry, s, d, client.DefaultOption)
	r.lock.Unlock()

	for {
		select {
		case <-closeSig:
			goto onEnd
		case ci := <-r.chanCall:
			r.exec(ci)
		case ri := <-r.chanRet:
			if ri.Error != nil {
				log.Error("failed to call: %v", ri.Error)
			} else {
				args := ri.Args.(*rpc.Args)
				if r.handlerFilter != nil {
					r.handlerFilter.After(ri.ServicePath + "." + ri.ServiceMethod, args.MsgId, &args.Session, ri.Reply)
				}
				r.responder.Cb(&args.Session, args.MsgId, ri.Reply)
			}
		}
	}
onEnd:
	log.Info("proxy component closed")
	r.pool.Close()
}

func (r *Proxy) SetRpcResponder(resp rpc.RpcResponder){
	r.responder = resp
}

func (r *Proxy) RpcCall(servicePath string, serviceMethod string, args *rpc.Args, reply interface{}) error {
	call := &rpc.Call{
		Session:     &args.Session,
		MsgId:       args.MsgId,
		ServicePath: servicePath,
		ServiceName: serviceMethod,
		MsgReq:      args.MsgReq,
		MsgResp:     reply,
		Done:        make(chan *client.Call, 1),
	}
	if r.rpcFilter != nil {
		r.rpcFilter.Before(servicePath + "." + serviceMethod, args.MsgId, &args.Session, args.MsgReq)
	}
	r.Go(call)

	done := <- call.Done.(chan *client.Call)
	if r.rpcFilter != nil {
		r.rpcFilter.After(servicePath + "." + serviceMethod, args.MsgId, &args.Session, done.Reply)
	}
	reply = done.Reply
	return nil
}

func (r *Proxy) Go(call *rpc.Call) {
	if call.Done == nil {
		if r.handlerFilter != nil {
			r.handlerFilter.Before(call.ServicePath + "." + call.ServiceName, call.MsgId, call.Session, call.MsgReq)
		}
	}
	r.chanCall <- call
	//select {
	//case r.ChanCall <- call:
	//	// ok
	//default:
	//	log.Debug("rpc: discarding Call due to insufficient Call chan capacity")
	//}
}

func (r *Proxy) exec(call *rpc.Call) {
	args := &rpc.Args{
		Session: *call.Session,
		MsgId: call.MsgId,
		MsgReq:  call.MsgReq,
	}

	xclient := r.pool.Get()
	c := r.chanRet
	if call.Done != nil {
		c = call.Done.(chan *client.Call)
	}
	if _,err := xclient.Go(context.TODO(), call.ServicePath, call.ServiceName, args, call.MsgResp, c); err != nil {
		log.Error("network call error: %v", err)
	}
}

// Set a filter for client handler
func (r *Proxy) SetHandlerFilter(f filter) {
	r.handlerFilter = f
}

// Set a filter for rpc
func (r *Proxy) SetRpcFilter(f filter) {
	r.rpcFilter = f
}