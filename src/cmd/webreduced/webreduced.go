package main

import (
	"webreduce/router"
	"webreduce/api"
	"net/http"
)

func main() {
	cfg := map[string]string{
		"db/name": "webreduce",
		"db/collection": "sinks",
	}

	as := api.AgentCollectionApi{cfg}

	p := "/sinks"
	r := router.NewRouter(p)
	r.AddRoute("", func(ctx map[string]string, w http.ResponseWriter, r *http.Request){as.GetList(ctx, w, r)}, "GET")

	http.Handle(p, &r)
	http.ListenAndServe(":8080", nil)
}
