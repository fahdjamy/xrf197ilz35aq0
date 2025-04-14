package handlers

import (
	"net/http"
	"xrf197ilz35aq0/core/service"
	"xrf197ilz35aq0/internal"
)

func SetupHandlers(serveMux *http.ServeMux, logger internal.Logger, services service.Services) {
	routes := make([]RoutesHandler, 0)

	// create handlers
	healthRoutesHandler := NewHealthRoutes(logger)
	authRoutesHandler := NewAuthHandler(logger, services)
	orgRoutesHandler := NewOrgHandler(logger, services.OrgService)
	userRoutesHandler := NewUserHandler(logger, services.UserService)
	permissionRoutesHandler := NewPermHandler(logger, services.PermissionService)

	routes = append(routes, orgRoutesHandler)
	routes = append(routes, userRoutesHandler)
	routes = append(routes, authRoutesHandler)
	routes = append(routes, healthRoutesHandler)
	routes = append(routes, permissionRoutesHandler)

	// register routes To the http mux server
	for _, handler := range routes {
		handler.RegisterRoutes(serveMux)
	}
}

type RoutesHandler interface {
	RegisterRoutes(serveMux *http.ServeMux)
}
