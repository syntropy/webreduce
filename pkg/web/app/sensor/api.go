package sensor

import (
	"io/ioutil"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"wr"
	"wr/messaging"
	"wr/web/app"
	"wr/web/router"
)

type Api struct {
	dbsession  *mgo.Session
	appcolname string
	colname    string
}

func NewApi(dburl, appcolname, colname string) (api Api, err error) {
	dbsession, err := mgo.Dial(dburl)
	if err != nil {
		return
	}

	api = Api{dbsession: dbsession, appcolname: appcolname, colname: colname}

	return
}

func (a *Api) Close() {
	a.dbsession.Close()
}

func (a *Api) AppCollection() *mgo.Collection {
	return a.dbsession.Clone().DB(wr.DBNAME).C(a.appcolname)
}

func (a *Api) Collection() *mgo.Collection {
	return a.dbsession.Clone().DB(wr.DBNAME).C(a.colname)
}

func (a *Api) RegisterRoutes(r *router.Router) {
	r.AddRoute("/sensors/<sensor>", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.Get(c, w, r) }, "GET")
	r.AddRoute("/sensors/<sensor>", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.PostMessage(c, w, r) }, "POST")
	r.AddRoute("/sensors/<sensor>", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.Put(c, w, r) }, "PUT")
	// r.AddRoute("/sensors/<sensor>", func(c wr.Context, w http.ResponseWriter, r *http.Request) { a.Delete(c, w, r) }, "DELETE")
	return
}

func (a *Api) Get(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	appname, _ := ctx.Get("app")
	sensorname, _ := ctx.Get("sensor")

	appcol := a.AppCollection()
	defer appcol.Database.Session.Close()

	app := &app.App{}
	if err := appcol.Find(bson.M{"name": appname.(string)}).Select(bson.M{"name": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	col := a.Collection()
	q := bson.M{"name": sensorname.(string)}
	sensor := &Sensor{}
	if err := col.Find(q).Select(bson.M{"name": 1}).One(&sensor); err != nil {
		http.NotFound(w, r)
		return
	}

	wr.WriteJsonResponse(w, 200, sensor)

	return
}

func (a *Api) Put(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	appname, _ := ctx.Get("app")
	sensorname, _ := ctx.Get("sensor")

	appcol := a.AppCollection()
	defer appcol.Database.Session.Close()

	app := &app.App{}
	if err := appcol.Find(bson.M{"name": appname.(string)}).Select(bson.M{"name": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	sensor := &Sensor{}
	if err := wr.ReadJsonRequest(r, sensor); err != nil || sensor.Name != sensorname || !sensor.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := a.Collection()
	defer col.Database.Session.Close()

	if _, err := col.Upsert(bson.M{"name": sensorname}, sensor); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	return
}

// func (a *Api) Delete(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
// 	name, found := ctx.Get("app")
// 	if !found {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	col := a.Collection()
// 	defer col.Database.Session.Close()

// 	if err := col.RemoveAll(bson.M{"name": name.(string)}); err != nil {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	return
// }

func (a *Api) PostMessage(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	appname, _ := ctx.Get("app")
	sensorname, _ := ctx.Get("sensor")

	appcol := a.AppCollection()
	defer appcol.Database.Session.Close()

	app := &app.App{}
	if err := appcol.Find(bson.M{"name": appname.(string)}).Select(bson.M{"name": 1}).One(&app); err != nil {
		http.NotFound(w, r)
		return
	}

	col := a.Collection()
	q := bson.M{"name": sensorname.(string)}
	sensor := &Sensor{}
	if err := col.Find(q).Select(bson.M{"name": 1}).One(&sensor); err != nil {
		http.NotFound(w, r)
		return
	}

	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	wr.MQ[appname.(string)].Pub <- messaging.NewMessage(msg)

	w.WriteHeader(http.StatusAccepted)

	return
}

// func (a *Api) GetIndex(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
// 	name, found := ctx.Get("app")
// 	if !found {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	col := a.Collection()
// 	defer col.Database.Session.Close()

// 	app := &App{}
// 	if err := col.Find(bson.M{"name": name.(string)}).Select(bson.M{"index": 1}).One(&app); err != nil {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	w.Header().Set("Content-type", "text/html")
// 	w.Write([]byte(app.Index))

// 	return
// }

// func (a *Api) PutIndex(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
// 	name, found := ctx.Get("app")
// 	if !found {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	col := a.Collection()
// 	defer col.Database.Session.Close()

// 	text, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 	}

// 	doc := bson.M{"$set": bson.M{"index": text}}

// 	if _, err := col.Upsert(bson.M{"name": name}, doc); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	return
// }
