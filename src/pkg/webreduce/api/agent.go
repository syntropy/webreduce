package api

import (
	"io"
	"net/http"
)

type AgentCollectionApi struct {
	Config map[string]string
}

func (api *AgentCollectionApi) GetList(ctx map[string]string, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

