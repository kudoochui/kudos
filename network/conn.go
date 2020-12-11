package network

import (
	"bytes"
	"net"
)

type Conn interface {
	Read(buffer []byte) (int, error)
	ReadMsg(buf *bytes.Buffer) error
	WriteMessage(buf []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
}
