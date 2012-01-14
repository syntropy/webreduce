package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type WorkerStore map[string]func(interface{}) error

type Server struct {
	Prefix string
	Workers WorkerStore
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handle errors by returning a HTTP 500 status code to the client
	defer func() {
		if e := recover(); e != nil {
			trace := debug.Stack()
			msg := fmt.Sprintf("Internal Server Error\n\n%s\n%s", trace, e)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}()

	method := req.Method
	path := req.URL.Path

	if method == "POST" {
		if path == srv.Prefix {
			srv.HandlePost(w, req)
			return
		}

		srv.HandleDataPost(w, req)
		return
	}

	http.NotFound(w, req)
}

func (srv *Server) HandlePost(w http.ResponseWriter, req *http.Request) {

}

func (srv *Server) HandleDataPost(w http.ResponseWriter, req *http.Request) {

}
