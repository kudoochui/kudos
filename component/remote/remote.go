package remote

import (
	"context"
	"github.com/kudoochui/kudos/component"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpcx/server"
	"github.com/kudoochui/kudos/rpcx/serverplugin"
	metrics "github.com/rcrowley/go-metrics"
	"time"
)

type Remote struct {
	opts    		*Options

	server			*server.Server
}

type RegisterPlugin interface {
	Start() error
}

func NewRemote(opts ...Option) *Remote {
	options := newOptions(opts...)

	return &Remote{
		opts: options,
	}
}

func (r *Remote) OnInit(s component.ServerImpl) {
	r.server = server.NewServer()
	r.addRegistryPlugin()
}

func (r *Remote) OnDestroy() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	r.server.Shutdown(ctx)
}

func (r *Remote) OnRun(closeSig chan bool) {
	go func() {
		err := r.server.Serve("tcp", r.opts.Addr)
		if err != nil {
			log.Info("rpcx serve %v", err)
		}
	}()
}

func (r *Remote) GetRemoteAddrs() string {
	return r.opts.Addr
}

func (r *Remote) RegisterName(nodeId string, rcvr interface{}, metadata string) error {
	//sname := reflect.TypeOf(rcvr).Elem().Name()
	//name := fmt.Sprintf("%s@%s", nodeId, sname)
	return r.server.RegisterName(nodeId, rcvr, metadata)
}

func (r *Remote) addRegistryPlugin() {

	var p RegisterPlugin
	switch r.opts.RegistryType {
	case "consul":
		p = &serverplugin.ConsulRegisterPlugin{
			ServiceAddress: "tcp@" + r.opts.Addr,
			ConsulServers:  []string{r.opts.RegistryAddr},
			BasePath:       r.opts.BasePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
	case "etcd":
		p = &serverplugin.EtcdRegisterPlugin{
			ServiceAddress: "tcp@" + r.opts.Addr,
			EtcdServers:    []string{r.opts.RegistryAddr},
			BasePath:       r.opts.BasePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
	case "etcdv3":
		p = &serverplugin.EtcdV3RegisterPlugin{
			ServiceAddress: "tcp@" + r.opts.Addr,
			EtcdServers:    []string{r.opts.RegistryAddr},
			BasePath:       r.opts.BasePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
	case "zookeeper":
		p = &serverplugin.ZooKeeperRegisterPlugin{
			ServiceAddress:   "tcp@" + r.opts.Addr,
			ZooKeeperServers: []string{r.opts.RegistryAddr},
			BasePath:         r.opts.BasePath,
			Metrics:          metrics.NewRegistry(),
			UpdateInterval:   time.Minute,
		}
	}

	err := p.Start()
	if err != nil {
		log.Error("%v", err)
	}
	r.server.Plugins.Add(p)
}