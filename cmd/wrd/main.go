package main

import (
	"net/http"
	"wr/web/app/sensor"
	"wr/web/router"
)

func main() {
	cfg := map[string]string{
		"db/url":             "localhost",
		"db/name":            "webreduce",
		"db/collection/name": "sensors",
	}

	sensor, err := sensor.NewSensorCollectionApi(cfg)
	if err != nil {
		panic(err)
	}
	defer sensor.Close()

	r := router.NewRouter("/sensors")
	r.AddRoute("", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { sensor.GetList(ctx, w, r) }, "GET")
	r.AddRoute("/<sensor>", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { sensor.GetSensor(ctx, w, r) }, "GET")
	r.AddRoute("/<sensor>", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { sensor.PutSensor(ctx, w, r) }, "PUT")
	r.AddRoute("/<sensor>", func(ctx map[string]string, w http.ResponseWriter, r *http.Request) { sensor.PostToSensor(ctx, w, r) }, "POST")

	http.Handle("/", &r)
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.Dir("./static/test"))))
	http.ListenAndServe(":8080", nil)
}
