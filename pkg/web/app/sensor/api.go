package sensor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"wr"
)

// The API for sensor collections
type SensorCollectionApi struct {
	config    wr.Context
	dbsession *mgo.Session
}

func NewSensorCollectionApi(config wr.Context) (a SensorCollectionApi, err error) {
	a.config = config

	dburl, found := a.config.Get("db/url")
	if !found {
		err = errors.New("Missing config entry 'db/url'")
		return
	}

	dbname, found := a.config.Get("db/name")
	if !found {
		err = errors.New("Missing config entry 'db/name'")
		return
	}

	colname, found := a.config.Get("db/collection/name")
	if !found {
		err = errors.New("Missing config entry 'db/collection/name'")
		return
	}

	a.dbsession, err = mgo.Dial(dburl.(string))
	if err != nil {
		return a, err
	}

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	col := a.dbsession.DB(dbname.(string)).C(colname.(string))
	if err := col.EnsureIndex(index); err != nil {
		return a, err
	}

	return
}

func (api *SensorCollectionApi) Close() {
	api.dbsession.Close()
}

func (api *SensorCollectionApi) Collection() *mgo.Collection {
	dbname, found := api.config.Get("db/name")
	if !found {
		panic("missing conf")
	}

	colname, found := api.config.Get("db/collection/name")
	if !found {
		panic("missing conf")
	}

	return api.dbsession.Copy().DB(dbname.(string)).C(colname.(string))
}

// Get a list of sensors
func (api *SensorCollectionApi) GetList(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	col := api.Collection()
	defer col.Database.Session.Close()

	query := col.Find(bson.M{})
	list := SensorList{Count: 0, Items: []Sensor{}}

	count, err := query.Count()
	if err == nil {
		list.Count = count
		query.All(&list.Items)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}

// Put a named sensor in the collection.
func (api *SensorCollectionApi) PutSensor(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, found := ctx.Get("sensor")
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	sensor := &Sensor{Name: name.(string)}
	err := decoder.Decode(&sensor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !(sensor.Name == name && sensor.Valid()) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.Collection()
	defer col.Database.Session.Close()

	if _, err := col.Upsert(bson.M{"name": name}, sensor); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get an sensor by name
func (api *SensorCollectionApi) GetSensor(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, found := ctx.Get("sensor")
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.Collection()
	defer col.Database.Session.Close()

	selector := bson.M{"name": name}
	sensor := &Sensor{}
	err := col.Find(selector).One(&sensor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(sensor)
}

// Post data to an sensor
func (api *SensorCollectionApi) PostToSensor(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, found := ctx.Get("sensor")
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.Collection()
	defer col.Database.Session.Close()

	sensor := &Sensor{}
	if err := col.Find(bson.M{"name": name}).One(&sensor); err != nil {
		http.NotFound(w, r)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := sensor.Call(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
