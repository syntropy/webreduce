package server

import (
	// "launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
)

type Server struct {
	Prefix    string
	DbSession *mgo.Session
	Db mgo.Database
	AppCollection mgo.Collection
}

type AppConfig struct {
	Name string
}

func NewServer(prefix string, dbaddr string, testmode bool) (srv *Server, err error) {
	session, err := mgo.Mongo(dbaddr)
	if err != nil {
		return
	}

	dbname := "webreduce"

	if testmode {
		dbname = "test"+dbname
	}

	db := session.DB(dbname)

	if testmode {
		if err := db.DropDatabase(); err != nil {
			panic(err)
		}
	}

	appCollection := db.C("apps")

	index := mgo.Index{
	    Key: []string{"name"},
	    Unique: true,
	    DropDups: true,
	    Background: true,
	    Sparse: true,
	}

	idxErr := appCollection.EnsureIndex(index)
	if idxErr != nil {
		panic(idxErr)
	}

	srv = &Server{prefix, session, db, appCollection}

	return
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handle errors by returning a HTTP 500 status code to the client
	defer func() {
		if e := recover(); e != nil {
			trace := debug.Stack()
			msg := fmt.Sprintf("Internal Server Error\n\n%s\n%s", trace, e)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}()

	method := req.Method
	path := req.URL.Path

	namePattern := "([a-z][a-z0-9]*)"
	regex := regexp.MustCompile(srv.Prefix+namePattern)

	matches := regex.FindStringSubmatch(path)
	if len(matches) < 2 {
		http.NotFound(w, req)
		return
	}

	appName := matches[1]

	if method == "PUT" {
		srv.HandlePutApp(appName, w, req)
		return
	}

	http.NotFound(w, req)
}

func (srv *Server) HandlePutApp(name string, w http.ResponseWriter, req *http.Request) {
	selector := &AppConfig{name}
	change := &AppConfig{name}

	id, err := srv.AppCollection.Upsert(selector, change)

	if err != nil {
		panic(err)
	}

	if id != nil {
		http.Redirect(w, req, srv.Prefix+name, http.StatusCreated)
		return
	}

	http.Error(w, "", http.StatusNoContent)
}
