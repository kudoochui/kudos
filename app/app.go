package app

import (
	"github.com/kudoochui/kudos/log"
	"os"
	"os/signal"
)

func Run(servers ...Server) {

	for i := 0; i < len(servers); i++ {
		Register(servers[i])
	}
	Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Warning("Server closing down (signal: %v)", sig)
	Destroy()
}