package main

import (
	"flag"
	"net/http"
	"webreduce/server"
)

var (
	httpAddr   string
	mongoAddr  string
)

func init() {
	flag.StringVar(&httpAddr, "listen", "localhost:7000", "http address")
	flag.StringVar(&mongoAddr, "dbaddr", "localhost", "database address")
}

func main() {
	prefix := "/apps"
	srv := &server.Server{prefix, make(server.WorkerStore)}
	http.Handle(prefix, srv)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
