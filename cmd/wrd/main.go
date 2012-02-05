package main

import (
	"net/http"
	"wr/web/app"
	"wr/web/app/sensor"
	"wr/web/router"
)

func main() {
	dburl := ""
	apps, err := app.NewApi(dburl, "apps")
	if err != nil {
		panic(err)
	}
	defer apps.Close()

	r := router.NewRouter("/<app>")
	apps.RegisterRoutes(&r)

	sensors, err := sensor.NewApi(dburl, "apps", "sensors")
	if err != nil {
		panic(err)
	}
	defer apps.Close()

	sensors.RegisterRoutes(&r)

	http.Handle("/", &r)
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.Dir("./static/test"))))
	http.ListenAndServe(":8080", nil)
}
