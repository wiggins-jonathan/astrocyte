package api

import (
	"encoding/json"
	"net/http"
)

type pushAPI struct{}

func NewPushAPI() *pushAPI {
	return &pushAPI{}
}

func (p *pushAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /_matrix/push/v1/notify", PushGatewayHandler)
}

func PushGatewayHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := struct {
		Rejected []string `json:"rejected"`
	}{
		Rejected: []string{"hello, from the push gateway!"},
	}

	json.NewEncoder(w).Encode(response)
}
