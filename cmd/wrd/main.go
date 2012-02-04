package main

import (
	"net/http"
	"wr"
	"wr/web/app/sensor"
	"wr/web/router"
)

func main() {
	cfg := &wr.StringContext{}
	data := map[string]string{
		"db/url":             "localhost",
		"db/name":            "webreduce",
		"db/collection/name": "sensors",
	}

	for k, v := range data {
		cfg.Set(k, v)
	}

	sensor, err := sensor.NewSensorCollectionApi(cfg)
	if err != nil {
		panic(err)
	}
	defer sensor.Close()

	r := router.NewRouter("/sensors")
	r.AddRoute("", func(ctx wr.Context, w http.ResponseWriter, r *http.Request) { sensor.GetList(ctx, w, r) }, "GET")
	r.AddRoute("/<sensor>", func(ctx wr.Context, w http.ResponseWriter, r *http.Request) { sensor.GetSensor(ctx, w, r) }, "GET")
	r.AddRoute("/<sensor>", func(ctx wr.Context, w http.ResponseWriter, r *http.Request) { sensor.PutSensor(ctx, w, r) }, "PUT")
	r.AddRoute("/<sensor>", func(ctx wr.Context, w http.ResponseWriter, r *http.Request) { sensor.PostToSensor(ctx, w, r) }, "POST")

	http.Handle("/", &r)
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.Dir("./static/test"))))
	http.ListenAndServe(":8080", nil)
}
