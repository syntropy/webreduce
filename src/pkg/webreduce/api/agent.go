package api

import (
	"encoding/json"
	"errors"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
)

// Agents represent a persistable behaviour that collect and/or emit data
type Agent struct {
	Name	 string
	Language string
	Code     string
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
	var list []Agent
	query := api.Collection().Find(bson.M{})
	query.All(&list)

	count, err := query.Count()
	if err != nil {
		count = 0
	}

	res := make(map[string]interface{})
	res["Count"] = count
	res["Result"] = list

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

// Put an agents
func (api *AgentCollectionApi) PutAgent(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	name, found := ctx["agent"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lang := "lua"
	code := ""

	col := api.dbsession.Copy().DB(api.config["db/name"]).C(api.config["db/collection/name"])
	selector := bson.M{"name": name}
	agent := &Agent{Name: name, Language: lang, Code: code}
	_, err := col.Upsert(selector, agent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Get an agent by name
func (api *AgentCollectionApi) GetAgent(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	name, found := ctx["agent"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	col := api.dbsession.Copy().DB(api.config["db/name"]).C(api.config["db/collection/name"])

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
