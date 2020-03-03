package remote

import (
	"context"
	"fmt"
	"github.com/kudoochui/kudos/log"
	"reflect"
	"time"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/kudoochui/rpcx/server"
	"github.com/kudoochui/rpcx/serverplugin"
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

func (r *Remote) OnInit() {
	r.server = server.NewServer()
	r.addRegistryPlugin()
}

func (r *Remote) OnDestroy() {

}

func (r *Remote) Run(closeSig chan bool) {
	go func() {
		err := r.server.Serve("tcp", r.opts.Addr)
		if err != nil {
			log.Info("rpcx serve %v", err)
		}
	}()

	<- closeSig

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	r.server.Shutdown(ctx)
}

func (r *Remote) GetRemoteAddrs() string {
	return r.opts.Addr
}

func (r *Remote) RegisterHandler(rcvr interface{}, metadata string) error {
	return r.server.Register(rcvr, metadata)
}

func (r *Remote) RegisterName(nodeId string, rcvr interface{}, metadata string) error {
	sname := reflect.TypeOf(rcvr).Elem().Name()
	name := fmt.Sprintf("%s@%s", nodeId, sname)
	return r.server.RegisterName(name, rcvr, metadata)
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