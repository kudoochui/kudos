package network

import (
	"encoding/binary"
	"errors"
	"github.com/kudoochui/kudos/protocol/protobuf/pkg"
	"io"
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	lenMsgLen	 uint32
	minMsgLen    uint32
	maxMsgLen    uint32
	littleEndian bool
}

func NewMsgParser() *MsgParser {
	p := new(MsgParser)
	p.lenMsgLen = 4
	p.minMsgLen = 1
	p.maxMsgLen = 4096
	p.littleEndian = false

	return p
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32) {
	p.lenMsgLen, p.minMsgLen, p.maxMsgLen = uint32(lenMsgLen), minMsgLen, maxMsgLen
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// goroutine safe
func (p *MsgParser) Read(conn *TCPConn) (uint32, []byte, error) {
	var b [4]byte

	// read len
	if _, err := io.ReadFull(conn, b[:]); err != nil {
		return 0, nil, err
	}

	// parse len
	var msgLen uint32
	if p.littleEndian {
		msgLen = binary.LittleEndian.Uint32(b[:])
	} else {
		msgLen = binary.BigEndian.Uint32(b[:])
	}

	// check len
	if msgLen > p.maxMsgLen {
		return 0, nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return 0, nil, errors.New("message too short")
	}

	// data
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return 0, nil, err
	}

	var pkgType uint32
	if p.littleEndian {
		pkgType = binary.LittleEndian.Uint32(msgData[:pkg.PKG_TYPE_BYTES])
	} else {
		pkgType = binary.BigEndian.Uint32(msgData[:pkg.PKG_TYPE_BYTES])
	}
	body := msgData[pkg.PKG_TYPE_BYTES:]

	return pkgType, body, nil
}

// goroutine safe
func (p *MsgParser) Write(conn *TCPConn, respId uint32, data []byte) error {
	// get len
	var msgLen = uint32(len(data) + 4)

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	msg := make([]byte, uint32(p.lenMsgLen)+msgLen)

	// write len
	if p.littleEndian {
		binary.LittleEndian.PutUint32(msg, msgLen)
		binary.LittleEndian.PutUint32(msg[p.lenMsgLen:], respId)
	} else {
		binary.BigEndian.PutUint32(msg, msgLen)
		binary.BigEndian.PutUint32(msg[p.lenMsgLen:], respId)
	}

	// write data
	l := p.lenMsgLen + 4
	copy(msg[l:], data)

	conn.WriteMessage(msg)
	return nil
}
