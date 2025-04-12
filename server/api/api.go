package api

import "net/http"

type API interface {
	RegisterRoutes(mux *http.ServeMux)
}
