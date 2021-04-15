package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"github.com/kudoochui/kudos/rpcx/log"
	"github.com/kudoochui/kudos/rpcx/mailbox"
	"github.com/kudoochui/kudos/rpcx/protocol"
	"github.com/kudoochui/kudos/rpcx/share"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

type agent struct {
	server  *Server
	conn 	net.Conn
	sessionMap map[int64]mailbox.Mailbox
}

func newAgent(s *Server, conn net.Conn) *agent {
	return &agent{
		server: s,
		conn: conn,
		sessionMap:make(map[int64]mailbox.Mailbox),
	}
}

func (a *agent)serveConn() {
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
		var mb mailbox.Mailbox
		var ok bool
		if mb, ok = a.sessionMap[sid]; !ok {
			mb = mailbox.Unbounded()
			a.sessionMap[sid] = mb
			mb.RegisterHandlers(a, mailbox.NewDefaultDispatcher(100))
		}
		mb.PostUserMessage(newMessageEnvelope(ctx, req))
	}
}

func (a *agent) InvokeSystemMessage(message interface{}) {

}

func (a *agent) InvokeUserMessage(message interface{}) {
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
	res, err := s.handleRequest(ctx, req)
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
