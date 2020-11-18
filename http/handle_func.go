package jhttp

import (
	"github.com/gorilla/mux"
	"net/http"
)


type HandleFunc func(c *Context) (reject error)

func (serve *Router) HandleFunc(method Method,path string,  action HandleFunc) {
	coreHandleFunc(serve, serve.router, method, path, action)
}

func coreHandleFunc(serve *Router, router *mux.Router, method Method,path string,  action HandleFunc) {
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r, serve)
		defer func() {
			r := recover()
			if r  != nil {
				c.CheckError(r) ; return
			}
		}()
		err := action(c)
		if err != nil {
			c.CheckError(err) ; return
		}
	}).Methods(method.String())
}

func (group *Group) HandleFunc(method Method,path string,  action HandleFunc) {
	coreHandleFunc(group.serve, group.router, method, path, action)
}