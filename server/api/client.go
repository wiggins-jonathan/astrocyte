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

func (c *clientAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/matrix/client", c.WellKnownClientHandler)
	mux.HandleFunc("GET /_matrix/client/versions", c.VersionsHandler)
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

func (c *clientAPI) VersionsHandler(w http.ResponseWriter, r *http.Request) {
	// the versions endpoint has different behavior depending on if we are
	// authenticated or not. Add auth stuff later

	w.WriteHeader(http.StatusOK)

	versions := []string{"v1.14"}
	response := struct {
		Versions         []string        `json:"versions"`
		UnstableFeatures map[string]bool `json:"unstable_features,omitempty"`
	}{
		Versions: versions,
	}

	// json.Encode returns an optional error. We should handle that error here
	// & in WellKnownClientHandler. To do this, we may have to pass in our
	// logger from the server struct
	json.NewEncoder(w).Encode(response)
}

func WithBaseURL(baseURL *url.URL) ClientOption {
	return func(c *clientAPI) {
		c.BaseURL = baseURL
	}
}
