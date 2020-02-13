package remote


type Option func(*Options)

type Options struct {
	Addr 			string
	RegistryType 	string
	RegistryAddr 	string
	BasePath 		string
}

func newOptions(opts ...Option) *Options {
	opt := &Options{

	}

	for _,o := range opts {
		o(opt)
	}
	return opt
}

// Address of rpc service
func Addr(s string) Option {
	return func(options *Options) {
		options.Addr = s
	}
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

// Base path of rpc service
func BasePath(s string) Option {
	return func(options *Options) {
		options.BasePath = s
	}
}