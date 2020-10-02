package juice

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Serve struct {
		addr string
		router *mux.Router
		session SessionStore
		OnCatchError func(c *Context, errInterface interface{}) error
}
func (serve Serve) Listen(addr string) error {
		serve.addr = addr
		log.Print("Listen http://127.0.0.1" +addr)
		return http.ListenAndServe(addr, serve.router)
}
type ServeOption struct {
	Session SessionStore
	OnCatchError func(c *Context, errInterface interface{}) error
}
func NewServe(opt ServeOption) *Serve {
	return &Serve{
		router: mux.NewRouter(),
		session: opt.Session,
		OnCatchError: opt.OnCatchError,
	}
}
