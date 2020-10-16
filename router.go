package juice

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
		router *mux.Router
		OnCatchError func(c *Context, errInterface interface{}) error
}
type RouterOption struct {
	OnCatchError func(c *Context, errInterface interface{}) error
}
func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
func NewRouter(opt RouterOption) *Router {
	r := mux.NewRouter()
	return &Router{
		router: r,
		OnCatchError: opt.OnCatchError,
	}
}

