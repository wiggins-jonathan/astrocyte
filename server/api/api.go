package api

import "net/http"

type API interface {
	Register(mux *http.ServeMux)
}
