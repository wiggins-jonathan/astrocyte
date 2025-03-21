package api

import (
	"encoding/json"
	"net/http"
)

type clientAPI struct{}

func NewClient() *clientAPI {
	return &clientAPI{}
}

func (c *clientAPI) Register(mux *http.ServeMux) {
	mux.HandleFunc("/.well-known/matrix/client", WellKnownClientHandler)
}

func WellKnownClientHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	type serverInfo struct {
		BaseURL string `json:"base_url"`
	}

	response := struct {
		HomeServer     serverInfo     `json:"m.homeserver"`
		IdentityServer serverInfo     `json:"m.identity_server,omitempty"`
		CustomProps    map[string]any `json:"-"`
	}{
		HomeServer: serverInfo{
			BaseURL: "matrix.org",
		},
	}

	json.NewEncoder(w).Encode(response)
}
