package network

import (
	"bytes"
	"net"
)

type Conn interface {
	ReadMsg(buf *bytes.Buffer) error
	WriteMessage(buf []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
}
