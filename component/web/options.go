package web

import "time"

type Option func(*Options)

type Options struct {
	ListenIp 		string
	ListenPort 		int
	WriteTimeout	time.Duration
	ReadTimeout		time.Duration
	IdleTimeout 	time.Duration
	CloseTimeout	time.Duration
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		ListenIp:		"0.0.0.0",
		ListenPort: 	5050,
		WriteTimeout: 	15 * time.Second,
		ReadTimeout: 	15 * time.Second,
		IdleTimeout: 	60 * time.Second,
		CloseTimeout: 	5 * time.Second,
	}

	for _,o := range opts {
		o(opt)
	}
	return opt
}

// Ip of listening. Default is "0.0.0.0"
func ListenIp(ip string) Option {
	return func(options *Options) {
		options.ListenIp = ip
	}
}

// Port of listening. Default is 5050
func ListenPort(port int) Option {
	return func(options *Options) {
		options.ListenPort = port
	}
}

// Maximum duration before timing out writes of the response. Default is 15s
func WriteTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.WriteTimeout = t
	}
}

// Maximum duration for reading the entire request. Default is 15s
func ReadTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.ReadTimeout = t
	}
}

// Maximum amount of time to wait for the next request. Default is 60s
func IdleTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.IdleTimeout = t
	}
}

// Deadline to wait for closing. Default is 5s
func CloseTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = t
	}
}