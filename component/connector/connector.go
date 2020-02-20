package connector

import (
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/network"
	"github.com/kudoochui/kudos/rpc"
)

type Connector struct {
	opts 			*Options
	sessions		*sessionMap
	sessionRemote	*SessionRemote
	channelRemote 	*ChannelRemote
	route           rpc.RpcRouter
	customerRoute 	rpc.CustomerRoute
	remote			rpc.HandlerRegister
	connection 		Connection
	timers 			*Timers
}

func NewConnector(opts ...Option) *Connector{
	options := newOptions(opts...)
	c := &Connector{
		opts:			options,
		sessions:		 &sessionMap{},
	}
	c.sessionRemote = NewSessionRemote(c)
	c.channelRemote = NewChannelRemote(c)
	c.timers = NewTimer()
	return c
}

func (c *Connector) OnInit() {}

func (c *Connector) Run(closeSig chan bool) {
	go c.timers.RunTimer(closeSig)

	var wsServer *network.WSServer
	if c.opts.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = c.opts.WSAddr
		wsServer.MaxConnNum = c.opts.MaxConnNum
		wsServer.PendingWriteNum = c.opts.PendingWriteNum
		wsServer.MaxMsgLen = c.opts.MaxMsgLen
		wsServer.HTTPTimeout = c.opts.HTTPTimeout
		wsServer.CertFile = c.opts.CertFile
		wsServer.KeyFile = c.opts.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := NewAgent(conn, c)
			//if c.AgentChanRPC != nil {
			//	c.AgentChanRPC.Go("NewAgent", a)
			//}
			c.sessions.AddSession(a)
			return a
		}
	}

	var tcpServer *network.TCPServer
	if c.opts.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = c.opts.TCPAddr
		tcpServer.MaxConnNum = c.opts.MaxConnNum
		tcpServer.PendingWriteNum = c.opts.PendingWriteNum
		tcpServer.LenMsgLen = c.opts.LenMsgLen
		tcpServer.MaxMsgLen = c.opts.MaxMsgLen
		tcpServer.LittleEndian = c.opts.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := NewAgent(conn, c)
			//if c.AgentChanRPC != nil {
			//	c.AgentChanRPC.Go("NewAgent", a)
			//}
			c.sessions.AddSession(a)
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
		log.Info("websocket server start at: %s", c.opts.WSAddr)
	}
	if tcpServer != nil {
		tcpServer.Start()
		log.Info("tcp server start at: %s", c.opts.TCPAddr)
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
		log.Info("websocket server %s closed", c.opts.WSAddr)
	}
	if tcpServer != nil {
		tcpServer.Close()
		log.Info("tcp server %s closed", c.opts.TCPAddr)
	}
}

func (c *Connector) OnDestroy() {}

func (c *Connector) SetRouter(route rpc.RpcRouter){
	c.route = route
}

func (c *Connector) Route(f rpc.CustomerRoute){
	c.customerRoute = f
}

func (c *Connector) Cb(session *rpc.Session, msgId int, msg interface{}) {
	agent, err := c.sessions.GetAgent(session.GetSessionId())
	if err != nil {
		log.Error("%v", err)
		return
	}

	agent.WriteMsg(msgId, msg)
}

// Register session service
func (c *Connector) SetRegisterServiceHandler(r rpc.HandlerRegister){
	c.remote = r
	c.remote.RegisterHandler(c.sessionRemote,"")
	c.remote.RegisterHandler(c.channelRemote,"")
}

func (c *Connector) SetConnectionListener(conn Connection) {
	c.connection = conn
}