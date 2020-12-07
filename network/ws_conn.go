package network

import (
	"bytes"
	"github.com/gorilla/websocket"
	"net"
	"sync"
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
}

func newWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn

	return wsConn
}

func (wsConn *WSConn) Close() {
	wsConn.Lock()
	defer wsConn.Unlock()
	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()
}

func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe
func (wsConn *WSConn) ReadMsg(buf *bytes.Buffer) error {
	//_, b, err := wsConn.conn.ReadMessage()
	messageType, r, err := wsConn.conn.NextReader()
	_ = messageType
	if err != nil {
		return err
	}
	_, err = buf.ReadFrom(r)
	return err
}

func (wsConn *WSConn) WriteMessage(buf []byte) error {
	return wsConn.conn.WriteMessage(websocket.BinaryMessage, buf)
}