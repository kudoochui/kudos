package app

import (
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"runtime"
	"sync"
)

type Server interface {
	OnStart()
	Run(closeSig chan bool)
	OnStop()
}


type server struct {
	s       Server
	closeSig chan bool
	wg       sync.WaitGroup
}

var servers []*server

func Register(s Server) {
	m := new(server)
	m.s = s
	m.closeSig = make(chan bool, 1)

	servers = append(servers, m)
}

func Init() {
	for i := 0; i < len(servers); i++ {
		servers[i].s.OnStart()
	}

	for i := 0; i < len(servers); i++ {
		s := servers[i]
		s.wg.Add(1)
		go run(s)
	}
}

func Destroy() {
	for i := len(servers) - 1; i >= 0; i-- {
		m := servers[i]
		//m.closeSig <- true
		close(m.closeSig)
		m.wg.Wait()
		destroy(m)
	}

	// Clean global objects
	rpc.Cleanup()
}

func run(m *server) {
	m.s.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *server) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			l := runtime.Stack(buf, false)
			log.Error("%v: %s", r, buf[:l])
		}
	}()

	m.s.OnStop()
}