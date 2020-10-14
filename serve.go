package juice

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"time"
)

type Serve struct {
		HttpServer *http.Server
		router *mux.Router
		OnCatchError func(c *Context, errInterface interface{}) error
}
// 供 juicetest 使用，其他场景不要使用此方法
func (serve *Serve) HttpTestRouter()  *mux.Router {
	return serve.router
}
func (serve *Serve) Listen(addr string) error {
	log.Print("Listen http://127.0.0.1" +addr)
	serve.HttpServer.Addr = addr
	return serve.HttpServer.ListenAndServe()
}
func (serve *Serve) ListenTLS(addr string, certFile string, keyFile string) error {
	log.Print("Listen http://127.0.0.1" +addr)
	serve.HttpServer.Addr = addr
	return serve.HttpServer.ListenAndServeTLS(certFile, keyFile)
}
type ServeOptionHttp struct {
	TLSConfig *tls.Config
	ReadTimeout time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
	MaxHeaderBytes int
	TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState func(net.Conn, http.ConnState)
	ErrorLog *log.Logger
	BaseContext func(net.Listener) context.Context
	ConnContext func(ctx context.Context, c net.Conn) context.Context
}
type ServeOption struct {
	Http ServeOptionHttp
	OnCatchError func(c *Context, errInterface interface{}) error
}
func NewServe(opt ServeOption) *Serve {
	r := mux.NewRouter()
	httpServer := &http.Server{
		Handler: r,
		TLSConfig: opt.Http.TLSConfig,
		ReadTimeout: opt.Http.ReadTimeout,
		ReadHeaderTimeout: opt.Http.ReadHeaderTimeout,
		WriteTimeout: opt.Http.WriteTimeout,
		IdleTimeout: opt.Http.IdleTimeout,
		MaxHeaderBytes: opt.Http.MaxHeaderBytes,
		TLSNextProto: opt.Http.TLSNextProto,
		ConnState: opt.Http.ConnState,
		ErrorLog: opt.Http.ErrorLog,
		BaseContext: opt.Http.BaseContext,
		ConnContext: opt.Http.ConnContext,
	}
	return &Serve{
		HttpServer: httpServer,
		router: r,
		OnCatchError: opt.OnCatchError,
	}
}

