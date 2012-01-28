package api

import (
	"errors"
	"io"
	"launchpad.net/mgo"
	"net/http"
)

// Agents represent a persistable behaviour that collect and/or emit data
type Agent struct {
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

// Get a list of agents
func (api *AgentCollectionApi) GetList(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}
