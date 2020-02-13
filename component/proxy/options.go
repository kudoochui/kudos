package proxy


type Option func(*Options)

type Options struct {
	RegistryType 	string
	RegistryAddr 	string
	BasePath 		string
	ChanCallSize	int
	ChanRetSize 	int
	RpcPoolSize 	int
}

func newOptions(opts ...Option) *Options  {
	opt := &Options{
		ChanCallSize: 100,
		ChanRetSize: 100,
		RpcPoolSize: 10,
	}

	for _,o := range opts {
		o(opt)
	}
	return opt
}

// Register service type. option is Consul, etcd ...
func RegistryType(s string) Option {
	return func(options *Options) {
		options.RegistryType = s
	}
}

// Address of register service
func RegistryAddr(s string) Option {
	return func(options *Options) {
		options.RegistryAddr = s
	}
}

// Base path of service
func BasePath(s string) Option {
	return func(options *Options) {
		options.BasePath = s
	}
}

// Size of agent call chan. Default is 20
func ChanCallSize(s int) Option {
	return func(options *Options) {
		options.ChanCallSize = s
	}
}

// Size of return chan which is callback of rpc. Default is 10
func ChanRetSize(s int) Option {
	return func(options *Options) {
		options.ChanRetSize = s
	}
}

// Size of the rpcx client pool. Default is 10.
func RpcPoolSize(s int) Option {
	return func(options *Options) {
		options.RpcPoolSize = s
	}
}