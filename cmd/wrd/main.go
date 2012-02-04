package main

import (
	"net/http"
	"wr/web/api"
	"wr/web/router"
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

	r := router.NewRouter("/agents")
	r.AddRoute("", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { as.GetList(ctx, w, r) }, "GET")
	r.AddRoute("/<agent>", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { as.GetAgent(ctx, w, r) }, "GET")
	r.AddRoute("/<agent>", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { as.PutAgent(ctx, w, r) }, "PUT")
	r.AddRoute("/<agent>", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { as.PostToAgent(ctx, w, r) }, "POST")

	http.Handle("/", &r)
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.Dir("./static/test"))))
	http.ListenAndServe(":8080", nil)
}
