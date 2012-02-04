package app

import (
	"net/http"
	"wr"
)

type Api struct {
}

func NewApi() (api Api, err error) {
	api = Api{}

	return
}

func (a *Api) Close() {}

func (a *Api) GetApp(ctx wr.Context, w http.ResponseWriter, r *http.Request) {
	name, found := ctx.Get("name")
	if !found {
		http.NotFound(w, r)
		return
	}

	app := App{name.(string)}

	wr.WriteJsonResponse(w, 200, app)

	return
}
