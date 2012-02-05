package app

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io/ioutil"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"wr"
	"wr/web/router"
)

type Api struct {
	dbsession *mgo.Session
	colname   string
}

func NewApi(dburl, colname string) (api Api, err error) {
	dbsession, err := mgo.Dial(dburl)
	if err != nil {
		return
	}

	api = Api{dbsession: dbsession, colname: colname}

	return
}

func (a *Api) Close() {
	a.dbsession.Close()
}

func (a *Api) Collection() *mgo.Collection {
	return a.dbsession.Clone().DB(wr.DBNAME).C(a.colname)
}

func (a *Api) RegisterRoutes(r *router.Router) {
	r.AddRoute("", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.Get(c, w, r) }, "GET")
	r.AddRoute("", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.Put(c, w, r) }, "PUT")
	r.AddRoute("", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.Delete(c, w, r) }, "DELETE")
	r.AddRoute("", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.PostSignal(c, w, r) }, "POST")
	r.AddRoute("/", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.GetIndex(c, w, r) }, "GET")
	r.AddRoute("/index.html", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.GetIndex(c, w, r) }, "GET")
	r.AddRoute("/index.html", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.PutIndex(c, w, r) }, "PUT")
	r.AddRoute("/ws", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.GetWebsocket(c, w, r) }, "GET")
	return
}

func (a *Api) Get(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, _ := ctx.Get("app")

	col := a.Collection()
	defer col.Database.Session.Close()

	app := &App{}
	if err := col.Find(bson.M{"name": name.(string)}).Select(bson.M{"name": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	wr.WriteJsonResponse(w, 200, app)

	return
}

func (a *Api) Put(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, _ := ctx.Get("app")

	app := &App{}
	if err := wr.ReadJsonRequest(r, app); err != nil || app.Name != name || !app.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := a.Collection()
	defer col.Database.Session.Close()

	if _, err := col.Upsert(bson.M{"name": name}, app); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	return
}

func (a *Api) Delete(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, _ := ctx.Get("app")

	col := a.Collection()
	defer col.Database.Session.Close()

	if err := col.RemoveAll(bson.M{"name": name.(string)}); err != nil {
		http.NotFound(w, r)
		return
	}
	return
}

func (a *Api) PostSignal(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	rawname, _ := ctx.Get("app")
	name := rawname.(string)

	col := a.Collection()
	defer col.Database.Session.Close()

	app := &App{}
	if err := col.Find(bson.M{"name": name}).Select(bson.M{"name": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	wr.MQ[name] = wr.NewPubSub(name)

	go func() {
		for msg := range wr.MQ[name].Sub() {
			fmt.Printf("Got: %v\n", msg)
		}
	}()

	w.WriteHeader(http.StatusAccepted)

	return
}

func (a *Api) GetIndex(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, _ := ctx.Get("app")

	col := a.Collection()
	defer col.Database.Session.Close()

	app := &App{}
	if err := col.Find(bson.M{"name": name.(string)}).Select(bson.M{"index": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.Write([]byte(app.Index))

	return
}

func (a *Api) PutIndex(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, _ := ctx.Get("app")

	col := a.Collection()
	defer col.Database.Session.Close()

	text, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	doc := bson.M{"$set": bson.M{"index": text}}

	if _, err := col.Upsert(bson.M{"name": name}, doc); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	return
}

func (a *Api) GetWebsocket(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, _ := ctx.Get("app")

	col := a.Collection()
	defer col.Database.Session.Close()

	app := &App{}
	if err := col.Find(bson.M{"name": name.(string)}).Select(bson.M{"name": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	go websocket.Handler(func(ws *websocket.Conn) {
		for msg := range wr.MQ[name.(string)].Sub() {
			ws.Write(msg.Payload)
		}
	}).ServeHTTP(w, r)

	return
}
