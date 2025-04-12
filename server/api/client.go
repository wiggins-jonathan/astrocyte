package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type clientAPI struct {
	BaseURL *url.URL
}

type ClientOption func(*clientAPI)

func NewClient(options ...ClientOption) *clientAPI {
	client := &clientAPI{}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *clientAPI) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/matrix/client", c.WellKnownClientHandler)
}

func (c *clientAPI) WellKnownClientHandler(w http.ResponseWriter, r *http.Request) {
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
			BaseURL: c.BaseURL.Path,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func WithBaseURL(baseURL *url.URL) ClientOption {
	return func(c *clientAPI) {
		c.BaseURL = baseURL
	}
}
