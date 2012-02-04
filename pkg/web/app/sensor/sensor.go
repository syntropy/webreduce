package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"wr/interpreter/lua"
)

// Agents represent a persistable behaviour that collect and/or emit data
type Agent struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

// Validates the agent.
func (a *Agent) Valid() bool {
	if len(a.Name) < 1 {
		return false
	}

	if a.Language != "lua" {
		return false
	}

	if _, err := lua.New().Eval(a.Code); err != nil {
		return false
	}

	return true
}

// Calls the agent with data
func (a *Agent) Call(data interface{}) (err error) {
	lctx := lua.New()
	lctx.RegisterEmitCallback(func(data []byte) { fmt.Printf("EMIT: %v\n", string(data)) })

	fn, err := lctx.Eval(a.Code)
	if err != nil {
		return
	}

	fn(data.([]byte), []byte{})
	return
}

// AgentList represents a list of persisted Agents
type AgentList struct {
	Count int     `json:"count"`
	Items []Agent `json:"items"`
}

// The API for agent collections
type AgentCollectionApi struct {
	config    map[string]string
	dbsession *mgo.Session
}

func NewAgentCollectionApi(config map[string]string) (a AgentCollectionApi, err error) {
	a.config = config

	dburl, found := a.config["db/url"]
	if !found {
		err = errors.New("Missing config entry 'db/url'")
		return
	}

	dbname, found := a.config["db/name"]
	if !found {
		err = errors.New("Missing config entry 'db/name'")
		return
	}

	colname, found := a.config["db/collection/name"]
	if !found {
		err = errors.New("Missing config entry 'db/collection/name'")
		return
	}

	a.dbsession, err = mgo.Dial(dburl)
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

	col := a.dbsession.DB(dbname).C(colname)
	if err := col.EnsureIndex(index); err != nil {
		return a, err
	}

	return
}

func (api *AgentCollectionApi) Close() {
	api.dbsession.Close()
}

func (api *AgentCollectionApi) Collection() *mgo.Collection {
	return api.dbsession.Copy().DB(api.config["db/name"]).C(api.config["db/collection/name"])
}

// Get a list of agents
func (api *AgentCollectionApi) GetList(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	col := api.Collection()
	defer col.Database.Session.Close()

	query := col.Find(bson.M{})
	list := AgentList{Count: 0, Items: []Agent{}}

	count, err := query.Count()
	if err == nil {
		list.Count = count
		query.All(&list.Items)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}

// Put a named agent in the collection.
func (api *AgentCollectionApi) PutAgent(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	name, found := ctx["agent"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	agent := &Agent{Name: name}
	err := decoder.Decode(&agent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !(agent.Name == name && agent.Valid()) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.Collection()
	defer col.Database.Session.Close()

	if _, err := col.Upsert(bson.M{"name": name}, agent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get an agent by name
func (api *AgentCollectionApi) GetAgent(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	name, found := ctx["agent"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.Collection()
	defer col.Database.Session.Close()

	selector := bson.M{"name": name}
	agent := &Agent{}
	err := col.Find(selector).One(&agent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(agent)
}

// Post data to an agent
func (api *AgentCollectionApi) PostToAgent(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	name, found := ctx["agent"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.Collection()
	defer col.Database.Session.Close()

	agent := &Agent{}
	if err := col.Find(bson.M{"name": name}).One(&agent); err != nil {
		http.NotFound(w, r)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := agent.Call(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
