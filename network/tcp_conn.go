package network

import (
	"bytes"
	"io"
	"net"
	"sync"
)

type ConnSet map[net.Conn]struct{}

type TCPConn struct {
	sync.Mutex
	conn      net.Conn
}

func newTCPConn(conn net.Conn) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn

	return tcpConn
}

func (tcpConn *TCPConn) Close() {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	tcpConn.conn.(*net.TCPConn).SetLinger(0)
	tcpConn.conn.Close()
}

func (tcpConn *TCPConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

func (tcpConn *TCPConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

// Read buffer length data
func (tcpConn *TCPConn) Read(buffer []byte) (int, error) {
	return io.ReadFull(tcpConn.conn, buffer)
}

func (tcpConn *TCPConn) ReadMsg(buf *bytes.Buffer) error {
	return nil
}

func (tcpConn *TCPConn) WriteMessage(buf []byte) error {
	_, err := tcpConn.conn.Write(buf)
	return err
}