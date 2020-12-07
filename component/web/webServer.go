package web

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kudoochui/kudos/log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

type WebServer struct {
	opts		*Options
	server 		*http.Server
	router 		*mux.Router
	routeMap 	map[string]Handler
	prefixRouteMap map[string]http.Handler
}

func NewWebServer(opts ...Option) *WebServer {
	options := newOptions(opts...)

	web := &WebServer{
		opts: options,
		routeMap: map[string]Handler{},
		prefixRouteMap: map[string]http.Handler{},
	}

	return web
}

func (w *WebServer) Route(path string, f Handler) {
	w.routeMap[path] = f
}

func (w *WebServer) PrefixRoute(path string, f http.Handler) {
	w.prefixRouteMap[path] = f
}

func (w *WebServer) OnInit() {

}

func (w *WebServer) OnDestroy() {
	// Create a deadline to wait for.
	ctx, _ := context.WithTimeout(context.Background(), w.opts.CloseTimeout)
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	w.server.Shutdown(ctx)
}

func (w *WebServer) OnRun(closeSig chan bool) {
	w.router = mux.NewRouter()
	for p,f := range w.routeMap {
		w.router.HandleFunc(p, f)
	}

	for p,f := range w.prefixRouteMap {
		w.router.PathPrefix(p).Handler(f)
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
}