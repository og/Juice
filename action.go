package juice

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Method string
func (m Method) String() string {
	return string(m)
}
const GET Method = "GET"
const POST Method = "POST"
type Action func(c *Context) (reject error)

func (serve *Serve) Action(method Method,path string,  action Action) {
	coreAction(serve, serve.router, method, path, action)
}

func coreAction(serve *Serve, router *mux.Router, method Method,path string,  action Action) {
	router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r, serve)
		err := action(c)
		if err != nil {
			c.CheckError(err) ; return
		}
	}).Methods(method.String())
}

func (group *Group) Action(method Method,path string,  action Action) {
	coreAction(group.serve, group.router, method, path, action)
}