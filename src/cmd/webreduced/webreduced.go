package main

import (
	"net/http"
	"webreduce/api"
	"webreduce/router"
)

func main() {
	cfg := map[string]string{
		"db/url":             "localhost",
		"db/name":            "webreduce",
		"db/collection/name": "agents",
	}

	as, err := api.NewAgentCollectionApi(cfg)
	if err != nil {
		panic(err)
	}
	defer as.Close()

	p := "/sinks"
	r := router.NewRouter(p)
	r.AddRoute("", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { as.GetList(ctx, w, r) }, "GET")

	http.Handle(p, &r)
	http.ListenAndServe(":8080", nil)
}
