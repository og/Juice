package juice

import "github.com/gorilla/mux"

type Group struct {
	serve *Serve
	router *mux.Router
}

func (serve *Serve) Group() Group {
	return Group{
		serve: serve,
		router: serve.router.PathPrefix("").Subrouter(),
	}
}