package web

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kudoochui/kudos/log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

type WebServer struct {
	opts		*Options
	server 		*http.Server
	router 		*mux.Router
	routeMap 	map[string]Handler
}

func NewWebServer(opts ...Option) *WebServer {
	options := newOptions(opts...)

	web := &WebServer{
		opts: options,
		routeMap: map[string]Handler{},
	}

	return web
}

func (w *WebServer) Route(path string, f Handler) {
	w.routeMap[path] = f
}

func (w *WebServer) OnInit() {

}

func (w *WebServer) OnDestroy() {

}

func (w *WebServer) Run(closeSig chan bool) {
	w.router = mux.NewRouter()
	for p,f := range w.routeMap {
		w.router.HandleFunc(p, f)
	}
	w.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", w.opts.ListenIp, w.opts.ListenPort),
		WriteTimeout: w.opts.WriteTimeout,
		ReadTimeout:  w.opts.ReadTimeout,
		IdleTimeout:  w.opts.IdleTimeout,
		Handler: w.router,
	}

	log.Info("web server listen at: %s:%d", w.opts.ListenIp, w.opts.ListenPort)

	go func() {
		if err := w.server.ListenAndServe(); err != nil {
			log.Info("web server: %s", err)
		}
	}()

	<-closeSig

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), w.opts.CloseTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	w.server.Shutdown(ctx)
}