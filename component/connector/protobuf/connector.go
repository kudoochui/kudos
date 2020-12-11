package protobuf

import (
	"github.com/kudoochui/kudos/component"
	"github.com/kudoochui/kudos/component/connector"
	"github.com/kudoochui/kudos/component/remote"
	"github.com/kudoochui/kudos/filter"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/network"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
)

type Connector struct{
	opts 			*Options
	nodeId			string
	sessions		*connector.SessionMap
	sessionRemote	*connector.SessionRemote
	channelRemote 	*connector.ChannelRemote
	customerRoute 	rpc.CustomerRoute
	remote			*remote.Remote
	//proxy 			*proxy.Proxy
	handlerFilter 	filter.Filter
	connection 		connector.Connection
	timers 			*connector.Timers
	wsServer 		*network.WSServer
	tcpServer 		*network.TCPServer
}

func NewConnector(opts ...Option) *Connector{
	options := newOptions(opts...)
	c := &Connector{
		opts:			options,
		sessions:		 &connector.SessionMap{},
	}
	c.sessionRemote = connector.NewSessionRemote(c)
	c.channelRemote = connector.NewChannelRemote(c)
	c.timers = connector.NewTimer()
	return c
}

func (c *Connector) OnInit(server component.ServerImpl) {
	if codecService.GetCodecType() != codecService.TYPE_CODEC_PROTOBUF {
		panic("Please use protobuf codec!")
	}
	c.nodeId = server.GetServerId()
	c.remote = server.GetComponent("remote").(*remote.Remote)
	//c.proxy = server.GetComponent("proxy").(*proxy.Proxy)
}

func (c *Connector) OnRun(closeSig chan bool) {
	c.remote.RegisterName(c.nodeId, c.sessionRemote,"")
	c.remote.RegisterName(c.nodeId, c.channelRemote,"")

	go c.timers.RunTimer(closeSig)

	if c.opts.WSAddr != "" {
		c.wsServer = new(network.WSServer)
		c.wsServer.Addr = c.opts.WSAddr
		c.wsServer.MaxConnNum = c.opts.MaxConnNum
		c.wsServer.MaxMsgLen = c.opts.MaxMsgLen
		c.wsServer.HTTPTimeout = c.opts.HTTPTimeout
		c.wsServer.CertFile = c.opts.CertFile
		c.wsServer.KeyFile = c.opts.KeyFile
		c.wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := NewAgent(conn, c)
			c.sessions.AddSession(a)
			return a
		}
	}

	if c.opts.TCPAddr != "" {
		c.tcpServer = new(network.TCPServer)
		c.tcpServer.Addr = c.opts.TCPAddr
		c.tcpServer.MaxConnNum = c.opts.MaxConnNum
		//c.tcpServer.PendingWriteNum = c.opts.PendingWriteNum
		//c.tcpServer.LenMsgLen = c.opts.LenMsgLen
		//c.tcpServer.MaxMsgLen = c.opts.MaxMsgLen
		//c.tcpServer.LittleEndian = c.opts.LittleEndian
		c.tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := NewAgent(conn, c)
			c.sessions.AddSession(a)
			return a
		}
	}

	if c.wsServer != nil {
		c.wsServer.Start()
		log.Info("websocket server start at: %s", c.opts.WSAddr)
	}
	if c.tcpServer != nil {
		c.tcpServer.Start()
		log.Info("tcp server start at: %s", c.opts.TCPAddr)
	}
}

func (c *Connector) OnDestroy() {
	if c.wsServer != nil {
		c.wsServer.Close()
		log.Info("websocket server %s closed", c.opts.WSAddr)
	}
	if c.tcpServer != nil {
		c.tcpServer.Close()
		log.Info("tcp server %s closed", c.opts.TCPAddr)
	}
}

func (c *Connector) Route(f rpc.CustomerRoute){
	c.customerRoute = f
}

func (c *Connector) SetConnectionListener(conn connector.Connection) {
	c.connection = conn
}

// Set a filter for client handler
func (c *Connector) SetHandlerFilter(f filter.Filter) {
	c.handlerFilter = f
}

func (c* Connector) GetSessionMap() *connector.SessionMap {
	return c.sessions
}