package timers

type Option func(*Options)

type Options struct {
	TimerDispatcherLen int
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		TimerDispatcherLen: 20,
	}

	for _,o := range opts {
		o(opt)
	}
	return opt
}

// Address of rpc service
func TimerDispatcherLen(length int) Option {
	return func(options *Options) {
		options.TimerDispatcherLen = length
	}
}
