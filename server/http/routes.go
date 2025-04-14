package http

import "net/http"

type RoutesHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

func SetUpRoutes(mux *http.ServeMux, routes []RoutesHandler) {
	for _, handler := range routes {
		handler.RegisterRoutes(mux)
	}
}
