package protobuf

import (
	"encoding/binary"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/network"
	"github.com/kudoochui/kudos/protocol"
	"github.com/kudoochui/kudos/protocol/protobuf/pkg"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
	"github.com/kudoochui/rpcx/client"
	"net"
	"reflect"
)

type agent struct {
	conn      network.Conn
	connector *Connector
	session   *rpc.Session
	userData  interface{}
	agentHandler	*agentHandler
	chanRet		chan *client.Call
	writeChan 	chan *[]byte
}

func NewAgent(conn network.Conn, connector *Connector) *agent{
	a := &agent{
		conn:      conn,
		connector: connector,
		session:   rpc.NewSession(connector.nodeId),
		userData:  nil,
		chanRet: make(chan *client.Call, 100),
		writeChan: make(chan *[]byte, 100),
	}
	a.agentHandler = NewAgentHandler(a)
	return a
}

func (a *agent) Run() {
	go func() {
		defer a.conn.Close()
		for {
			select {
			case ri := <-a.chanRet:
				if ri.Error != nil {
					log.Error("failed to call: %v", ri.Error)
				} else {
					args := ri.Args.(*rpc.Args)
					if a.connector.handlerFilter != nil {
						a.connector.handlerFilter.After(ri.ServicePath + "." + ri.ServiceMethod, ri)
					}

					a.WriteResponse(args.MsgId, ri.Reply)
				}
			case b := <-a.writeChan:
				if b == nil {
					return
				}

				err := a.conn.WriteMessage(*b)
				protocol.FreePoolBuffer(b)
				if err != nil {
					log.Error("ws WriteMessage: %s", err.Error())
					return
				}
			}
		}
	}()

	switch a.conn.(type) {
	case *network.WSConn:
		for {
			buffer := protocol.GetPoolMsg()
			err := a.conn.ReadMsg(buffer)
			if err != nil {
				log.Debug("read message: %v", err)
				break
			}

			_, pkgType, body := pkg.Decode(buffer.Bytes())
			a.agentHandler.Handle(pkgType, body)
			buffer.Reset()
			protocol.FreePoolMsg(buffer)
		}
		break
	case *network.TCPConn:
		for {
			headBuffer := protocol.GetUint32PoolData()

			// read len
			if _, err := a.conn.Read(*headBuffer); err != nil {
				protocol.PutUint32PoolData(headBuffer)
				break
			}

			// parse len
			var msgLen uint32
			if pkg.GetByteOrder() {
				msgLen = binary.LittleEndian.Uint32(*headBuffer)
			} else {
				msgLen = binary.BigEndian.Uint32(*headBuffer)
			}

			// check len
			//if msgLen > p.maxMsgLen {
			//	return nil, errors.New("message too long")
			//} else if msgLen < p.minMsgLen {
			//	return nil, errors.New("message too short")
			//}

			// data
			payloadBuffer := protocol.GetPoolBuffer(int(msgLen))
			protocol.PutUint32PoolData(headBuffer)
			if _, err := a.conn.Read(*payloadBuffer); err != nil {
				protocol.FreePoolBuffer(payloadBuffer)
				break
			}

			//pkgType, body := pkg.Decode(*payloadBuffer)
			var pkgType uint32
			if pkg.GetByteOrder() {
				pkgType = binary.LittleEndian.Uint32((*payloadBuffer)[:pkg.PKG_TYPE_BYTES])
			} else {
				pkgType = binary.BigEndian.Uint32((*payloadBuffer)[:pkg.PKG_TYPE_BYTES])
			}
			body := (*payloadBuffer)[pkg.PKG_TYPE_BYTES:]
			a.agentHandler.Handle(pkgType, body)
			protocol.FreePoolBuffer(payloadBuffer)
		}
		break
	}

	close(a.writeChan)
}

func (a *agent) OnClose() {
	if a.agentHandler.timerHandler != nil {
		a.connector.timers.ClearTimeout(a.agentHandler.timerHandler)
	}
	a.connector.connection.OnDisconnect(a.session)
	a.connector.sessions.DelSession(a)
}

func (a *agent) WriteResponse(msgId int, msg interface{}) {
	_codec := codecService.GetCodecService()
	if _codec != nil {
		data, err := _codec.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		buffer := pkg.Encode(uint32(msgId), data)
		err = a.conn.WriteMessage(*buffer)
		protocol.FreePoolBuffer(buffer)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

// Write to channel. Make sure buffer from protocol.GetPoolBuffer()
func (a *agent) Write(data *[]byte) {
	a.writeChan <- data
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

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) GetSession() *rpc.Session {
	return a.session
}

func (a *agent) PushMessage(routeId uint32, data []byte) {
	buffer := pkg.Encode(routeId, data)
	a.Write(buffer)
}

func (a *agent) KickMessage(reason string) {
	ret := &pkg.RespResult{Code:int32(pkg.EErrorCode_ERROR_KICK_BY_SERVER), Msg:reason}
	buffer,_ := codecService.GetCodecService().Marshal(ret)
	a.Write(pkg.Encode(uint32(pkg.EMsgType_TYPE_KICK_BY_SERVER), buffer))
}