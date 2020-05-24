package connector

import (
	"net"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/network"
	"github.com/kudoochui/kudos/protocol/message"
	"github.com/kudoochui/kudos/protocol/pkg"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
	"reflect"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}

type Connection interface {
	OnDisconnect(session *rpc.Session)
}

type agent struct {
	conn      network.Conn
	connector *Connector
	session   *rpc.Session
	userData  interface{}
	agentHandler	*agentHandler
}

func NewAgent(conn network.Conn, connector *Connector) *agent{
	a := &agent{
		conn:      conn,
		connector: connector,
		session:   rpc.NewSession(connector.remote.GetRemoteAddrs()),
		userData:  nil,
	}
	a.agentHandler = NewAgentHandler(a)
	return a
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		pkgType, body := pkg.Decode(data)
		a.agentHandler.Handle(pkgType, body)
	}
}

func (a *agent) OnClose() {
	if a.agentHandler.timerHandler != nil {
		a.connector.timers.ClearTimeout(a.agentHandler.timerHandler)
	}
	a.connector.connection.OnDisconnect(a.session)
	a.connector.sessions.DelSession(a)
}

func (a *agent) WriteMsg(msgId int, msg interface{}) {
	_codec := codecService.GetCodecService()
	if _codec != nil {
		data, err := _codec.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		//routeId := msgService.GetMsgService().GetRouteId(route)
		buffer := message.Encode(msgId, message.TYPE_RESPONSE, 0, data)
		err = a.conn.WriteMsg(pkg.Encode(pkg.TYPE_DATA, buffer...)...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) Write(data ...[]byte) {
	err := a.conn.WriteMsg(data...)
	if err != nil {
		log.Error("write data error: %v", err)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) GetSession() *rpc.Session {
	return a.session
}