package main

import (
	"flag"
	"net/http"
	"webreduce/server"
)

var (
	httpAddr  string
	mongoAddr string
	testMode  bool
)

func init() {
	flag.StringVar(&httpAddr, "listen", "localhost:5000", "http address")
	flag.StringVar(&mongoAddr, "dbaddr", "localhost", "database address")
	flag.BoolVar(&testMode, "testmode", true, "test mode")
}

func main() {
	prefix := "/apps/"
	srv, err := server.NewServer(prefix, mongoAddr, testMode)
	if err != nil {
		panic(err)
	}

	http.Handle(prefix, srv)

	err = http.ListenAndServe(httpAddr, nil)
	if err != nil {
		panic(err)
	}
}
