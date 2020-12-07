package rpcClientService


type Option func(*Options)

type Options struct {
	RegistryType 	string
	RegistryAddr 	string
	BasePath 		string
	SelectMode 		string
}

func newOptions(opts ...Option) *Options  {
	opt := &Options{
	}

	for _,o := range opts {
		o(opt)
	}
	return opt
}

// Register service type. option is consul, etcd, etcdv3, zookeeper.
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

// Select mode of service. options is "RandomSelect","RoundRobin","WeightedRoundRobin",
// "WeightedICMP","ConsistentHash","Closest".
func SelectMode(s string) Option {
	return func(options *Options) {
		options.SelectMode = s
	}
}
