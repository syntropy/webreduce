package app

import (
	"encoding/json"
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
	return
}

func (a *Api) Get(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, found := ctx.Get("app")
	if !found {
		http.NotFound(w, r)
		return
	}

	col := a.Collection()
	defer col.Database.Session.Close()

	app := &App{}
	if err := col.Find(bson.M{"name": name.(string)}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	wr.WriteJsonResponse(w, 200, app)

	return
}

func (a *Api) Put(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, found := ctx.Get("app")
	if !found {
		http.NotFound(w, r)
		return
	}

	decoder := json.NewDecoder(r.Body)
	app := &App{}
	if err := decoder.Decode(&app); err != nil || app.Name != name || !app.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := a.Collection()
	defer col.Database.Session.Close()

	if _, err := col.Upsert(bson.M{"name": name}, app); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
