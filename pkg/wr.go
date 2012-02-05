package wr

import (
	"encoding/json"
	"net/http"
	"wr/messaging"
)

const (
	DBNAME = "webreduce"
)

var MQ map[string]*PubSub

func init() {
	MQ = make(map[string]*PubSub)
}

type Context interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}) error
}

type StringContext struct {
	data map[string]string
}

func (c *StringContext) Get(key string) (val interface{}, found bool) {
	val, found = c.data[key]

	return
}

func (c *StringContext) Set(key string, value interface{}) (err error) {
	if c.data == nil {
		c.data = map[string]string{}
	}

	c.data[key] = value.(string)

	return
}

func WriteJsonResponse(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(obj)
}

func ReadJsonRequest(r *http.Request, obj interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&obj)

	return
}

type PubSub struct {
	Pub  chan *messaging.Message
	dev  *messaging.Device
	name string
}

func NewPubSub(name string) *PubSub {
	dev, _ := messaging.NewDevice()
	dev.StartPub(name)
	ps := &PubSub{Pub: dev.Out, dev: dev, name: name}

	return ps
}

func (ps *PubSub) Sub() chan *messaging.Message {
	dev, _ := messaging.NewDevice()
	dev.StartSub(ps.name)

	return dev.In
}
