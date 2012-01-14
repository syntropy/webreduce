package main

import (
	"flag"
	"fmt"
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

	flag.Parse()
}

func main() {
	prefix := "/apps/"
	srv, err := server.NewServer(prefix, mongoAddr, testMode)
	if err != nil {
		panic(err)
	}

	http.Handle(prefix, srv)

	fmt.Println("Listen ", httpAddr)

	err = http.ListenAndServe(httpAddr, nil)
	if err != nil {
		panic(err)
	}
}
