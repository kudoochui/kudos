package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"github.com/kudoochui/kudos/rpcx/log"
	"github.com/kudoochui/kudos/rpcx/mailbox"
	"github.com/kudoochui/kudos/rpcx/protocol"
	"github.com/kudoochui/kudos/rpcx/share"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type userDataCache struct {
	userDataMap 		sync.Map
}

var globalUserDataCache = &userDataCache{}

type ConnAgent struct {
	server  *Server
	conn 	net.Conn
	sessionMap sync.Map

	userDataMap	map[int64]interface{}				//data attachment: uid <=> playerData
	dataMutex sync.RWMutex
}

type TimeTickCallback func(*ServerSession)

func newAgent(s *Server, conn net.Conn) *ConnAgent {
	return &ConnAgent{
		server: s,
		conn: conn,
		userDataMap: make(map[int64]interface{}),
	}
}

func (a *ConnAgent) OnClose() {
	// save user data to global
	for uid,c := range a.userDataMap {
		globalUserDataCache.userDataMap.Store(uid, c)
	}
}

func (a *ConnAgent)serveConn() {
	s := a.server
	conn := a.conn

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
		}
		if share.Trace {
			log.Debugf("server closed conn: %v", conn.RemoteAddr().String())
		}

		closeChannel(s, conn)

		s.Plugins.DoPostConnClose(conn)
	}()

	if isShutdown(s) {
		closeChannel(s, conn)
		return
	}

	if tlsConn, ok := conn.(*tls.Conn); ok {
		if d := s.readTimeout; d != 0 {
			conn.SetReadDeadline(time.Now().Add(d))
		}
		if d := s.writeTimeout; d != 0 {
			conn.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			log.Errorf("rpcx: TLS handshake error from %s: %v", conn.RemoteAddr(), err)
			return
		}
	}

	r := bufio.NewReaderSize(conn, ReaderBuffsize)

	for {
		if isShutdown(s) {
			closeChannel(s, conn)
			return
		}

		t0 := time.Now()
		if s.readTimeout != 0 {
			conn.SetReadDeadline(t0.Add(s.readTimeout))
		}

		ctx := share.WithValue(context.Background(), RemoteConnContextKey, conn)

		req, err := s.readRequest(ctx, r)
		if err != nil {
			if err == io.EOF {
				log.Infof("client has closed this connection: %s", conn.RemoteAddr().String())
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				log.Infof("rpcx: connection %s is closed", conn.RemoteAddr().String())
			} else {
				log.Warnf("rpcx: failed to read request: %v", err)
			}
			return
		}

		if s.writeTimeout != 0 {
			conn.SetWriteDeadline(t0.Add(s.writeTimeout))
		}

		if share.Trace {
			log.Debugf("server received an request %+v from conn: %v", req, conn.RemoteAddr().String())
		}

		ctx = share.WithLocalValue(ctx, StartRequestContextKey, time.Now().UnixNano())
		closeConn := false
		if !req.IsHeartbeat() {
			err = s.auth(ctx, req)
			closeConn = err != nil
		}

		if err != nil {
			if !req.IsOneway() {
				res := req.Clone()
				res.SetMessageType(protocol.Response)
				if len(res.Payload) > 1024 && req.CompressType() != protocol.None {
					res.SetCompressType(req.CompressType())
				}
				handleError(res, err)
				s.Plugins.DoPreWriteResponse(ctx, req, res, err)
				data := res.EncodeSlicePointer()
				_, err := conn.Write(*data)
				protocol.PutData(data)
				s.Plugins.DoPostWriteResponse(ctx, req, res, err)
				protocol.FreeMsg(res)
			} else {
				s.Plugins.DoPreWriteResponse(ctx, req, nil, err)
			}
			protocol.FreeMsg(req)
			// auth failed, closed the connection
			if closeConn {
				log.Infof("auth failed for conn %s: %v", conn.RemoteAddr().String(), err)
				return
			}
			continue
		}

		if req.IsHeartbeat() {
			s.Plugins.DoHeartbeatRequest(ctx, req)
			req.SetMessageType(protocol.Response)
			data := req.EncodeSlicePointer()
			conn.Write(*data)
			protocol.PutData(data)
		}

		sid := req.SessionId
		a.PostMessage(ctx, sid, req)
	}
}

func (a *ConnAgent) PostMessage(ctx *share.Context, sid int64, req *protocol.Message) {
	var mb interface{}
	var ok bool
	if mb, ok = a.sessionMap.Load(sid); !ok {
		mb = mailbox.Unbounded()
		a.sessionMap.Store(sid, mb)
		mb.(mailbox.Mailbox).RegisterHandlers(a, mailbox.NewDefaultDispatcher(100))
	}
	mb.(mailbox.Mailbox).PostUserMessage(newMessageEnvelope(ctx, req))
}

// Run in user actor goroutine
func (a *ConnAgent) RemoveSession(sessionId int64) {
	if mb, ok := a.sessionMap.Load(sessionId); ok {
		mb.(mailbox.Mailbox).PostSystemMessage(&mailbox.SuspendMailbox{})
	}
	a.sessionMap.Delete(sessionId)
}

// Local call, no return
func (a *ConnAgent) Go(route string, session protocol.ISession, args interface{}) error {
	rr := strings.Split(route, ".")
	if len(rr) < 3 {
		return errors.New("invalid route")
	}

	ctx := share.WithValue(context.Background(), StartSendRequestContextKey, time.Now().UnixNano())
	a.server.Plugins.DoPreWriteRequest(ctx)

	req := protocol.GetPooledMsg()
	req.SetMessageType(protocol.Request)

	seq := atomic.AddUint64(&a.server.seq, 1)
	req.SetSeq(seq)
	req.SetOneway(true)
	req.SetSerializeType(protocol.MsgPack)
	req.ServicePath = rr[1]
	req.ServiceMethod = rr[2]
	req.NodeId = session.GetNodeId()
	req.SessionId = session.GetSessionId()
	req.UserId = session.GetUserId()

	// TODO: local call, no need to pack
	codec := share.Codecs[protocol.MsgPack]
	data, err := codec.Encode(args)
	if err != nil {
		return err
	}
	req.Payload = data

	a.server.Plugins.DoPostWriteRequest(ctx, req, err)
	a.PostMessage(ctx, session.GetSessionId(), req)
	return nil
}

func (a *ConnAgent) RegisterTimeTick(session *ServerSession, cb TimeTickCallback) {
	if mb, ok := a.sessionMap.Load(session.GetSessionId()); ok {
		mb.(mailbox.Mailbox).PostTimeMessage(newTimeEnvelope(session, cb))
	} else {
		log.Errorf("RegisterTimeTick error %+v", session)
	}
}

// Tick every 100ms
func (a *ConnAgent) OnTimeTick(message interface{}) {
	msg := message.(*TimeEnvelope)
	msg.Cb(msg.Session)
}

func (a *ConnAgent) InvokeSystemMessage(message interface{}) {

}

func (a *ConnAgent) InvokeUserMessage(message interface{}) {
	s := a.server
	conn := a.conn
	msg := message.(*MessageEnvelope)
	ctx := msg.Context
	req := msg.Request

	resMetadata := make(map[string]string)
	ctx = share.WithLocalValue(share.WithLocalValue(ctx, share.ReqMetaDataKey, req.Metadata),
		share.ResMetaDataKey, resMetadata)

	cancelFunc := parseServerTimeout(ctx, req)
	if cancelFunc != nil {
		defer cancelFunc()
	}

	s.Plugins.DoPreHandleRequest(ctx, req)

	if share.Trace {
		log.Debugf("server handle request %+v from conn: %v", req, conn.RemoteAddr().String())
	}
	res, err := s.handleRequest(a, ctx, req)
	if err != nil {
		log.Warnf("rpcx: failed to handle request: %v", err)
	}

	s.Plugins.DoPreWriteResponse(ctx, req, res, err)
	if !req.IsOneway() {
		if len(resMetadata) > 0 { // copy meta in context to request
			meta := res.Metadata
			if meta == nil {
				res.Metadata = resMetadata
			} else {
				for k, v := range resMetadata {
					if meta[k] == "" {
						meta[k] = v
					}
				}
			}
		}

		if len(res.Payload) > 1024 && req.CompressType() != protocol.None {
			res.SetCompressType(req.CompressType())
		}
		data := res.EncodeSlicePointer()
		conn.Write(*data)
		protocol.PutData(data)
	}
	s.Plugins.DoPostWriteResponse(ctx, req, res, err)

	if share.Trace {
		log.Debugf("server write response %+v for an request %+v from conn: %v", res, req, conn.RemoteAddr().String())
	}

	protocol.FreeMsg(req)
	protocol.FreeMsg(res)
}

func (a *ConnAgent) GetData(userId int64) interface{} {
	a.dataMutex.RLock()
	if c, ok := a.userDataMap[userId]; ok {
		a.dataMutex.RUnlock()
		return c
	}
	a.dataMutex.RUnlock()

	// find global
	if c, ok := globalUserDataCache.userDataMap.Load(userId); ok {
		a.SetData(userId, c)
		globalUserDataCache.userDataMap.Delete(userId)
		return c
	}

	return nil
}

func (a *ConnAgent) SetData(userId int64, data interface{})  {
	a.dataMutex.Lock()
	defer a.dataMutex.Unlock()

	a.userDataMap[userId] = data
}