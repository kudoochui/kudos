package protobuf

import (
	"time"
)

type Option func(*Options)

type Options struct {
	MaxConnNum      int
	MaxMsgLen       uint32
	HeartbeatTimeout 	time.Duration

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
}

func newOptions(opts...Option) *Options {
	opt := &Options{
		MaxConnNum:      20000,
		MaxMsgLen:       4096,
		HTTPTimeout:     10 * time.Second,
		HeartbeatTimeout: 20 * time.Second,
	}

	for _,o := range opts {
		o(opt)
	}

	return opt
}

// Max connections support. Default is 20000
func MaxConnNum(num int) Option {
	return func(options *Options) {
		options.MaxConnNum = num
	}
}

// Max message length. If a message exceeds the limit, the connection sends a close message to the peer. Default is 4096
func MaxMsgLen(length uint32) Option {
	return func(options *Options) {
		options.MaxMsgLen = length
	}
}

// Address of tcp server as "host:port"
func TCPAddr(addr string) Option {
	return func(options *Options) {
		options.TCPAddr = addr
	}
}

// Address of websocket server as "host:port"
func WSAddr(addr string) Option {
	return func(options *Options) {
		options.WSAddr = addr
	}
}

// Timeout for http handshake. Default is 10s
func HTTPTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.HTTPTimeout = t
	}
}

// Cert file for https
func CertFile(f string) Option {
	return func(options *Options) {
		options.CertFile = f
	}
}

// Key file for https
func KeyFile(f string) Option {
	return func(options *Options) {
		options.KeyFile = f
	}
}

// Heartbeat timeout. Default is 10s. Disconnect if after 2*t
func HeartbeatTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.HeartbeatTimeout = t
	}
}