package juice

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Next func() error
type Middleware func(c *Context, next Next) (reject error)
func (serve *Serve) Use(middleware Middleware) {
	middlewareUse(serve, serve.router, middleware)
}
func middlewareUse(serve *Serve, router *mux.Router, middleware Middleware) {
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print(r.Method, " ", r.URL.RawPath)
			c := NewContext(w,r, serve)
			mwErr := middleware(c, func() error {
				handler.ServeHTTP(w, r)
				return nil
			})
			if mwErr != nil {
				c.CheckError(mwErr) ; return
			}
		})
	})
}
func (group *Group) Use(middleware Middleware) {
	middlewareUse(group.serve, group.router, middleware)
}