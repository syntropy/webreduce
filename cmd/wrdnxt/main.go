package main

import (
	"net/http"
	"wr"
	"wr/web/app"
	"wr/web/router"
)

func main() {
	apps, err := app.NewApi()
	if err != nil {
		panic(err)
	}
	defer apps.Close()

	r := router.NewRouter("/<app>")

	r.AddRoute("/", func(c wr.Context, w http.ResponseWriter, r *http.Request) {
		apps.GetApp(c, w, r)
	}, "GET")

	http.Handle("/", &r)
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.Dir("./static/test"))))
	http.ListenAndServe(":8080", nil)
}
